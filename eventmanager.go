// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0

package vlc

// #include "glue.h"
//
// extern void goEventCB(const struct libvlc_event_t*, void*);
//
// static int goAttach(libvlc_event_manager_t* em, libvlc_event_type_t et, void* userdata) {
//    return libvlc_event_attach(em, et, goEventCB, userdata);
// }
// static void goDetach(libvlc_event_manager_t* em, libvlc_event_type_t et, void* userdata) {
//    libvlc_event_detach(em, et, goEventCB, userdata);
// }
import "C"
import (
	"sync"
	"unsafe"
)

// A libvlc instance has an event manager which can be used to hook event callbacks,
type EventManager struct {
	ptr    *C.libvlc_event_manager_t
	events map[int]*eventData
	m      *sync.Mutex
}

func NewEventManager(p *C.libvlc_event_manager_t) *EventManager {
	v := new(EventManager)
	v.ptr = p
	v.m = new(sync.Mutex)
	v.events = make(map[int]*eventData)
	return v
}

// Attach registers the given event handler and returns a unique id
// we can use to detach the event at a later point.
func (this *EventManager) Attach(et EventType, cb EventHandler, userdata interface{}) (id int, err error) {
	if this.ptr == nil {
		return 0, &VLCError{"EventManager is nil"}
	}

	id = this.getUniqId()

	this.m.Lock()
	this.events[id] = &eventData{C.libvlc_event_type_t(et), cb, userdata}
	this.m.Unlock()

	if C.goAttach(this.ptr, this.events[id].t, unsafe.Pointer(this.events[id])) != 0 {
		err = checkError()
	}

	return
}

// Detach unregisters the given event id.
func (this *EventManager) Detach(id int) (err error) {
	if this.ptr == nil {
		return &VLCError{"EventManager is nil"}
	}

	var ed *eventData
	var ok bool

	this.m.Lock()
	if ed, ok = this.events[id]; !ok {
		this.m.Unlock()
		return &VLCError{"No event with that id"}
	}

	delete(this.events, id)
	this.m.Unlock()

	C.goDetach(this.ptr, ed.t, unsafe.Pointer(ed))
	return
}

// getUniqId finds and returns a unique event id.
func (this *EventManager) getUniqId() int {
	var id int
	var ok bool

	this.m.Lock()
	defer this.m.Unlock()

	for {
		if _, ok = this.events[id]; !ok {
			break
		}
		id++
	}

	return id
}
