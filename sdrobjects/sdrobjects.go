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
package sdrobjects

type VitaPacketType int

type SdrWaterfallTile struct {
	FirstPixelFreq uint64
	BinBandwidth   uint64
	LineDurationMS uint32
	Width          uint16
	Height         uint16
	Timecode       uint32
	AutoBlackLevel uint32
	Data           []uint16
}

type SdrFFTPacket struct {
	StartBin_index uint32
	NumBins        uint32
	BinSize        uint32
	FrameIndex     uint32
	Payload        []uint16
}

type SdrMeterPacket struct {
	Ids  []uint16
	Vals []int16
}

type SdrIfData struct{
	Stream_id   uint32
	Data []byte
}


