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

// A single libvlc instance.
type Instance struct {
	ptr *C.libvlc_instance_t
}

// New creates and initializes a new VLC instance with the given parameters.
// Returns nil and a possible error if no instance could be created.
func New(argv []string) (i *Instance, err os.Error) {
	cstr := make([]*C.char, len(argv))

	for i := range cstr {
		cstr[i] = C.CString(argv[i])
	}

	if c := C.libvlc_new(C.int(len(argv)), *(***C.char)(unsafe.Pointer(&cstr))); c != nil {
		i = &Instance{c}
	} else {
		err = checkError()
	}

	for i := range cstr {
		C.free(unsafe.Pointer(cstr[i]))
	}
	return
}

// Retain increments the reference count of the Instance.
// The initial reference count is 1 after vlc.New() returns.
func (this *Instance) Retain() (err os.Error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_retain(this.ptr)
	return
}

// Release decreases the reference count of the instance and destroys it
// when it reaches zero.
func (this *Instance) Release() (err os.Error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_release(this.ptr)
	return
}

// StartUI tries to start a user interface for the Instance.
// Specify an empty name to use the default.
func (this *Instance) StartUI(name string) (err os.Error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	c := C.CString(name)
	defer C.free(unsafe.Pointer(c))

	if C.libvlc_add_intf(this.ptr, c) < 0 {
		err = checkError()
	}

	return
}

// SetName sets the human-readable application name (e.g. "FooBar player 1.2.3")
// and user-agent name (e.g. "FooBar/1.2.3 Python/2.6.0"). LibVLC passes this as
// the user agent string when a protocol requires it.
func (this *Instance) SetName(appname, httpname string) (err os.Error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	ca := C.CString(appname)
	ch := C.CString(httpname)

	C.libvlc_set_user_agent(this.ptr, ca, ch)

	C.free(unsafe.Pointer(ca))
	C.free(unsafe.Pointer(ch))
	return
}

// Wait waits until an interface causes the instance to exit.
// You should start at least one interface first, using Instance.StartUI().
func (this *Instance) Wait() os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_wait(this.ptr)
	return nil
}

// LogVerbosity returns the VLC messaging verbosity level.
func (this *Instance) LogVerbosity() uint {
	if this.ptr == nil {
		return 0
	}
	return uint(C.libvlc_get_log_verbosity(this.ptr))
}

// SetLogVerbosity sets the VLC messaging verbosity level.
func (this *Instance) SetLogVerbosity(v uint) {
	if this.ptr != nil {
		C.libvlc_set_log_verbosity(this.ptr, C.uint(v))
	}
}

// OpenLog opens a VLC message log instance.
func (this *Instance) OpenLog() (*Log, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_log_open(this.ptr); c != nil {
		l := new(Log)
		l.fromC(c)
		return l, nil
	}

	return nil, checkError()
}

// OpenMediaUri loads a media instance from the given uri.
func (this *Instance) OpenMediaUri(uri string) (*Media, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	c := C.CString(uri)
	defer C.free(unsafe.Pointer(c))

	if m := C.libvlc_media_new_location(this.ptr, c); m != nil {
		return &Media{m}, nil
	}

	return nil, checkError()
}

// OpenMediaFile loads a media instance from the given filesystem path.
func (this *Instance) OpenMediaFile(path string) (*Media, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	c := C.CString(path)
	defer C.free(unsafe.Pointer(c))

	if m := C.libvlc_media_new_path(this.ptr, c); m != nil {
		return &Media{m}, nil
	}

	return nil, checkError()
}

// OpenMediaFd creates a media instance for an open file descriptor.
// The file descriptor shall be open for reading (or reading and writing).
//
// Regular file descriptors, pipe read descriptors and character device
// descriptors (including TTYs) are supported on all platforms.
// Block device descriptors are supported where available.
// Directory descriptors are supported on systems that provide fdopendir().
// Sockets are supported on all platforms where they are file descriptors,
// i.e. all except Windows.
//
// Note: This library will /not/ automatically close the file descriptor
// under any circumstance. Nevertheless, a file descriptor can usually only be
// rendered once in a media player. To render it a second time, the file
// descriptor should probably be rewound to the beginning with lseek().
func (this *Instance) OpenMediaFd(fd int) (*Media, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if m := C.libvlc_media_new_fd(this.ptr, C.int(fd)); m != nil {
		return &Media{m}, nil
	}

	return nil, checkError()
}

// OpenMediaNode creates a media instance as an empty node with a given name.
func (this *Instance) OpenMediaNode(name string) (*Media, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	c := C.CString(name)
	defer C.free(unsafe.Pointer(c))

	if m := C.libvlc_media_new_as_node(this.ptr, c); m != nil {
		return &Media{m}, nil
	}

	return nil, checkError()
}

// NewPlayer creates an empty media player object.
func (this *Instance) NewPlayer() (*Player, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_media_player_new(this.ptr); c != nil {
		return &Player{c}, nil
	}

	return nil, checkError()
}

// NewList creates and initializes a new media list.
func (this *Instance) NewList() (*MediaList, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_media_list_new(this.ptr); c != nil {
		return &MediaList{c}, nil
	}

	return nil, checkError()
}

// NewListPlayer creates an empty media list player object.
func (this *Instance) NewListPlayer() (*ListPlayer, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_media_list_player_new(this.ptr); c != nil {
		return &ListPlayer{c}, nil
	}

	return nil, checkError()
}

// NewLibrary creates an empty media library.
func (this *Instance) NewLibrary() (*Library, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_media_library_new(this.ptr); c != nil {
		return &Library{c}, nil
	}

	return nil, checkError()
}

// Discoverer creates a new discover media service by name.
func (this *Instance) Discoverer(name string) (*Discoverer, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))

	if c := C.libvlc_media_discoverer_new_from_name(this.ptr, s); c != nil {
		return &Discoverer{c}, nil
	}

	return nil, checkError()
}

// VlmRelease releases the vlm instance associated with this instance.
func (this *Instance) VlmRelease() os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_vlm_release(this.ptr)
	return checkError()
}

// VlmAddBroadcast adds a broadcast with given input.
//
//     name: The name of the new broadcast.
//     input: The input MRL.
//     output: The output MRL (the parameter to the "sout" variable).
//     options: Additional options.
//     enabled: Enable the new broadcast?
//     loop: Should this broadcast be played in loop?
//
func (this *Instance) VlmAddBroadcast(name, input, output string, options []string, enabled, loop bool) os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	a := C.CString(name)
	b := C.CString(input)
	c := C.CString(output)
	d := C.int(len(options))
	e := make([]*C.char, len(options))

	for i := range e {
		e[i] = C.CString(options[i])
	}

	var f, g C.int
	if enabled {
		f = 1
	}

	if loop {
		g = 1
	}

	C.libvlc_vlm_add_broadcast(this.ptr, a, b, c, d, *(***C.char)(unsafe.Pointer(&e)), f, g)

	C.free(unsafe.Pointer(a))
	C.free(unsafe.Pointer(b))
	C.free(unsafe.Pointer(c))

	for i := range e {
		C.free(unsafe.Pointer(e[i]))
	}

	return checkError()
}

// VlmAddVOD adds a VOD with given input.
//
//     name: The name of the new broadcast.
//     input: The input MRL.
//     options: Additional options.
//     mux: The muxer of the vod media.
//     enabled: Enable the new broadcast?
//
func (this *Instance) VlmAddVOD(name, input, output, mux string, options []string, enabled bool) os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	a := C.CString(name)
	b := C.CString(input)
	c := C.int(len(options))
	d := make([]*C.char, len(options))

	for i := range d {
		d[i] = C.CString(options[i])
	}

	var e C.int
	if enabled {
		e = 1
	}

	f := C.CString(mux)

	C.libvlc_vlm_add_vod(this.ptr, a, b, c, *(***C.char)(unsafe.Pointer(&d)), e, f)

	C.free(unsafe.Pointer(a))
	C.free(unsafe.Pointer(b))

	for i := range d {
		C.free(unsafe.Pointer(d[i]))
	}

	C.free(unsafe.Pointer(f))
	return checkError()
}

// VlmDelete deletes the given media (VOD or broadcast).
func (this *Instance) VlmDelete(name string) os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	c := C.CString(name)
	C.libvlc_vlm_del_media(this.ptr, c)
	C.free(unsafe.Pointer(c))

	return checkError()
}

// VlmSetEnabled enables or disables the given media (VOD or broadcast).
func (this *Instance) VlmSetEnabled(name string, toggle bool) os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	c := C.CString(name)
	if toggle {
		C.libvlc_vlm_set_enabled(this.ptr, c, 1)
	} else {
		C.libvlc_vlm_set_enabled(this.ptr, c, 0)
	}
	C.free(unsafe.Pointer(c))

	return checkError()
}

// VlmSetLoop enables or disables the given media's loop state (VOD or broadcast).
func (this *Instance) VlmSetLoop(name string, toggle bool) os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	c := C.CString(name)
	if toggle {
		C.libvlc_vlm_set_loop(this.ptr, c, 1)
	} else {
		C.libvlc_vlm_set_loop(this.ptr, c, 0)
	}
	C.free(unsafe.Pointer(c))

	return checkError()
}

// VlmSetOutput sets the output for the given media (VOD or broadcast).
func (this *Instance) VlmSetOutput(name, output string) os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	a := C.CString(name)
	b := C.CString(output)
	C.libvlc_vlm_set_output(this.ptr, a, b)
	C.free(unsafe.Pointer(a))
	C.free(unsafe.Pointer(b))
	return checkError()
}

// VlmSetInput sets the input MRL for the given media (VOD or broadcast).
// This will delete all existing inputs and add the specified one.
func (this *Instance) VlmSetInput(name, input string) os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	a := C.CString(name)
	b := C.CString(input)
	C.libvlc_vlm_set_input(this.ptr, a, b)
	C.free(unsafe.Pointer(a))
	C.free(unsafe.Pointer(b))
	return checkError()
}

// VlmAddInput adds an input MRL for the given media (VOD or broadcast).
func (this *Instance) VlmAddInput(name, input string) os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	a := C.CString(name)
	b := C.CString(input)
	C.libvlc_vlm_add_input(this.ptr, a, b)
	C.free(unsafe.Pointer(a))
	C.free(unsafe.Pointer(b))
	return checkError()
}

// VlmSetMux sets a media's VOD muxer.
func (this *Instance) VlmSetMux(name, mux string) os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	a := C.CString(name)
	b := C.CString(mux)
	C.libvlc_vlm_set_mux(this.ptr, a, b)
	C.free(unsafe.Pointer(a))
	C.free(unsafe.Pointer(b))
	return checkError()
}

// VlmChangeMedia edits the parameters of a media. This will delete all existing
// inputs and add the specified one.
//
//     name: The name of the new broadcast.
//     input: The input MRL.
//     output: The output MRL (the parameter to the "sout" variable).
//     options: Additional options.
//     enabled: Enable the new broadcast?
//     loop: Should this broadcast be played in loop?
//
func (this *Instance) VlmChangeMedia(name, input, output string, options []string, enabled, loop bool) os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	a := C.CString(name)
	b := C.CString(input)
	c := C.CString(output)
	d := C.int(len(options))
	e := make([]*C.char, len(options))

	for i := range e {
		e[i] = C.CString(options[i])
	}

	var f, g C.int
	if enabled {
		f = 1
	}

	if loop {
		g = 1
	}

	C.libvlc_vlm_change_media(this.ptr, a, b, c, d, *(***C.char)(unsafe.Pointer(&e)), f, g)

	C.free(unsafe.Pointer(a))
	C.free(unsafe.Pointer(b))
	C.free(unsafe.Pointer(c))

	for i := range e {
		C.free(unsafe.Pointer(e[i]))
	}

	return checkError()
}

// VlmPlay plays the named broadcast.
func (this *Instance) VlmPlay(name string) os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	a := C.CString(name)
	C.libvlc_vlm_play_media(this.ptr, a)
	C.free(unsafe.Pointer(a))
	return checkError()
}

// VlmStop halts playback of the named broadcast.
func (this *Instance) VlmStop(name string) os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	a := C.CString(name)
	C.libvlc_vlm_stop_media(this.ptr, a)
	C.free(unsafe.Pointer(a))
	return checkError()
}

// VlmPause pauses playback of the named broadcast.
func (this *Instance) VlmPause(name string) os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	a := C.CString(name)
	C.libvlc_vlm_pause_media(this.ptr, a)
	C.free(unsafe.Pointer(a))
	return checkError()
}

// VlmSeek seeks in the named broadcast.
func (this *Instance) VlmSeek(name string, percentage float32) os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	a := C.CString(name)
	C.libvlc_vlm_seek_media(this.ptr, a, C.float(percentage))
	C.free(unsafe.Pointer(a))
	return checkError()
}

// VlmMediaInfo returns information about the named media as a JSON string.
//
// Note: This function is mainly intended for debugging use,
func (this *Instance) VlmMediaInfo(name string) (s string, err os.Error) {
	if this.ptr == nil {
		return "", os.EINVAL
	}

	a := C.CString(name)

	if c := C.libvlc_vlm_show_media(this.ptr, a); c != nil {
		s = C.GoString(c)
		C.free(unsafe.Pointer(c))
	} else {
		err = checkError()
	}

	C.free(unsafe.Pointer(a))
	return
}

// VlmPosition returns the instance position by name or instance id.
func (this *Instance) VlmPosition(name string, id int) (float32, os.Error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}

	a := C.CString(name)
	defer C.free(unsafe.Pointer(a))
	return float32(C.libvlc_vlm_get_media_instance_position(this.ptr, a, C.int(id))), checkError()
}

// VlmTime returns the instance time by name or instance id.
func (this *Instance) VlmTime(name string, id int) (int, os.Error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}

	a := C.CString(name)
	defer C.free(unsafe.Pointer(a))
	return int(C.libvlc_vlm_get_media_instance_time(this.ptr, a, C.int(id))), checkError()
}

// VlmLength returns the instance length by name or instance id.
func (this *Instance) VlmLength(name string, id int) (int, os.Error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}

	a := C.CString(name)
	defer C.free(unsafe.Pointer(a))
	return int(C.libvlc_vlm_get_media_instance_length(this.ptr, a, C.int(id))), checkError()
}

// VlmRate returns the instance playback rate by name or instance id.
func (this *Instance) VlmRate(name string, id int) (int, os.Error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}

	a := C.CString(name)
	defer C.free(unsafe.Pointer(a))
	return int(C.libvlc_vlm_get_media_instance_rate(this.ptr, a, C.int(id))), checkError()
}

// VlmEvents returns an event manager for a VLM instance.
func (this *Instance) VlmEvents() (*EventManager, os.Error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_vlm_get_event_manager(this.ptr); c != nil {
		return NewEventManager(c), nil
	}

	return nil, checkError()
}
