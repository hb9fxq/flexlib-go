# flexlib-go
Multi platform tools to interact with flexradio 6k series radio.

Currently, any tools, except the MQTT adapter require an instance of SmartSDR Windows/OSX/IOS to be running. DAX and DAX IQ Data is currently not "headless".

## Installation
Option A) Binary Download
* Download the latest binary release from https://github.com/hb9fxq/flexlib-go/releases

Option B) Install from source

* Install GO https://golang.org/doc/install
* Install CMDs
<pre>
go get -u github.com/hb9fxq/flexlib-go/cmd/smartsdr-iqtransfer
go get -u github.com/hb9fxq/flexlib-go/cmd/smartsdr-mqttadapter
go get -u github.com/hb9fxq/flexlib-go/cmd/smartsdr-daxclient
</pre>

## Tools

### Binary "smartsdr-iqtransfer"
Tool to transfer DAX IQ data from a FRS 6K Radio to any platform.

When you run smartsdr-iqtransfer, make sure, that you select the matching DAX IQ channel in Smartsdr.

fCenter is the center of the GUI Panadapter

_Options:_
* **RADIO** IP address of the radio
* **MYUDP** UDP port on local machine that the radio will send VITA49 traffic. Must be a free port on your machine. Check your firewall! 
* **CH** DAX-IQ channel to stream.
* **FWD** Endpoint to send the Float32 IQ data to. If not supplied, the data is written to stdout and can be used for piping. You can find a sample for GNU Radio under https://github.com/hb9fxq/flexlib-go/tree/master/GRC/iq-transfer
* **RATE** SampleRate in kHz, Possible Values: 24000 48000 96000 192000

__e.g.__

Send raw IQ data 127.0.0.1:2345 <pre>./smartsdr-iqtransfer --RADIO=192.168.92.8 --MYUDP=5999 --RATE=192000 --CH=1 --FWD=127.0.0.1:2345</pre>
 
record IQ Data to a file 
<pre>./smartsdr-iqtransfer  --RADIO=192.168.92.8 --MYUDP=5999 --RATE=192000 --CH=1 --FWD=127.0.0.1:2345 > "$(date +"%FT%T").raw"</pre>

![alt text](https://github.com/hb9fxq/flexlib-go/raw/master/assets/grc_sample.png "FFT with GRC using iq-transfer util")

### Binary "smartsdr-daxclient"

Receives RAW DAX audio streams (RX Channels 1-6)

_Options:_
* **RADIO** IP address of the radio
* **MYUDP** UDP port on local machine that the radio will send VITA49 traffic. Must be a free port on your machine. Check your firewall! 
* **CH** DAX audio channel to stream.
* **FWD** Endpoint to send the Float32 IQ data to. If not supplied, the data is written to stdout and can be used for piping. You can find a sample for GNU Radio under https://github.com/hb9fxq/flexlib-go/tree/master/GRC/iq-transfer

e.g.
Forward raw DAX audio stream from channel 1 to a computer on the network (FWD)
<pre>./smartsdr-daxclient --RADIO=192.168.92.8 --MYUDP=5999 --CH=1 --FWD=127.0.0.1:2345</pre>

Play RAW audio to the speaker. 2 Channels, 32 Bit float, big endian
<pre>socat -u udp-recv:2345 - | play -q -t f32 -r 24k --endian big -c 2 -</pre>

![alt text](https://github.com/hb9fxq/flexlib-go/raw/master/assets/wsjtx_use_case.png "Pulling DAX Audio to WSJX-T on Ubuntu")

### Binary "smartsdr-mqttadapter"

Tool to reflect most important radio status, like Slices, Panadapters and connected clients to a MQTT broker. Useful for status monitoring, dashboards or advanced radio integration.

<pre>./smartsdr-mqttadapter --RADIO=192.168.92.8 --MQTTBROKER=tcp://192.168.92.7:1883 --MQTTCLIENTID=flexdev --MQTTTOPIC=flexdev</pre>


![alt text](https://github.com/hb9fxq/flexlib-go/raw/master/assets/mqtt_sample.png "DAX IQ setting in SmartSDR")

## Some experiments
The library is currently able to parse most of the VITA 49 types, that the FRS is using... 

Reconstructed waterfall tile data from pcap:

![alt text](https://raw.githubusercontent.com/hb9fxq/flexlib-go/master/assets/test_output/waterfall.png "waterfall from pcap")

Reconstructed opus audio from pcap: 

https://soundcloud.com/frank-werner-hb9fxq-14069568/opus-decoded

Reconstructed FFT Plot (all captured fft points aggregated)

![alt text](https://github.com/hb9fxq/flexlib-go/raw/master/assets/test_output/fft.png "fft from pcap")
