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
	"testing"
	"github.com/krippendorf/flexlib-go/obj"
	"fmt"
	"time"
	"strconv"
	"github.com/krippendorf/flexlib-go/sdrobjects"
)

func TestRadioInitIntegration(t *testing.T) {

	ctx := new(obj.RadioContext)
	ctx.RadioAddr = "192.168.92.8"
	ctx.MyUdpEndpointPort = "4700"
	ctx.ChannelRadioResponse = make(chan string)

	go obj.InitRadioContext(ctx)


	go func(ctx *obj.RadioContext) {
		for {
			fmt.Println(">" +  <-ctx.ChannelRadioResponse + "<")
		}
	}(ctx)

	for{
		if(len(ctx.RadioHandle)>0){
			break
		}
		time.Sleep(500)
	}


	forever := make(chan bool)
	forever <- true
}


func TestRadioSubIqInitIntegration(t *testing.T) {
	ctx := new(obj.RadioContext)
	ctx.RadioAddr = "192.168.92.8"
	ctx.MyUdpEndpointPort = "4700"
	ctx.ChannelRadioResponse = make(chan string)
	ctx.ChannelVitaIfData = make(chan *sdrobjects.SdrIfData)

	go obj.InitRadioContext(ctx)


	go func(ctx *obj.RadioContext) {
		for {
			fmt.Println(">" +  <-ctx.ChannelRadioResponse + "<")
		}
	}(ctx)

	go func(ctx *obj.RadioContext) {

		cnt := 0

		for {
			<-ctx.ChannelVitaIfData
			cnt++
			if(cnt%500 == 0){
				fmt.Println("VitaIfDataPacket count: " + strconv.Itoa(cnt))
			}
		}
	}(ctx)

	for{
		if(len(ctx.RadioHandle)>0){
			break
		}
		time.Sleep(500)
	}

	obj.SendRadioCommand(ctx, "stream create daxiq=1ip=" + ctx.MyUdpEndpointIP.String() + " port=" + ctx.MyUdpEndpointPort)

	forever := make(chan bool)
	forever <- true
}