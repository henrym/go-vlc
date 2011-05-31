// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0

package vlc

// #include "glue.h"
import "C"
import (
	"os"
	"unsafe"
)

// Medis discovery service.
type Discoverer struct {
	ptr *C.libvlc_media_discoverer_t
}

// Release media discover object
func (this *Discoverer) Release() {
	if this.ptr != nil {
		C.libvlc_media_discoverer_release(this.ptr)
		this.ptr = nil
	}
}

// LocalizedName return the localzied discovery service name.
func (this *Discoverer) LocalizedName() (s string, err os.Error) {
	if this.ptr != nil {
		return "", os.EINVAL
	}

	if c := C.libvlc_media_discoverer_localized_name(this.ptr); c != nil {
		s = C.GoString(c)
		C.free(unsafe.Pointer(c))
		return
	}

	return "", checkError()
}

// MediaList returns a list of media items.
func (this *Discoverer) MediaList() (m *MediaList, err os.Error) {
	if this.ptr != nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_media_discoverer_media_list(this.ptr); c != nil {
		return &MediaList{c}, nil
	}

	return nil, checkError()
}

// Events returns an event manager for this instance.
// Note: This method does not increment the media reference count.
func (this *Discoverer) Events() (*EventManager, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_media_discoverer_event_manager(this.ptr); c != nil {
		return NewEventManager(c), nil
	}

	return nil, checkError()
}

// IsRunning returns true if the discovery service is currently running.
func (this *Discoverer) IsRunning() (bool, os.Error) {
	if this.ptr == nil {
		return false, os.EINVAL
	}
	return C.libvlc_media_discoverer_is_running(this.ptr) != 0, checkError()
}

// Description for audio output.
type AudioOutput struct {
	ptr *C.libvlc_audio_output_t
}

// Name returns the track name.
func (this *AudioOutput) Name() string { return C.GoString(this.ptr.psz_name) }

// Description returns the track description.
func (this *AudioOutput) Description() string { return C.GoString(this.ptr.psz_description) }

func (this *AudioOutput) Release() {
	if this.ptr != nil {
		C.libvlc_audio_output_list_release(this.ptr)
		this.ptr = nil
	}
}

// List of track descriptions.
type AudioOutputList []*AudioOutput

func (this *AudioOutputList) fromC(p *C.libvlc_audio_output_t) {
	for ; p != nil; p = (*C.libvlc_audio_output_t)(p.p_next) {
		*this = append(*this, &AudioOutput{p})
	}
}

func (this *AudioOutputList) Release() {
	if len(*this) == 0 {
		return
	}

	(*this)[0].Release()
	*this = nil
}

// Rectangle type for video geometry.
type Rect struct {
	ptr *C.libvlc_rectangle_t
}

// Top returns the top coordinate.
func (this *Rect) Top() int { return int(this.ptr.top) }

// Left returns the left coordinate.
func (this *Rect) Left() int { return int(this.ptr.left) }

// Bottom returns the bottom coordinate.
func (this *Rect) Bottom() int { return int(this.ptr.bottom) }

// Right returns the right coordinate.
func (this *Rect) Right() int { return int(this.ptr.right) }
