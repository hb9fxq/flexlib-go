/* 2017 by Frank Werner-Krippendorf / HB9FXQ, mail@hb9fxq.ch
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package vita

type VitaPacketType int

const (
	IFData            VitaPacketType = iota
	IFDataWithStream  VitaPacketType = iota
	ExtData           VitaPacketType = iota
	ExtDataWithStream VitaPacketType = iota
	IFContext         VitaPacketType = iota
	ExtContext        VitaPacketType = iota
)

type VitaTimeStampIntegerType uint

const (
	NoneTsi VitaTimeStampIntegerType = iota
	UTC     VitaTimeStampIntegerType = iota
	GPS     VitaTimeStampIntegerType = iota
	Other   VitaTimeStampIntegerType = iota
)

type VitaTimeStampFractionalType uint

const (
	NoneTsf     VitaTimeStampFractionalType = iota
	SampleCount VitaTimeStampFractionalType = iota
	RealTime    VitaTimeStampFractionalType = iota
	FreeRunning VitaTimeStampFractionalType = iota
)

type VitaClassID struct {
	OUI                  uint32
	InformationClassCode uint16
	PacketClassCode      uint16
}

type VitaTrailer struct {
	CalibratedTimeEnable       bool
	ValidDataEnable            bool
	ReferenceLockEnable        bool
	AGCMGCEnable               bool
	DetectedSignalEnable       bool
	SpectralInversionEnable    bool
	OverrangeEnable            bool
	SampleLossEnable           bool
	CalibratedTimeIndicator    bool
	ValidDataIndicator         bool
	ReferenceLockIndicator     bool
	AGCMGCIndicator            bool
	DetectedSignalIndicator    bool
	SpectralInversionIndicator bool
	OverrangeIndicator         bool
	SampleLossIndicator        bool
}

type VitaHeader struct {
	Pkt_type             VitaPacketType
	C                    bool
	T                    bool
	Tsi                  VitaTimeStampIntegerType
	Tsf                  VitaTimeStampFractionalType
	Packet_count         uint16
	Packet_size          uint16
	payload_cutoff_bytes int
}

type VitaPacketPreamble struct {
	Header         *VitaHeader
	Stream_id      uint32
	Class_id       *VitaClassID
	Timestamp_int  uint32
	Timestamp_frac uint64
}

type VitaIfData struct {
	Header         *VitaHeader
	Stream_id      uint32
	Class_id_h     uint32
	Class_id_l     uint32
	Timestamp_int  uint32
	Timestamp_frac uint64
	Payload        []float32
}

const (
	SL_VITA_DISCOVERY_CLASS      = uint16(0xFFFF)
	SL_VITA_METER_CLASS          = uint16(0x8002)
	SL_VITA_FFT_CLASS            = uint16(0x8003)
	SL_VITA_WATERFALL_CLASS      = uint16(0x8004)
	SL_VITA_OPUS_CLASS           = uint16(0x8005)
	SL_VITA_IF_NARROW_CLASS      = uint16(0x03E3)
	SL_VITA_IF_WIDE_CLASS_24kHz  = uint16(0x02E3)
	SL_VITA_IF_WIDE_CLASS_48kHz  = uint16(0x02E4)
	SL_VITA_IF_WIDE_CLASS_96kHz  = uint16(0x02E5)
	SL_VITA_IF_WIDE_CLASS_192kHz = uint16(0x02E6)
	MAX_VITA_PACKET_SIZE         = uint16(16384)
	FLEX_OUI                     = uint16(0x1C2D)
)
