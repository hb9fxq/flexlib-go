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
package main

import (
	"../flexlib-go/sdrobjects"
	"../flexlib-go/vita"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/lucasb-eyer/go-colorful"
	"gopkg.in/hraban/opus.v2"
	"image"
	"image/draw"
	"image/png"
	"os"
	"testing"
)

type GradientTable []struct {
	Col colorful.Color
	Pos float64
}

/* thx https://github.com/lucasb-eyer/go-colorful/blob/master/doc/gradientgen/gradientgen.go*/
func (self GradientTable) GetInterpolatedColorFor(t float64) colorful.Color {
	for i := 0; i < len(self)-1; i++ {
		c1 := self[i]
		c2 := self[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			t := (t - c1.Pos) / (c2.Pos - c1.Pos)
			return c1.Col.BlendHcl(c2.Col, t).Clamped()
		}
	}

	return self[len(self)-1].Col
}

func TestParsePcap(t *testing.T) {

	// package counters
	_countFFT := 0
	_countRXOpus := 0
	_countDAX := 0
	_countMeter := 0
	_countWaterfall := 0
	_countUnknown := 0
	_countIf := 0

	TCP_FRAGMENTATION_SIZE := 1514

	// waterfall render img canvas
	var img = image.NewRGBA(image.Rect(0, 0, 2460, 560*3))

	keypoints := GradientTable{
		{MustParseHex("#000000"), 0.0},
		{MustParseHex("#0000ff"), 0.15},
		{MustParseHex("#00FF00"), 0.30},
		{MustParseHex("#ffff00"), 0.45},
		{MustParseHex("#ff0000"), 0.60},
		{MustParseHex("#800080"), 0.75},
		{MustParseHex("#ffffff"), 1.0},
	}

	// opus stream test output
	f, err := os.Create("../../test_output/opus_decoded_float_32_LE_24000.raw")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// pcap input
	if handle, err := pcap.OpenOffline("../../test_input/flex.pcap"); err != nil {
		panic(err)
	} else {

		dec, err := opus.NewDecoder(24e3, 2)

		if err != nil {
			panic(err)
		}

		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

		var buff []byte
		previous_fragment := false

		for packet := range packetSource.Packets() {

			temp := packet.ApplicationLayer().Payload()
			packet.Dump()

			// reassemble fragmented packages
			if len(packet.Data()) == TCP_FRAGMENTATION_SIZE {

				offset := 0

				if !previous_fragment {
					buff = []byte{}
					offset = 8
				}

				buff = append(buff, temp[offset:]...)
				previous_fragment = true
				continue
			}

			if previous_fragment {
				buff = append(buff, temp...)
				previous_fragment = false

			} else {
				buff = temp
			}

			// parse preamble
			err, preamble, payload := vita.ParseVitaPreamble(buff)

			if err != nil || preamble.Class_id == nil {
				continue
			}

			switch preamble.Header.Pkt_type {

			case vita.ExtDataWithStream:

				switch preamble.Class_id.PacketClassCode {

				case vita.SL_VITA_FFT_CLASS:
					_countFFT++
					break
				case vita.SL_VITA_OPUS_CLASS:
					// decode opus and output raw PCM
					frameSizeMs := 10 // ms
					frameSize := 2 * frameSizeMs * 24e3 / 1000
					pcm := make([]float32, frameSize)
					dec.DecodeFloat32(payload, pcm)
					for sample := range pcm {
						f.Write(sdrobjects.Float32ToBytes(pcm[sample]))
					}
					_countRXOpus++
					break
				case vita.SL_VITA_IF_NARROW_CLASS:
					_countDAX++
					break
				case vita.SL_VITA_METER_CLASS:
					_countMeter++
					break
				case vita.SL_VITA_WATERFALL_CLASS:
					tile := vita.ParseVitaWaterfall(payload, preamble)
					renderAppend(_countWaterfall*3, tile, img, keypoints)
					_countWaterfall++
					break
				default:
					_countUnknown++
					break
				}

				break

			case vita.IFDataWithStream:
				switch preamble.Class_id.InformationClassCode {
				case vita.SL_VITA_IF_WIDE_CLASS_24kHz:
				case vita.SL_VITA_IF_WIDE_CLASS_48kHz:
				case vita.SL_VITA_IF_WIDE_CLASS_96kHz:
				case vita.SL_VITA_IF_WIDE_CLASS_192kHz:
					_countIf++
				}

				break
			}
		}

		fmt.Printf("_countFFT %d\n", _countFFT)
		fmt.Printf("_countRXOpus %d\n", _countRXOpus)
		fmt.Printf("_countDAX %d\n", _countDAX)
		fmt.Printf("_countMeter %d\n", _countMeter)
		fmt.Printf("_countWaterfall %d\n", _countWaterfall)
		fmt.Printf("_countUnknown %d\n", _countUnknown)
		fmt.Printf("_countIf %d\n", _countIf)

		f, _ := os.OpenFile("../../test_output/waterfall.png", os.O_WRONLY|os.O_CREATE, 0600)
		defer f.Close()
		png.Encode(f, img)

	}
}
func renderAppend(y int, tile *sdrobjects.SdrWaterfallTile, img *image.RGBA, keypoints GradientTable) {
	i := 0
	cBlackLevel := keypoints.GetInterpolatedColorFor(0.0)

	for value := range tile.Data {
		gain := 1.125
		pVal := (float64(tile.Data[value]))
		cv := (1.0 / (65535.0)) * (pVal * gain)
		c := cBlackLevel

		if (tile.Data[value] - uint16(tile.AutoBlackLevel)) >= 1 {
			c = keypoints.GetInterpolatedColorFor(cv)
		}

		draw.Draw(img, image.Rect(i, y, i+1, y+3), &image.Uniform{c}, image.ZP, draw.Src)
		i++
	}
}

func MustParseHex(s string) colorful.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		panic("MustParseHex: " + err.Error())
	}
	return c
}