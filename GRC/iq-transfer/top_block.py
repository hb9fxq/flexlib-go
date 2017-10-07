#!/usr/bin/env python2
# -*- coding: utf-8 -*-
##################################################
# GNU Radio Python Flow Graph
# Title: Top Block
# Generated: Sat Oct  7 20:29:38 2017
##################################################

if __name__ == '__main__':
    import ctypes
    import sys
    if sys.platform.startswith('linux'):
        try:
            x11 = ctypes.cdll.LoadLibrary('libX11.so')
            x11.XInitThreads()
        except:
            print "Warning: failed to XInitThreads()"

from gnuradio import blocks
from gnuradio import eng_notation
from gnuradio import gr
from gnuradio import wxgui
from gnuradio.eng_option import eng_option
from gnuradio.fft import window
from gnuradio.filter import firdes
from gnuradio.wxgui import fftsink2
from gnuradio.wxgui import waterfallsink2
from grc_gnuradio import wxgui as grc_wxgui
from optparse import OptionParser
import wx


class top_block(grc_wxgui.top_block_gui):

    def __init__(self):
        grc_wxgui.top_block_gui.__init__(self, title="Top Block")
        _icon_path = "/usr/share/icons/hicolor/32x32/apps/gnuradio-grc.png"
        self.SetIcon(wx.Icon(_icon_path, wx.BITMAP_TYPE_ANY))

        ##################################################
        # Variables
        ##################################################
        self.samp_rate = samp_rate = 192e3

        ##################################################
        # Blocks
        ##################################################
        self.wxgui_waterfallsink2_0 = waterfallsink2.waterfall_sink_c(
        	self.GetWin(),
        	baseband_freq=0,
        	dynamic_range=100,
        	ref_level=0,
        	ref_scale=2.0,
        	sample_rate=samp_rate,
        	fft_size=4096,
        	fft_rate=20,
        	average=False,
        	avg_alpha=None,
        	title="Flexradio IQ Stream 192kHz",
        )
        self.Add(self.wxgui_waterfallsink2_0.win)
        self.wxgui_fftsink2_0 = fftsink2.fft_sink_c(
        	self.GetWin(),
        	baseband_freq=0,
        	y_per_div=10,
        	y_divs=10,
        	ref_level=0,
        	ref_scale=2.0,
        	sample_rate=samp_rate,
        	fft_size=4096,
        	fft_rate=15,
        	average=True,
        	avg_alpha=None,
        	title="FFT Plot Flexradio IQ Stream 192kHz",
        	peak_hold=False,
        	win=window.hamming,
        )
        self.Add(self.wxgui_fftsink2_0.win)
        self.blocks_udp_source_0 = blocks.udp_source(gr.sizeof_float*1, "127.0.0.1", 2345, 4096, True)
        self.blocks_float_to_complex_0 = blocks.float_to_complex(1)
        self.blocks_deinterleave_0 = blocks.deinterleave(gr.sizeof_float*1, 1)

        ##################################################
        # Connections
        ##################################################
        self.connect((self.blocks_deinterleave_0, 0), (self.blocks_float_to_complex_0, 0))    
        self.connect((self.blocks_deinterleave_0, 1), (self.blocks_float_to_complex_0, 1))    
        self.connect((self.blocks_float_to_complex_0, 0), (self.wxgui_fftsink2_0, 0))    
        self.connect((self.blocks_float_to_complex_0, 0), (self.wxgui_waterfallsink2_0, 0))    
        self.connect((self.blocks_udp_source_0, 0), (self.blocks_deinterleave_0, 0))    

    def get_samp_rate(self):
        return self.samp_rate

    def set_samp_rate(self, samp_rate):
        self.samp_rate = samp_rate
        self.wxgui_fftsink2_0.set_sample_rate(self.samp_rate)
        self.wxgui_waterfallsink2_0.set_sample_rate(self.samp_rate)


def main(top_block_cls=top_block, options=None):

    tb = top_block_cls()
    tb.Start(True)
    tb.Wait()


if __name__ == '__main__':
    main()
