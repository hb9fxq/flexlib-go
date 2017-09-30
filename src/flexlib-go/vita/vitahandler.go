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

import (
	"../sdrobjects"
	"encoding/binary"
	"errors"
)

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

	if(header.T){
		index += 4
	}

	header.payload_bytes = (index *-1) + header.payload_bytes

	return nil, &vitaPacketPreamble, data[index:]
}

func ParseVitaWaterfall(data []byte, preamble *VitaPacketPreamble) (*sdrobjects.SdrWaterfallTile) {
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


	for i := 0; i < (len(data)-index); i+=2 {
		wftile.Data = append(wftile.Data, binary.BigEndian.Uint16(data[i : i +2]))
	}

	return &wftile
}



func ParseVitaOpus(data []byte, preamble *VitaPacketPreamble) ([]byte) {
	return data[:]
}
