# flexlib-go
Go lib to interact with flexradio 6k series

## Status
Started with VITA49 parsing, see the https://github.com/krippendorf/flexlib-go/blob/master/src/flexlib-go/pcap_test.go file for VITA49 handling. 

Currently working on a [pcap](https://github.com/krippendorf/flexlib-go/tree/master/test_input) file, captured with tcpdump to not permanently stress the radio.

Reconstructed Waterfall data from pcap: 

![alt text](https://github.com/krippendorf/flexlib-go/raw/master/test_output/waterfall.png "waterfall from pcap")

Reconstruced opus audio from pcap: 

https://soundcloud.com/frank-werner-krippendorf-14069568/opus-decoded
