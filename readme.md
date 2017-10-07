# flexlib-go
Go lib to interact with flexradio 6k series

## Status
Started with VITA49 parsing, see the https://github.com/krippendorf/flexlib-go/blob/master/src/flexlib-go/pcap_test.go file for VITA49 handling. 

Currently working on a [pcap](https://github.com/krippendorf/flexlib-go/tree/master/test_input) file, captured with tcpdump to not permanently stress the radio.

Reconstructed waterfall tile data from pcap:

![alt text](https://github.com/krippendorf/flexlib-go/raw/master/test_output/waterfall.png "waterfall from pcap")

Reconstruced opus audio from pcap: 

https://soundcloud.com/frank-werner-krippendorf-14069568/opus-decoded

Reconstructed FFT Plot (all captured values added)

![alt text](https://github.com/krippendorf/flexlib-go/raw/master/test_output/fft.png "fft from pcap")


## Tools

### iq-transfer
Tool to transfer DAX IQ data from a FRS 6K Radio on any platform. Get the latest file for your platform from https://github.com/krippendorf/flexlib-go/releases

_Options:_
* **RADIO** IP address of the radio
* **MYUDP** UDP port on local machine that the radio will send VITA49 traffic. Must be a free port on your machine. Check your firewall! 
* **CH** DAX-IQ channel to stream.
* **FWD** Endpoint to send the Float32 IQ data to. You can find a sample for GNU Radio under https://github.com/krippendorf/flexlib-go/tree/master/GRC/iq-transfer

__e.g.__

**./iq-transfer_linux64  --RADIO=192.168.92.8 --MYUDP=7799 --CH=1 --RATE=192 --FWD=127.0.0.1:2345**

![alt text](https://github.com/krippendorf/flexlib-go/raw/master/GRC/iq-transfer/iq_transfer_fft.png "FFT with GRC using iq-transfer util")


