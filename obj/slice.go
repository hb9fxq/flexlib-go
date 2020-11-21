package obj

type Slice struct {
	Id           string
	InUse        bool
	RfFrequency  float64
	ClientHandle string
	IndexLetter  string
	Mode         string
	TxAnt        string
	RxAnt        string
	Panadapter   string
	Dax          string
}

/*
impl:
in_use=1
RF_frequency=18.100100
index_letter=A
client_handle=0x736ABCCB
rxant=ANT1
mode=USB
txant=ANT1
*/

/*
F:SCDD36271|
	slice 0
in_use=1
RF_frequency=18.100100
client_handle=0x736ABCCB
index_letter=A
rit_on=0
rit_freq=0
xit_on=0
xit_freq=0
rxant=ANT1
mode=USB
wide=0
filter_lo=100
filter_hi=2800
step=100
step_list=1,10,50,100,500,1000,2000,3000
agc_mode=med
agc_threshold=65
agc_off_level=10
pan=0x40000000
txant=ANT1
loopa=0
loopb=0
qsk=0
dax=0
dax_clients=0
lock=0
tx=1
active=0
audio_level=50
audio_pan=50
audio_mute=0
record=0
play=disabled
record_time=0.0
anf=0
anf_level=0
nr=0
nr_level=0
nb=0
nb_level=50
wnb=0
wnb_level=0
apf=0
apf_level=0
squelch=1
squelch_level=20
diversity=0
diversity_parent=0
diversity_child=0
diversity_index=1342177293
ant_list=ANT1,ANT2,RX_A,XVTR
mode_list=LSB,USB,AM,CW,DIGL,DIGU,SAM,FM,NFM,DFM,RTTY
fm_tone_mode=OFF
fm_tone_value=67.0
fm_repeater_offset_freq=0.000000
tx_offset_freq=0.000000
repeater_offset_dir=SIMPLEX
fm_tone_burst=0
fm_deviation=5000
dfm_pre_de_emphasis=0
post_demod_low=300
post_demod_high=3300
rtty_mark=2125
rtty_shift=170
digl_offset=2210
digu_offset=1500
post_demod_bypass=0
rfgain=0
tx_ant_list=ANT1,ANT2,XVTR
*/
