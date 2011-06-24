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

type Media struct {
	ptr *C.libvlc_media_t
}

// Retain increments the reference count of this Media instance.
func (this *Media) Retain() (err os.Error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_retain(this.ptr)
	return
}

// Release decreases the reference count of this Media instance and destroys it
// when it reaches zero. It will send out a MediaFreed event to all listeners.
// If the media descriptor object has been released it should not be used again.
func (this *Media) Release() (err os.Error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_release(this.ptr)
	return
}

// Duplicate duplicates the media object.
func (this *Media) Duplicate() (*Media, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_media_duplicate(this.ptr); c != nil {
		return &Media{c}, nil
	}

	return nil, checkError()
}

// Add an option to the media.
//
// This option will be used to determine how the media player will read the
// media. This allows us to use VLC's advanced reading/streaming options on a
// per-media basis.
//
// The options are detailed in vlc --full-help, for instance "--sout-all"
func (this *Media) AddOption(options string) os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	c := C.CString(options)
	C.libvlc_media_add_option(this.ptr, c)
	C.free(unsafe.Pointer(c))

	return checkError()
}

// Add an option to the media with configurable flags.
//
// This option will be used to determine how the media player will read the
// media. This allows us to use VLC's advanced reading/streaming options on a
// per-media basis.
//
// The options are detailed in vlc --full-help, for instance "--sout-all"
func (this *Media) AddOptionFlag(options string, flags uint32) os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	c := C.CString(options)
	C.libvlc_media_add_option_flag(this.ptr, c, C.uint(flags))
	C.free(unsafe.Pointer(c))

	return checkError()
}

// Mrl returns the media resource locator (mrl) from a media descriptor object.
func (this *Media) Mrl() (s string) {
	if this.ptr == nil {
		return
	}

	if c := C.libvlc_media_get_mrl(this.ptr); c != nil {
		s = C.GoString(c)
		C.free(unsafe.Pointer(c))
	}

	return
}

// Meta reads the specified metadata property of the media.
//
// If the media has not yet been parsed this will return an empty string.
//
// This method automatically calls Media.ParseAsync(), so after calling
// it you may receive a MediaMetaChanged event. If you prefer a synchronous
// version, ensure that you call Media.Parse() before Media.Meta().
func (this *Media) Meta(mp MetaProperty) (s string) {
	if this.ptr == nil {
		return
	}

	if c := C.libvlc_media_get_meta(this.ptr, C.libvlc_meta_t(mp)); c != nil {
		s = C.GoString(c)
		C.free(unsafe.Pointer(c))
	}

	return
}

// SetMeta sets the metadata for this media instance.
// Note: This method does not save the metadata. Call Media.SaveMeta() for this purpose.
func (this *Media) SetMeta(mp MetaProperty, v string) {
	if this.ptr == nil {
		return
	}

	c := C.CString(v)
	C.libvlc_media_set_meta(this.ptr, C.libvlc_meta_t(mp), c)
	C.free(unsafe.Pointer(c))
}

// SaveMeta saves the previously changed metadata.
func (this *Media) SaveMeta() (err os.Error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	if C.libvlc_media_save_meta(this.ptr) != 0 {
		err = checkError()
	}

	return
}

// State returns the current media state.
func (this *Media) State() MediaState {
	if this.ptr == nil {
		return MSError
	}
	return MediaState(C.libvlc_media_get_state(this.ptr))
}

// Stats returns media statistics.
func (this *Media) Stats() (*Stats, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	var c C.libvlc_media_stats_t

	if C.libvlc_media_get_stats(this.ptr, &c) != 0 {
		return nil, checkError()
	}

	return &Stats{&c}, nil
}

// SubItems returns subitems of this media instance. This will increment
// the reference count of this media instance. Use MediaList.Release() to
// decrement the reference count.
func (this *Media) SubItems() (*MediaList, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_media_subitems(this.ptr); c != nil {
		return &MediaList{c}, nil
	}

	return nil, checkError()
}

// Events returns an event manager for this media instance.
// Note: This method does not increment the media reference count.
func (this *Media) Events() (*EventManager, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_media_event_manager(this.ptr); c != nil {
		return NewEventManager(c), nil
	}

	return nil, checkError()
}

// Duration returns the duration in milliseconds for the current media instance.
func (this *Media) Duration() int64 {
	if this.ptr == nil {
		return 0
	}
	return int64(C.libvlc_media_get_duration(this.ptr))
}

// Parse the current media source.
//
// This fetches (local) meta data and track information.
// The method is synchronous version of Media.ParseAsync().
func (this *Media) Parse() {
	if this.ptr != nil {
		C.libvlc_media_parse(this.ptr)
	}
}

// Parse the current media source.
//
// This fetches (local) meta data and track information.
// The method is the asynchronous version of Media.Parse()
//
// To determine when this routine finishes, you can listen for a MediaParsedChanged
// event. However if the media was already parsed you will not receive this event.
func (this *Media) ParseAsync() {
	if this.ptr != nil {
		C.libvlc_media_parse_async(this.ptr)
	}
}

// IsParsed returns true if the media's metadata has already been parsed.
func (this *Media) IsParsed() bool {
	if this.ptr != nil {
		return C.libvlc_media_is_parsed(this.ptr) != 0
	}
	return false
}

// UserData returns the media descriptor's user_data. user_data is specialized
// data accessed by the host application, VLC.framework uses it as a pointer to
// a native object that references a libvlc_media_t pointer.
//
//TODO(jimt): I have no idea what this comment means. Presumably its a roundabout
// way of saying that the data specified in here will survive roundtrips through
// event callback handlers. So you can pass it anything you need.
func (this *Media) UserData() interface{} {
	if this.ptr != nil {
		return C.libvlc_media_get_user_data(this.ptr)
	}
	return nil
}

// SetUserData sets the media descriptor's user_data. user_data is specialized
// data accessed by the host application, VLC.framework uses it as a pointer to
// a native object that references a libvlc_media_t pointer.
//
//TODO(jimt): I have no idea what this comment means. Presumably its a roundabout
// way of saying that the data specified in here will survive roundtrips through
// event callback handlers. So you can pass it anything you need.
func (this *Media) SetUserData(v interface{}) {
	if this.ptr != nil {
		C.libvlc_media_set_user_data(this.ptr, unsafe.Pointer(&v))
	}
}

// TrackInfo yields the media descriptor's elementary stream descriptions.
//
// Note: You need to play the media _one_ time with --sout="#description"
// Not doing this will result in an empty array, and doing it more than once
// will duplicate the entries in the array each time. Something like this:
//
//     player, _ := media.NewPlayer()
//     media.AddOption("sout=#description")
//     player.Play()
//     // ... wait until playing
//     player.Release()
//
// This is very likely to change in next release, and will be done at the
// parsing phase instead.
func (this *Media) TrackInfo() ([]*TrackInfo, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	var c *C.libvlc_media_track_info_t
	if size := C.libvlc_media_get_tracks_info(this.ptr, &c); size > 0 {
		list := make([]*TrackInfo, size)
		addr := uintptr(unsafe.Pointer(c))
		sz := int(unsafe.Sizeof(c))

		for i := range list {
			list[i] = &TrackInfo{(*C.libvlc_media_track_info_t)(unsafe.Pointer(addr + uintptr(i*sz)))}
		}

		return list, nil
	}

	return nil, checkError()
}

// NewPlayer a media player from this media instance.
// After creating the player, you can destroy this Media instance, unless you
// really need it for something. It is not necessary to perform actual playback.
func (this *Media) NewPlayer() (*Player, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_media_player_new_from_media(this.ptr); c != nil {
		return &Player{c}, nil
	}

	return nil, checkError()
}
