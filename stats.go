// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0

package vlc

// #include "glue.h"
import "C"

// Supplies a set of statistics on playing media.
type Stats struct {
	ptr *C.libvlc_media_stats_t
}

// ReadBytes returns the amount of bytes read from the input source.
func (this *Stats) ReadBytes() int { return int(this.ptr.i_read_bytes) }

// InputBitRate returns the input transfer rate.
func (this *Stats) InputBitRate() float32 { return float32(this.ptr.f_input_bitrate) }

func (this *Stats) DemuxReadBytes() int     { return int(this.ptr.i_demux_read_bytes) }
func (this *Stats) DemuxBitRate() float32   { return float32(this.ptr.f_demux_bitrate) }
func (this *Stats) DemuxCorrupted() int     { return int(this.ptr.i_demux_corrupted) }
func (this *Stats) DemuxDiscontinuity() int { return int(this.ptr.i_demux_discontinuity) }
func (this *Stats) DecodedVideo() int       { return int(this.ptr.i_decoded_video) }
func (this *Stats) DecodedAudio() int       { return int(this.ptr.i_decoded_audio) }

// DisplayedPictures returns the amount of displayed pictures.
func (this *Stats) DisplayedPictures() int { return int(this.ptr.i_displayed_pictures) }

// LostPictures returns the amount of lost pictures.
func (this *Stats) LostPictures() int { return int(this.ptr.i_lost_pictures) }

// PlayedAudioBuffers returns the amount of played audio buffers.
func (this *Stats) PlayedAudioBuffers() int { return int(this.ptr.i_played_abuffers) }

// LostAudioBuffers returns the amount of lost audio buffers.
func (this *Stats) LostAudioBuffers() int { return int(this.ptr.i_lost_abuffers) }

// SentPackets returns the amount of packets sent.
func (this *Stats) SentPackets() int { return int(this.ptr.i_sent_packets) }

// SentBytes returns the amount of bytes sent.
func (this *Stats) SentBytes() int { return int(this.ptr.i_sent_bytes) }

// SendBitRate returns the transfer bitrate.
func (this *Stats) SendBitRate() float32 { return float32(this.ptr.f_send_bitrate) }
