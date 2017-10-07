package obj

import (
	"github.com/krippendorf/flexlib-go/sdrobjects"
	"github.com/krippendorf/flexlib-go/vita"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type RadioData struct {
	Preampble *vita.VitaPacketPreamble
	Payload   []byte
	LastErr   error
}

type RadioContext struct {
	RadioAddr                string
	RadioCmdSeqNumber        int
	RadioConn                *net.TCPConn
	ChannelRadioData         chan *RadioData
	ChannelRadioResponse     chan string
	RadioHandle              string
	MyUdpEndpointIP          *net.IP
	MyUdpEndpointPort        string // we need strings for all cmds....
	ChannelVitaFFT           chan *sdrobjects.SdrFFTPacket
	ChannelVitaOpus          chan []byte
	ChannelVitaIfData        chan *sdrobjects.SdrIfData
	ChannelVitaMeter         chan *sdrobjects.SdrMeterPacket
	ChannelVitaWaterfallTile chan *sdrobjects.SdrWaterfallTile
	Panadapters              map[string]Panadapter
	Debug                    bool
}

func getNextCommandPrefix(ctx *RadioContext) (string, int) {
	ctx.RadioCmdSeqNumber += 1
	return "C" + strconv.Itoa(ctx.RadioCmdSeqNumber) + "|", ctx.RadioCmdSeqNumber
}

func SendRadioCommand(ctx *RadioContext, cmd string) int {

	prefixString, sequence := getNextCommandPrefix(ctx)
	_, err := ctx.RadioConn.Write([]byte(prefixString + cmd + "\r"))

	if err != nil {
		panic(err)
	}

	return sequence
}

func GetOutboundIP() *net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return &localAddr.IP
}

func InitRadioContext(ctx *RadioContext) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", ctx.RadioAddr+":4992")

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	// dial TCP connection to radio
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	ctx.RadioConn = conn

	ctx.MyUdpEndpointIP = GetOutboundIP()

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	if err != nil {
		log.Println(err)
		panic(err)
	}

	go subscribeRadioUdp(ctx)
	go subscribeRadioUpdates(conn, ctx)

	// Subscribe data from radio
	SendRadioCommand(ctx, "sub tx all")
	SendRadioCommand(ctx, "sub atu all")
	SendRadioCommand(ctx, "sub amplifier all")
	SendRadioCommand(ctx, "sub meter all")
	SendRadioCommand(ctx, "sub pan all")
	SendRadioCommand(ctx, "sub slice all")
	SendRadioCommand(ctx, "sub gps all")
	SendRadioCommand(ctx, "sub audio_stream all")
	SendRadioCommand(ctx, "sub cwx all")
	SendRadioCommand(ctx, "sub xvtr all")
	SendRadioCommand(ctx, "sub memories all")
	SendRadioCommand(ctx, "sub daxiq all")
	SendRadioCommand(ctx, "sub dax all")
	SendRadioCommand(ctx, "sub usb_cable all")
}

func subscribeRadioUpdates(conn *net.TCPConn, ctx *RadioContext) {

	l := log.New(os.Stderr, "RADIO_MSG", 0)
	buf := make([]byte, 4096)

	for {
		n, err := conn.Read(buf)

		if err != nil {
			continue
		}

		response := string(buf[:n])

		if len(response) == 0 {
			continue
		}

		lines := strings.Split(response, "\n")

		for _, responseLine := range lines {

			if len(strings.Trim(responseLine, " ")) == 0 {
				continue
			}

			if len(ctx.RadioHandle) == 0 && strings.HasPrefix(strings.ToUpper(responseLine), "H") {
				ctx.RadioHandle = responseLine[1:]
				l.Println("\nMY_RADIO_HANDLE>>" + ctx.RadioHandle)
			} else {

				if nil == ctx.ChannelRadioResponse {
					l.Println("Respnse Channel not bound: " + responseLine)
				} else {
					ctx.ChannelRadioResponse <- responseLine
					parseResponseLine(ctx, responseLine)
					if(ctx.Debug){
						l.Println("DEBU:RESP:" + responseLine)
					}
				}
			}
		}

		if err != nil {
			l.Println(err)
		}
	}
}
func parseResponseLine(context *RadioContext, respLine string) {
	if strings.Contains(respLine, "display pan") {
		parsePanAdapterParams(context, respLine)
	}
}
func parsePanAdapterParams(context *RadioContext, i string) {

	if context.Panadapters == nil {
		context.Panadapters = map[string]Panadapter{}
	}

	tokens := strings.Split(i, " ")

	var pan Panadapter

	if val, ok := context.Panadapters[tokens[2]]; ok {
		pan = val
	}

	for rngAttr := range tokens[3:] {

		if strings.Index(tokens[rngAttr+3], "=") < 0 {
			continue
		}

		attrName := strings.Split(tokens[rngAttr+3], "=")[0]
		val := strings.Split(tokens[rngAttr+3], "=")[1]

		switch attrName {
		case "center":
			floatVal, _ := strconv.ParseFloat(val, 32)
			pan.center = float32(floatVal)
			break
		}

	}

	context.Panadapters[tokens[2]] = pan

}

func subscribeRadioUdp(ctx *RadioContext) {

	FLexBroadcastAddr, err := net.ResolveUDPAddr("udp", ctx.MyUdpEndpointIP.String()+":"+ctx.MyUdpEndpointPort)

	if err != nil {
		panic(err)
	}

	ServerConn, err := net.ListenUDP("udp", FLexBroadcastAddr)

	if err != nil {
		panic(err)
	}

	defer ServerConn.Close()
	buf := make([]byte, 64000)

	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	for {
		n, _, _ := ServerConn.ReadFromUDP(buf)
		radioData := new(RadioData)
		radioData.LastErr, radioData.Preampble, radioData.Payload = vita.ParseVitaPreamble(buf[:n])
		if ctx.ChannelRadioData != nil {
			ctx.ChannelRadioData <- radioData
		}

		dispatchDataToChannels(ctx, radioData)
	}
}

func dispatchDataToChannels(ctx *RadioContext, data *RadioData) {
	switch data.Preampble.Header.Pkt_type {

	case vita.ExtDataWithStream:

		switch data.Preampble.Class_id.PacketClassCode {

		case vita.SL_VITA_FFT_CLASS:
			if nil != ctx.ChannelVitaFFT {
				ctx.ChannelVitaFFT <- vita.ParseVitaFFT(data.Payload, data.Preampble)
			}
			break
		case vita.SL_VITA_OPUS_CLASS:
			if nil != ctx.ChannelVitaOpus {
				ctx.ChannelVitaOpus <- data.Payload[:len(data.Payload)-data.Preampble.Header.Payload_cutoff_bytes]
			}

			break
		case vita.SL_VITA_IF_NARROW_CLASS:
			if nil != ctx.ChannelVitaIfData {
				vita.ParseFData(data.Payload, data.Preampble)
			}
			break
		case vita.SL_VITA_METER_CLASS:
			if nil != ctx.ChannelVitaMeter {
				ctx.ChannelVitaMeter <- vita.ParseVitaMeterPacket(data.Payload, data.Preampble)
			}

			break
		case vita.SL_VITA_DISCOVERY_CLASS:
			// maybe later - we use static addresses
			break
		case vita.SL_VITA_WATERFALL_CLASS:
			if nil != ctx.ChannelVitaWaterfallTile {
				vita.ParseVitaWaterfall(data.Payload, data.Preampble)
			}
			break
		default:
			break
		}

		break

	case vita.IFDataWithStream:
		switch data.Preampble.Class_id.PacketClassCode {
		case vita.SL_VITA_IF_WIDE_CLASS_24kHz:
			fallthrough
		case vita.SL_VITA_IF_WIDE_CLASS_48kHz:
			fallthrough
		case vita.SL_VITA_IF_WIDE_CLASS_96kHz:
			fallthrough
		case vita.SL_VITA_IF_WIDE_CLASS_192kHz:
			if nil != ctx.ChannelVitaIfData {
				ctx.ChannelVitaIfData <- vita.ParseFData(data.Payload, data.Preampble)
			}
		}
		break
	}
}
