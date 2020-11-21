# flexlib-go
Go lib to interact with flexradio 6k series

## Installation

1. Install GO https://golang.org/doc/install

2. Install CMDs
<pre>
go get -u github.com/hb9fxq/flexlib-go/cmd/smartsdr-iqtransfer
go get -u github.com/hb9fxq/flexlib-go/cmd/smartsdr-mqttadapter
go get -u github.com/hb9fxq/flexlib-go/cmd/smartsdr-daxclient
</pre>


## Tools

### smartsdr-iqtransfer
Tool to transfer DAX IQ data from a FRS 6K Radio on any platform. Get the latest file for your platform from https://github.com/hb9fxq/flexlib-go/releases

Run SmartSDR for Windows, Mac or IOS on any machine, while pulling the IQ data at any other computer in the network.

When you run iq-transfer, make sure, that you select the matching DAX IQ channel in Smartsdr. (See option CH in options below) (Yes, it must run to use IQ-transfer ...for now)

_Options:_
* **RADIO** IP address of the radio
* **MYUDP** UDP port on local machine that the radio will send VITA49 traffic. Must be a free port on your machine. Check your firewall! 
* **CH** DAX-IQ channel to stream.
* **FWD** Endpoint to send the Float32 IQ data to. If not supplied, the data is written to stdout and can be used for piping. You can find a sample for GNU Radio under https://github.com/hb9fxq/flexlib-go/tree/master/GRC/iq-transfer
* **RATE** SampleRate in kHz, Possible Values: 24000 48000 96000 192000

__e.g.__

send raw IQ data 127.0.0.1:2345 **./smartsdr-iqtransfer  --RADIO=192.168.92.8 --MYUDP=5999 --RATE=192000 --CH=1 --FWD=127.0.0.1:2345**
 
record IQ Data to a file **./smartsdr-iqtransfer  --RADIO=192.168.92.8 --MYUDP=5999 --RATE=192000 --CH=1 --FWD=127.0.0.1:2345 > "$(date +"%FT%T").raw"**

![alt text](https://github.com/hb9fxq/flexlib-go/raw/master/assets/grc_sample.png "FFT with GRC using iq-transfer util")


### smartsdr-daxclient

Receives RAW DAX audio streams (RX Channels 1-6)

_Options:_
* **RADIO** IP address of the radio
* **MYUDP** UDP port on local machine that the radio will send VITA49 traffic. Must be a free port on your machine. Check your firewall! 
* **CH** DAX audio channel to stream.
* **FWD** Endpoint to send the Float32 IQ data to. If not supplied, the data is written to stdout and can be used for piping. You can find a sample for GNU Radio under https://github.com/hb9fxq/flexlib-go/tree/master/GRC/iq-transfer

e.g.
**./smartsdr-daxclient --RADIO=192.168.92.8 --MYUDP=5999 --CH=1 --FWD=127.0.0.1:2345**

### smartsdr-daxclient

Tool to reflect most important radio status, like Slices, Panadapters and connected clients to an MQTT broker. Useful for status monitoring, dashboards or advanced radio integration

**./smartsdr-mqttadapter --RADIO=192.168.92.8 --MQTTBROKER=tcp://192.168.92.7:1883 --MQTTCLIENTID=flexdev --MQTTTOPIC=flexdev**

![alt text](https://github.com/hb9fxq/flexlib-go/raw/master/assets/mqtt_sample.png "DAX IQ setting in SmartSDR")



## Experiments
The library is currently able to parse most of the VITA 49 types, that the FRS is using... 

Some experiments:

Reconstructed waterfall tile data from pcap:

![alt text](https://raw.githubusercontent.com/hb9fxq/flexlib-go/master/assets/test_output/waterfall.png "waterfall from pcap")

Reconstruced opus audio from pcap: 

https://soundcloud.com/frank-werner-hb9fxq-14069568/opus-decoded

Reconstructed FFT Plot (all captured fft points aggregated)

![alt text](https://github.com/hb9fxq/flexlib-go/raw/master/assets/test_output/fft.png "fft from pcap")
