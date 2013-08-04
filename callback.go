// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0

package vlc

// #include "glue.h"
import "C"
import (
	"unsafe"
)

// Used when hooking/unhooking events.
type eventData struct {
	t C.libvlc_event_type_t
	f EventHandler
	d interface{}
}

// Event callback handler.
type EventHandler func(evt *Event, userdata interface{})

//export goEventCB
func goEventCB(e *C.libvlc_event_t, userdata unsafe.Pointer) {
	evt := &Event{Type: EventType(e._type)}
	evt.b.Write(e.u[:])

	if rd := (*eventData)(userdata); rd.f != nil {
		rd.f(evt, rd.d)
	}
}

// Used in Player.SetCallbacks() to render video to a custom memory location.
type memRenderReq struct {
	lh LockHandler
	uh UnlockHandler
	dh DisplayHandler
	ud interface{}
}

// Whenever a new video frame needs to be decoded, the lock callback is
// invoked. Depending on the video chroma, one or three pixel planes of
// adequate dimensions must be returned. Those planes must be aligned on
// 32-bytes boundaries.
//
// void* (*lock) (void** plane, void* userdata)
type LockHandler func(plane uintptr, userdata interface{}) uintptr

//export goLockCB
func goLockCB(userdata, plane unsafe.Pointer) unsafe.Pointer {
	if req := (*memRenderReq)(userdata); req.lh != nil {
		return unsafe.Pointer(req.lh(uintptr(plane), req.ud))
	}
	return nil
}

// When the video frame is decoded, the unlock callback is invoked. The
// second parameter to the callback is the return value of the lock callback.
// The third parameter conveys the pixel planes for convenience.
//
// void (*unlock) (void* picture, void* const* plane, void* userdata)
type UnlockHandler func(picture, plane uintptr, userdata interface{})

//export goUnlockCB
func goUnlockCB(userdata, picture, plane unsafe.Pointer) {
	if req := (*memRenderReq)(userdata); req.uh != nil {
		req.uh(uintptr(picture), uintptr(plane), req.ud)
	}
}

type DisplayHandler func(picture uintptr, userdata interface{})

//export goDisplayCB
//
// void (*display) (void* picture, void* userdata)
func goDisplayCB(userdata, picture unsafe.Pointer) {
	if req := (*memRenderReq)(userdata); req.dh != nil {
		req.dh(uintptr(picture), req.ud)
	}
}
