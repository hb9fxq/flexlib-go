package main

import "C"

import (
	"bufio"
	"fmt"
	"github.com/hb9fxq/flexlib-go/obj"
	"github.com/hb9fxq/flexlib-go/sdrobjects"
	"github.com/smallnest/ringbuffer"
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
	readBuffer                 *ringbuffer.RingBuffer
}

var gloabalAppCtx *AppContext
var initSuccess bool
var firstPan = ""
var radioContext *obj.RadioContext
var l = log.New(os.Stderr, "IQTRANSFERLIB ", 0)

//export InitRadio
func InitRadio(radioAddr string, myudp string, channel string, rate string) {

	if initSuccess {
		return
	}

	appContext := new(AppContext)
	appContext.radioAddr = radioAddr
	appContext.myPort = myudp
	appContext.daxIqChan = channel
	appContext.sampleRate = rate
	appContext.readBuffer = ringbuffer.New(16 * 100000)
	gloabalAppCtx = appContext

	if appContext.sampleRate != "24000" && appContext.sampleRate != "48000" && appContext.sampleRate != "96000" && appContext.sampleRate != "192000" {
		panic("Invalid Sample Rate! Allowed values 24000, 48000, 96000, 192000")
	}

	if len(appContext.forwardAddess) > 0 {
		appContext.forwardConnection, _ = net.Dial("udp", appContext.forwardAddess)
	}

	radioContext = new(obj.RadioContext)
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
	obj.SendRadioCommand(radioContext, "client set enforce_network_mtu=1 network_mtu=1420")
	l.Println("Requesting UDP VITA data to be sent to " + appContext.myPort)

	obj.SendRadioCommand(radioContext, "client udpport "+appContext.myPort)

	appContext.RadioReponseStreamSequence = obj.SendRadioCommand(radioContext, "stream create type=dax_iq daxiq_channel=1")
	obj.SendRadioCommand(radioContext, "client set enforce_network_mtu=1 network_mtu=1420")
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

	initSuccess = true

}

var dataAvail bool

//export SetFrequency
func SetFrequency(freq int64) {

	if firstPan == "" {
		return
	}

	freqMhz := float64(freq) / 1000000

	s := fmt.Sprintf("%.6f", freqMhz)

	l.Println(s)
	l.Println("Instructed pan " + firstPan + " to tune to " + s + " MHz")
	obj.SendRadioCommand(radioContext, "display pan set "+firstPan+" center="+s)

}

//export ReadStream3
func ReadStream3(elements int) (C.size_t, *C.uchar) {

	size := elements
	p := C.malloc(C.size_t(size))

	bytes := (*[1<<30 - 1]C.uchar)(p)[:size:size]

	for {
		if bytesWritten >= size {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if bytesWritten < size {
		return C.size_t(0), (*C.uchar)(p)
	}

	buf := make([]byte, size)
	readSizze, _ := gloabalAppCtx.readBuffer.Read(buf)
	bytesWritten -= elements

	for i := 0; i < readSizze; i++ {
		bytes[i] = C.uchar(buf[i])
	}

	//fmt.Printf("%s", hex.Dump(buf))

	return C.size_t(size), (*C.uchar)(p)
}

var LastBlock []byte

var bytesWritten = 0

func handleData(appctx *AppContext, ifDataPackage sdrobjects.SdrIfData) {

	if uint64(ifDataPackage.Stream_id) != appctx.RadioResponseStream {
		return
	}
	//fmt.Printf("%s", hex.Dump(ifDataPackage.Data))
	appctx.readBuffer.Write(ifDataPackage.Data)
	bytesWritten += len(ifDataPackage.Data)
}

func main() {}
