package main

import (
	"flag"
	"github.com/krippendorf/flexlib-go/obj"
	"github.com/krippendorf/flexlib-go/sdrobjects"
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
	forwardConnection          net.Conn
}

func main() {

	l := log.New(os.Stderr, "RADIO_MSG", 0)

	appContext := new(AppContext)
	flag.StringVar(&appContext.radioAddr, "RADIO", "", "IP ADDRESS OF THE RADIO e.g 192.168.41.8")
	flag.StringVar(&appContext.myPort, "MYUDP", "", "LOCAL UDP PORT 7788")
	flag.StringVar(&appContext.daxIqChan, "CH", "", "DAX IQ CHANNEL NUMBER e.g. ")
	flag.StringVar(&appContext.sampleRate, "RATE", "", "DAX IQ sample rate in kHz - 24 / 48 / 96 / 192")
	flag.StringVar(&appContext.forwardAddess, "FWD", "", "If empty, IQ data will be written to stdout. UDP Forward address for the IQ samples with port, e.g. 192.168.50.5:5000")
	flag.Parse()

	if appContext.sampleRate != "24" && appContext.sampleRate != "48" && appContext.sampleRate != "96" && appContext.sampleRate != "192" {
		panic("Invalid Sample Rate! Allowed values 24, 48, 96, 192")
	}

	if len(appContext.forwardAddess) > 0 {
		appContext.forwardConnection, _ = net.Dial("udp", appContext.forwardAddess)
	}

	radioContext := new(obj.RadioContext)
	radioContext.RadioAddr = appContext.radioAddr
	radioContext.MyUdpEndpointPort = appContext.myPort
	radioContext.ChannelRadioResponse = make(chan string)
	radioContext.ChannelVitaIfData = make(chan *sdrobjects.SdrIfData)
	radioContext.Debug = true

	go obj.InitRadioContext(radioContext)

	go func(ctx *obj.RadioContext) {
		for {
			response := <-ctx.ChannelRadioResponse

			if strings.HasPrefix(response, "R"+strconv.Itoa(appContext.RadioReponseStreamSequence)) {
				cmd := "daxiq set" + appContext.daxIqChan + " rate=" + appContext.sampleRate + "000"
				obj.SendRadioCommand(radioContext, cmd)
			}
		}
	}(radioContext)

	go func(ctx *obj.RadioContext) {
		for { /* we'll only receive the samples for the stream requested on that port so we can ignore the stream id*/
			handleData(appContext, *<-ctx.ChannelVitaIfData)
		}
	}(radioContext)

	for {
		if len(radioContext.RadioHandle) > 0 { // wait until we got our handle
			break
		}
		time.Sleep(500)
	}

	appContext.RadioReponseStreamSequence = obj.SendRadioCommand(radioContext, "stream create daxiq=1ip="+radioContext.MyUdpEndpointIP.String()+" port="+appContext.myPort)

	var centerFrequencyOfStream int32

	for {

		if val, ok := radioContext.IqStreams.Load(appContext.daxIqChan); ok {
			iqStream := val.(obj.IqStream)
			if len(iqStream.Pan) > 0 {

				if panVal, panOk := radioContext.Panadapters.Load(iqStream.Pan); panOk {
					pan := panVal.(obj.Panadapter)
					if centerFrequencyOfStream != pan.Center {
						centerFrequencyOfStream = pan.Center
						l.Println("CENTER_FREQ_CHANGE " + strconv.Itoa(int(centerFrequencyOfStream)))
					}
				}

			}
		}
	}

	forever := make(chan bool)
	forever <- true
}

func handleData(appctx *AppContext, ifDataPackage sdrobjects.SdrIfData) {
	if len(appctx.forwardAddess) > 0 {
		appctx.forwardConnection.Write(ifDataPackage.Data)
	} else {
		os.Stdout.Write(ifDataPackage.Data)
	}
}

func FloatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}