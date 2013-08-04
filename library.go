// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0

package vlc

// #include "glue.h"
import "C"

// A media library
type Library struct {
	ptr *C.libvlc_media_library_t
}

// Retain increments the reference count of the instance.
func (this *Library) Retain() (err error) {
	if this.ptr == nil {
		return &VLCError{"Library is nil"}
	}
	C.libvlc_media_library_retain(this.ptr)
	return
}

// Release decreases the reference count of the instance and destroys it when it reaches zero.
func (this *Library) Release() (err error) {
	if this.ptr == nil {
		return &VLCError{"Library is nil"}
	}

	C.libvlc_media_library_release(this.ptr)
	return
}

// Load loads the library contents.
func (this *Library) Load() error {
	if this.ptr == nil {
		return &VLCError{"Library is nil"}
	}
	C.libvlc_media_library_load(this.ptr)
	return checkError()
}

// Items returns a list of all the media items in this library.
func (this *Library) Items() (*MediaList, error) {
	if this.ptr == nil {
		return nil, &VLCError{"Library is nil"}
	}

	if c := C.libvlc_media_library_media_list(this.ptr); c != nil {
		return &MediaList{c}, nil
	}

	return nil, checkError()
}
