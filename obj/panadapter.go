package obj

type Panadapter struct {
	Id           string
	Center       int32
	ClientHandle string
}

type IqStream struct {
	Id   int
	Pan  string
	Rate int
}

/*
F:SCDD36271|display pan
0x40000000
client_handle=0x736ABCCB
wnb=0
wnb_level=0
wnb_updating=1
band_zoom=0
segment_zoom=0
x_pixels=1047
y_pixels=510
center=18.119837
bandwidth=0.162466
min_dbm=-137.66
max_dbm=-42.66
fps=12
average=50
weighted_average=0
rfgain=0
rxant=ANT1 wide=0
loopa=0
loopb=0
band=17
daxiq_channel=0
waterfall=0x42000000
min_bw=0.004920
max_bw=14.745601
xvtr=
pre=
ant_list=ANT1,ANT2,RX_A,XVTR


*/
