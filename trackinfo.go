// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0

package vlc

// #include "glue.h"
import "C"
import (
	"bytes"
	"encoding/binary"
)

// Represents a single media track. Can be audio or video.
// Access Audio() or Video() depending on the value of Type.
type TrackInfo struct {
	ptr *C.libvlc_media_track_info_t
}

func (this *TrackInfo) Codec() uint32                  { return uint32(this.ptr.i_codec) }
func (this *TrackInfo) Id() int                        { return int(this.ptr.i_id) }
func (this *TrackInfo) Type() TrackType                { return TrackType(this.ptr.i_type) }
func (this *TrackInfo) Profile() int                   { return int(this.ptr.i_profile) }
func (this *TrackInfo) Level() int                     { return int(this.ptr.i_level) }
func (this *TrackInfo) Audio() (channels, rate uint32) { return this.readU32_2() }
func (this *TrackInfo) Video() (width, height uint32)  { return this.readU32_2() }

func (this *TrackInfo) readU32_2() (uint32, uint32) {
	var a, b uint32
	buf := bytes.NewBuffer(this.ptr.u[:])
	binary.Read(buf, binary.LittleEndian, &a)
	binary.Read(buf, binary.LittleEndian, &b)
	return a, b
}

// Description for video, audio tracks and subtitles.
type TrackDescription struct {
	ptr *C.libvlc_track_description_t
}

// Release releases memory for this instance.
func (this *TrackDescription) Release() {
	if this.ptr != nil {
		C.libvlc_track_description_release(this.ptr)
		this.ptr = nil
	}
}

// Id returns the track Id.
func (this *TrackDescription) Id() int { return int(this.ptr.i_id) }

// Name returns the track name.
func (this *TrackDescription) Name() string { return C.GoString(this.ptr.psz_name) }


// List of track descriptions.
type TrackDescriptionList []*TrackDescription

func (this *TrackDescriptionList) fromC(p *C.libvlc_track_description_t) {
	for ; p != nil; p = (*C.libvlc_track_description_t)(p.p_next) {
		*this = append(*this, &TrackDescription{p})
	}
}

// Release recursively releases memory for this list.
func (this *TrackDescriptionList) Release() {
	if len(*this) == 0 {
		return
	}

	(*this)[0].Release()
	*this = nil
}
