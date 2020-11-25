/* 2017 by Frank Werner-hb9fxq / HB9FXQ, mail@hb9fxq.ch
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

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/hb9fxq/flexlib-go/sdrobjects"
	"math"
)

var ONE_OVER_ZERO_DBFS = 1.0 / math.Pow(2, 15)

func ParseVitaPreamble(data []byte) (error, *VitaPacketPreamble, []byte) {

	if len(data) < 20 {
		return errors.New("not a VITA49 Package"), nil, nil
	}

	var vitaPacketPreamble VitaPacketPreamble

	var header VitaHeader
	var vitaClassId VitaClassID

	vitaPacketPreamble.Header = &header

	index := 0
	rHeader := binary.BigEndian.Uint32(data[0:4])
	index += 4
	header.Pkt_type = VitaPacketType((rHeader >> 28))
	header.C = ((rHeader & 0x08000000) != 0)
	header.T = ((rHeader & 0x04000000) != 0)
	header.Tsi = VitaTimeStampIntegerType(((rHeader >> 22) & 0x03))
	header.Tsf = VitaTimeStampFractionalType(((rHeader >> 20) & 0x03))
	header.Packet_count = uint16(((rHeader >> 16) & 0x0F))
	header.Packet_size = uint16(rHeader & 0xFFFF)

	if header.Pkt_type == IFDataWithStream || header.Pkt_type == ExtDataWithStream {
		vitaPacketPreamble.Stream_id = binary.BigEndian.Uint32(data[index : index+4])
		index += 4
	}

	if header.C {
		vitaPacketPreamble.Class_id = &vitaClassId
		temp := binary.BigEndian.Uint32(data[index : index+4])
		index += 4
		vitaClassId.OUI = temp & 0x00FFFFFF

		temp = binary.BigEndian.Uint32(data[index : index+4])
		index += 4
		vitaClassId.InformationClassCode = uint16((temp >> 16))
		vitaClassId.PacketClassCode = uint16(temp)
	}

	if header.Tsi != NoneTsi {
		vitaPacketPreamble.Timestamp_int = binary.BigEndian.Uint32(data[index : index+4])
		index += 4
	}

	if header.Tsf != NoneTsf {
		vitaPacketPreamble.Timestamp_frac = binary.BigEndian.Uint64(data[index : index+8])
		index += 8
	}

	if header.T {
		//index += 4
		header.Payload_cutoff_bytes = 4
	}

	return nil, &vitaPacketPreamble, data[index:]
}

func ParseVitaFFT(data []byte, preamble *VitaPacketPreamble) *sdrobjects.SdrFFTPacket {

	index := 0
	var fftPacket sdrobjects.SdrFFTPacket

	fftPacket.StartBin_index = binary.BigEndian.Uint16(data[index : index+2])
	index += 2

	fftPacket.NumBins = binary.BigEndian.Uint16(data[index : index+2])
	index += 2

	fftPacket.BinSize = binary.BigEndian.Uint16(data[index : index+2])
	index += 2

	fftPacket.TotalBinsInFrame = binary.BigEndian.Uint16(data[index : index+2])
	index += 2

	fftPacket.FrameIndex = binary.BigEndian.Uint32(data[index : index+4])
	index += 4

	for i := 0; i < int(fftPacket.NumBins)*2; i += 2 {
		fftPacket.Payload = append(fftPacket.Payload, binary.BigEndian.Uint16(data[i+index:i+index+2]))
	}

	return &fftPacket
}

func ParseVitaMeterPacket(data []byte, preamble *VitaPacketPreamble) *sdrobjects.SdrMeterPacket {
	index := 0
	var meterPacket sdrobjects.SdrMeterPacket

	numberOfMeters := (len(data) - preamble.Header.Payload_cutoff_bytes) / 4

	for i := 0; i < numberOfMeters; i++ {

		meterPacket.Ids = append(meterPacket.Ids, binary.BigEndian.Uint16(data[index:index+2]))
		index += 2
		buf := bytes.NewBuffer(data[index : index+2])
		var valueRes int16
		binary.Read(buf, binary.BigEndian, &valueRes)
		index += 2
		meterPacket.Vals = append(meterPacket.Vals, valueRes)
	}

	return &meterPacket

}

func ParseVitaWaterfall(data []byte, preamble *VitaPacketPreamble) *sdrobjects.SdrWaterfallTile {
	index := 0
	var wftile sdrobjects.SdrWaterfallTile

	wftile.FirstPixelFreq = binary.BigEndian.Uint64(data[index:8]) >> 20
	index += 8

	wftile.BinBandwidth = binary.BigEndian.Uint64(data[index:index+8]) >> 20
	index += 8

	wftile.LineDurationMS = binary.BigEndian.Uint32(data[index : index+4])
	index += 4

	wftile.Width = binary.BigEndian.Uint16(data[index : index+2])
	index += 2

	wftile.Height = binary.BigEndian.Uint16(data[index : index+2])
	index += 2

	wftile.Timecode = binary.BigEndian.Uint32(data[index : index+4])
	index += 4

	wftile.AutoBlackLevel = binary.BigEndian.Uint32(data[index : index+4])
	index += 4

	for i := 0; i < (len(data))-preamble.Header.Payload_cutoff_bytes-index; i += 2 {
		wftile.Data = append(wftile.Data, binary.BigEndian.Uint16(data[i+index:i+index+2]))
	}

	return &wftile
}

func ParseVitaOpus(data []byte, preamble *VitaPacketPreamble) []byte {
	return data[:len(data)-preamble.Header.Payload_cutoff_bytes]
}

func ParseFData(data []byte, preamble *VitaPacketPreamble) *sdrobjects.SdrIfData {

	payload := data[:len(data)-preamble.Header.Payload_cutoff_bytes]

	var res sdrobjects.SdrIfData
	res.Stream_id = preamble.Stream_id

	switch preamble.Class_id.PacketClassCode { // dax audio
	case SL_VITA_IF_NARROW_CLASS:
		res.Data = payload
		return &res
	}

	for i := 0; i <= (len(payload)-4)/4; i++ {
		fVal := getFloat32fromLE(data[i*4:(i*4)+4]) * float32(ONE_OVER_ZERO_DBFS)
		res.Data = append(res.Data, getBytesFromFloat(fVal)...)
	}
	return &res
}

func ParseDiscoveryPackage(data []byte, preamble *VitaPacketPreamble) string {
	return string(data[:len(data)-4])
}

func getFloat32fromLE(bytes []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(bytes))
}

func getBytesFromFloat(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}
