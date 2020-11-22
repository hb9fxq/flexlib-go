package main

import (
	"bufio"
	"flag"
	"github.com/hb9fxq/flexlib-go/obj"
	"github.com/hb9fxq/flexlib-go/sdrobjects"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type AppContext struct {
	radioAddr                  string
	myPort                     string
	daxIqChan                  string
	sampleRate                 string
	forwardAddess              string
	RadioReponseStreamSequence int
	RadioResponseStream        uint64
	forwardConnection          net.Conn
	fCenterArg                 string
}

func main() {

	l := log.New(os.Stderr, "IQTRANSFER ", 0)

	appContext := new(AppContext)
	flag.StringVar(&appContext.radioAddr, "RADIO", "", "IP ADDRESS OF THE RADIO e.g 192.168.41.8")
	flag.StringVar(&appContext.myPort, "MYUDP", "", "LOCAL UDP PORT 7788")
	flag.StringVar(&appContext.daxIqChan, "CH", "", "DAX IQ CHANNEL NUMBER e.g. ")
	flag.StringVar(&appContext.sampleRate, "RATE", "", "DAX IQ sample rate in kHz - 24 / 48 / 96 / 192")
	flag.StringVar(&appContext.forwardAddess, "FWD", "", "If empty, IQ data will be written to stdout. UDP Forward address for the IQ samples with port, e.g. 192.168.50.5:5000")
	flag.StringVar(&appContext.fCenterArg, "FCENTER", "", "(Optional) Tune panadapter to initial fCenter (Mhz) e.g. 7.1")
	flag.Parse()

	if appContext.sampleRate != "24000" && appContext.sampleRate != "48000" && appContext.sampleRate != "96000" && appContext.sampleRate != "192000" {
		panic("Invalid Sample Rate! Allowed values 24000, 48000, 96000, 192000")
	}

	if len(appContext.forwardAddess) > 0 {
		appContext.forwardConnection, _ = net.Dial("udp", appContext.forwardAddess)
	}

	radioContext := new(obj.RadioContext)
	radioContext.RadioAddr = appContext.radioAddr
	radioContext.MyUdpEndpointPort = appContext.myPort
	radioContext.ChannelRadioResponse = make(chan string)
	radioContext.ChannelVitaIfData = make(chan *sdrobjects.SdrIfData)
	radioContext.Debug = false

	go func(ctx *obj.RadioContext) {
		for {
			response := <-ctx.ChannelRadioResponse

			if strings.HasPrefix(response, "R"+strconv.Itoa(appContext.RadioReponseStreamSequence)) {
				streamHexString := strings.Split(response, "|")[2]
				l.Println("Stream filter streamId 0x" + streamHexString)
				stream, _ := strconv.ParseUint(streamHexString, 16, 64)
				appContext.RadioResponseStream = stream
				obj.SendRadioCommand(radioContext, "stream set 0x"+streamHexString+" daxiq_rate="+appContext.sampleRate)
			}
		}
	}(radioContext)

	go func(ctx *obj.RadioContext) {
		for { /* we'll only receive the samples for the stream requested on that port so we can ignore the stream id*/
			handleData(appContext, *<-ctx.ChannelVitaIfData)
		}
	}(radioContext)

	go obj.InitRadioContext(radioContext)

	for {
		if len(radioContext.RadioHandle) > 0 { // wait until we got our handle
			break
		}
		time.Sleep(500)
	}

	// wait for first clientId
	var firstClient = ""

	if radioContext.Debug {
		l.Println("waiting for first client")
	}

	for {

		radioContext.Clients.Range(func(k interface{}, value interface{}) bool {
			firstClient = value.(obj.Client).ClientId
			return true
		})

		if firstClient != "" {
			break
		}
	}

	// wait for first panadapter
	var firstPan = ""

	if radioContext.Debug {
		l.Println("waiting for first panadapter")
	}

	for {

		radioContext.Panadapters.Range(func(k interface{}, value interface{}) bool {
			firstPan = value.(obj.Panadapter).Id
			return true
		})

		if firstPan != "" {
			break
		}
	}
	l.Println("Binding to client_id " + firstClient)
	obj.SendRadioCommand(radioContext, "client bind client_id="+firstClient)
	l.Println("Requesting UDP VITA data to be sent to " + appContext.myPort)

	obj.SendRadioCommand(radioContext, "client udpport "+appContext.myPort)

	appContext.RadioReponseStreamSequence = obj.SendRadioCommand(radioContext, "stream create type=dax_iq daxiq_channel=1")

	l.Println("binding to panadapter " + firstPan)

	obj.SendRadioCommand(radioContext, "dax iq set 1 pan="+firstPan+" rate="+appContext.sampleRate)

	if len(appContext.forwardAddess) > 0 {
		l.Println("Forwarding data to " + appContext.forwardAddess)
	}

	go func(ctx *obj.RadioContext) {

		if appContext.fCenterArg != "" {
			obj.SendRadioCommand(radioContext, "display pan set "+firstPan+" center="+appContext.fCenterArg)
			l.Println("Instructed pan " + firstPan + " to tune to " + appContext.fCenterArg + " MHz")
		}

		for {

			reader := bufio.NewReader(os.Stdin)
			l.Print("Listening for tuning instruction (MHz) at stdin")
			text, _ := reader.ReadString('\n')
			text = strings.Trim(text, " ")
			text = strings.Trim(text, "\n")
			if _, err := strconv.ParseFloat(text, 32); err == nil {
				//obj.SendRadioCommand(radioContext, "display pan s")
				obj.SendRadioCommand(radioContext, "display pan set "+firstPan+" center="+text)
				l.Println("Instructed pan " + firstPan + " to tune to" + text + " MHz")

			}

		}
	}(radioContext)

	forever := make(chan bool)
	forever <- true
}

func handleData(appctx *AppContext, ifDataPackage sdrobjects.SdrIfData) {

	if uint64(ifDataPackage.Stream_id) != appctx.RadioResponseStream {
		return
	}

	if len(appctx.forwardAddess) > 0 {
		appctx.forwardConnection.Write(ifDataPackage.Data)
	} else {
		os.Stdout.Write(ifDataPackage.Data)
	}
}
