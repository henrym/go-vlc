// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0

// Go bindings for libVLC 1.1.9.
package vlc

// #cgo        LDFLAGS: -lvlc
// #cgo  linux  CFLAGS: -I/usr/local/include
// #cgo  linux LDFLAGS: -L/usr/local/lib
// #cgo darwin  CFLAGS: -I/usr/local/include
// #cgo darwin LDFLAGS: -L/usr/local/lib
// #include "glue.h"
import "C"
import (
	"os"
	"unsafe"
)

// libVLC version numbers.
const (
	VersionMajor    = 1
	VersionMinor    = 1
	VersionRevision = 9
	VersionExtra    = 0

	// Version as a single integer. Practical for version comparison.
	Version = (VersionMajor << 24) | (VersionMinor << 16) | (VersionRevision << 8) | VersionExtra
)

// Version returns the libVLC version as a human-readable string.
func VersionString() string { return C.GoString(C.libvlc_get_version()) }

func (this EventType) String() string {
	return C.GoString(C.libvlc_event_type_name(C.libvlc_event_type_t(this)))
}

// Clears the LibVLC error status for the current thread. This is optional.
// By default, the error status is automatically overriden when a new error
// occurs, and destroyed when the thread exits.
func ClearError() { C.libvlc_clearerr() }

// Compiler returns the compiler used to build libvlc.
func Compiler() string { return C.GoString(C.libvlc_get_compiler()) }

// ChangeSet returns the change set for the libvlc build.
func ChangeSet() string { return C.GoString(C.libvlc_get_changeset()) }

// checkError checks if there is a new error message available. If so, return
// it as an os.Error. For internal use only.
func checkError() (err os.Error) {
	if c := C.libvlc_errmsg(); c != nil {
		err = os.NewError(C.GoString(c))
		C.free(unsafe.Pointer(c))
	}
	return
}
