// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0

package vlc

// #include "glue.h"
import "C"
import (
	"os"
)

// Maintains a list of Media items.
type MediaList struct {
	ptr *C.libvlc_media_list_t
}

// Retain increments the reference count of this MediaList instance.
func (this *MediaList) Retain() os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_list_retain(this.ptr)
	return checkError()
}

// Release cleans up any memory used by this list and decrements the
// reference counter for the Media instance this came from.
func (this *MediaList) Release() os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_list_release(this.ptr)
	this.ptr = nil
	return checkError()
}

// Set associates a media instance with this media list.
// If another media instance was present it will be released.
//
// Note: MediaList.Lock() should NOT be held upon entering this function.
func (this *MediaList) Set(m *Media) os.Error {
	if this.ptr == nil || m.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_list_set_media(this.ptr, m.ptr)
	return checkError()
}

// Get returns a media instance from this list.
// This action will increase the reference count on the media instance.
//
// Note: MediaList.Lock() should NOT be held upon entering this function.
func (this *MediaList) Get() (*Media, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_media_list_media(this.ptr); c != nil {
		return &Media{c}, nil
	}

	return nil, checkError()
}

// Add adds a media instance to this list.
//
// Note: MediaList.Lock() SHOULD be held upon entering this function.
func (this *MediaList) Add(m *Media) os.Error {
	if this.ptr == nil || m.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_list_add_media(this.ptr, m.ptr)
	return checkError()
}

// Insert adds a media instance to the list at the given position.
//
// Note: MediaList.Lock() SHOULD be held upon entering this function.
func (this *MediaList) Insert(m *Media, pos int) os.Error {
	if this.ptr == nil || m.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_list_insert_media(this.ptr, m.ptr, C.int(pos))
	return checkError()
}

// Remove removes a media instance at the given position from the list.
//
// Note: MediaList.Lock() SHOULD be held upon entering this function.
func (this *MediaList) Remove(pos int) os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_list_remove_index(this.ptr, C.int(pos))
	return checkError()
}

// Count returns the number if items in the list.
//
// Note: MediaList.Lock() SHOULD be held upon entering this function.
func (this *MediaList) Count() (int, os.Error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}

	return int(C.libvlc_media_list_count(this.ptr)), checkError()
}

// At returns the media at the given list position.
// This action will increase the reference count on the media instance.
//
// Note: MediaList.Lock() SHOULD be held upon entering this function.
func (this *MediaList) At(pos int) (*Media, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_media_list_item_at_index(this.ptr, C.int(pos)); c != nil {
		return &Media{c}, nil
	}

	return nil, checkError()
}

// Index returns the position of the given media in the list.
//
// Note: MediaList.Lock() SHOULD be held upon entering this function.
func (this *MediaList) Index(m *Media) (int, os.Error) {
	if this.ptr == nil || m.ptr == nil {
		return 0, os.EINVAL
	}

	return int(C.libvlc_media_list_index_of_item(this.ptr, m.ptr)), checkError()
}

// IsReadOnly returns true if this list is readonly for a user.
func (this *MediaList) IsReadOnly() (bool, os.Error) {
	if this.ptr == nil {
		return false, os.EINVAL
	}
	return C.libvlc_media_list_is_readonly(this.ptr) == 0, checkError()
}

// Lock gets a lock on the list items.
func (this *MediaList) Lock() os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}
	C.libvlc_media_list_lock(this.ptr)
	return checkError()
}

// Unlock removes a lock on the list items.
func (this *MediaList) Unlock() os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}
	C.libvlc_media_list_unlock(this.ptr)
	return checkError()
}

// Events returns an Eventmanager for this list.
func (this *MediaList) Events() (*EventManager, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_media_list_event_manager(this.ptr); c != nil {
		return NewEventManager(c), nil
	}

	return nil, checkError()
}
