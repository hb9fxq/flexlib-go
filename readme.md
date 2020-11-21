# flexlib-go
Go lib to interact with flexradio 6k series

## Installation

1. Install GO https://golang.org/doc/install

2. Install CMDs
<pre>
go get -u github.com/hb9fxq/flexlib-go/cmd/smartsdr-iqtransfer
go get -u github.com/hb9fxq/flexlib-go/cmd/smartsdr-mqttadapter
</pre>


## Tools

### iq-transfer
Tool to transfer DAX IQ data from a FRS 6K Radio on any platform. Get the latest file for your platform from https://github.com/hb9fxq/flexlib-go/releases

Run SmartSDR for Windows on any machine, while pulling the IQ data at any other computer in the network.
When you run iq-transfer, make sure, that you select the matching DAX IQ channel in Smartsdr for windows. (See option CH in options below) (Yes, it must run to use IQ-transfer ...for now)

_Options:_
* **RADIO** IP address of the radio
* **MYUDP** UDP port on local machine that the radio will send VITA49 traffic. Must be a free port on your machine. Check your firewall! 
* **CH** DAX-IQ channel to stream.
* **FWD** Endpoint to send the Float32 IQ data to. If not supplied, the data is written to stdout and can be used for piping. You can find a sample for GNU Radio under https://github.com/hb9fxq/flexlib-go/tree/master/GRC/iq-transfer
* **RATE** SampleRate in kHz, Possible Values: 24 48 96 192

__e.g.__

send data 127.0.0.1:2345 **./iq-transfer_linux64  --RADIO=192.168.92.8 --MYUDP=7799 --CH=1 --RATE=192 --FWD=127.0.0.1:2345**
 
record IQ Data to a file **./iq-transfer_linux64  --RADIO=192.168.92.8 --MYUDP=7799 --CH=1 --RATE=192 > "$(date +"%FT%T").raw"**

![alt text](https://github.com/hb9fxq/flexlib-go/raw/master/assets/GRC/iq-transfer/iq_transfer_fft.png "FFT with GRC using iq-transfer util")

![alt text](https://github.com/hb9fxq/flexlib-go/raw/master/assets/GRC/iq-transfer/2017-10-07_20_15_54-SmartSDR-Win.png "DAX IQ setting in SmartSDR")

#### DAX IQ Data to OpenWebRx ### 

Use the following device command in config_webrx.py
```
start_rtl_command="/home/f102/wrk/flexlib-go/bin/iq-transfer_linux64 --RADIO=192.168.92.8 --MYUDP=7799 --CH=1 --RATE=192"
format_conversion=""
```

![alt text](https://raw.githubusercontent.com/hb9fxq/flexlib-go/master/assets/GRC/iq-transfer/openwebrx.png "DAX IQ to OpenWebRX")


## Status
Started with VITA49 parsing, see the https://github.com/hb9fxq/flexlib-go/blob/master/src/flexlib-go/pcap_test.go file for VITA49 handling. 

Currently working on a [pcap](https://github.com/hb9fxq/flexlib-go/tree/master/test_input) file, captured with tcpdump to not permanently stress the radio.

Reconstructed waterfall tile data from pcap:

![alt text](https://raw.githubusercontent.com/hb9fxq/flexlib-go/master/assets/test_output/waterfall.png "waterfall from pcap")

Reconstruced opus audio from pcap: 

https://soundcloud.com/frank-werner-hb9fxq-14069568/opus-decoded

Reconstructed FFT Plot (all captured values added)

![alt text](https://github.com/hb9fxq/flexlib-go/raw/master/assets/test_output/fft.png "fft from pcap")
