// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0

package vlc

// #include "glue.h"
import "C"
import (
	"os"
)

// This player is meant for playlist playback.
// This is basically a wrapper for vlc.Player that takes care of playlist rotation.
type ListPlayer struct {
	ptr *C.libvlc_media_list_player_t
}

// Release decreases the reference count of the instance and destroys it when it reaches zero.
func (this *ListPlayer) Release() (err os.Error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_list_player_release(this.ptr)
	return
}

// Events returns an Eventmanager for this player.
func (this *ListPlayer) Events() (*EventManager, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_media_list_player_event_manager(this.ptr); c != nil {
		return NewEventManager(c), nil
	}

	return nil, checkError()
}

// Replace replaces the Player instance in this listplayer with a new one.
func (this *ListPlayer) Replace(p *Player) os.Error {
	if this.ptr == nil || p.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_list_player_set_media_player(this.ptr, p.ptr)
	return checkError()
}

// Set sets the MediaList associated with this player.
func (this *ListPlayer) Set(l *MediaList) os.Error {
	if this.ptr == nil || l.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_list_player_set_media_list(this.ptr, l.ptr)
	return checkError()
}

// Play plays the entries in the media list.
func (this *ListPlayer) Play() os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_list_player_play(this.ptr)
	return checkError()
}

// Pause pauses playback.
func (this *ListPlayer) Pause() os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_list_player_pause(this.ptr)
	return checkError()
}

// IsPlaying returns true if the player is currently playing.
func (this *ListPlayer) IsPlaying() (bool, os.Error) {
	if this.ptr == nil {
		return false, os.EINVAL
	}
	return C.libvlc_media_list_player_is_playing(this.ptr) != 0, checkError()
}

// State returns the current media state.
func (this *ListPlayer) State() (MediaState, os.Error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return MediaState(C.libvlc_media_list_player_get_state(this.ptr)), checkError()
}

// PlayAt plays the entry at the given list index.
func (this *ListPlayer) PlayAt(pos int) os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_list_player_play_item_at_index(this.ptr, C.int(pos))
	return checkError()
}

// PlayItem plays the given entry.
//
// Note: The supplied Media must be part of this list.
func (this *ListPlayer) PlayItem(m *Media) os.Error {
	if this.ptr == nil || m.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_list_player_play_item(this.ptr, m.ptr)
	return checkError()
}

// Stop halts playback.
func (this *ListPlayer) Stop() os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_list_player_stop(this.ptr)
	return checkError()
}

// Next plays the next item in the list if applicable.
func (this *ListPlayer) Next() os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_list_player_next(this.ptr)
	return checkError()
}

// Prev plays the previous item in the list if applicable.
func (this *ListPlayer) Prev() os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_list_player_previous(this.ptr)
	return checkError()
}

// SetMode sets the current playback mode.
// Any of: PMDefault, PMLoop or PMRepeat.
func (this *ListPlayer) SetMode(pm PlaybackMode) os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_list_player_set_playback_mode(this.ptr, C.libvlc_playback_mode_t(pm))
	return checkError()
}
