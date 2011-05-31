// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0

package vlc

// #include "glue.h"
import "C"
import (
	"os"
)

type Log struct {
	ptr  *C.libvlc_log_t
	iter *C.libvlc_log_iterator_t
}

func (this *Log) fromC(p *C.libvlc_log_t) {
	if p != nil {
		this.ptr = p
		this.iter = C.libvlc_log_get_iterator(p)
	}
}

// Close closes the log and releases resources associated with it.
func (this *Log) Close() os.Error {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_log_close(this.ptr)
	this.ptr = nil

	C.libvlc_log_iterator_free(this.iter)
	this.iter = nil
	return nil
}

// Count returns the number of messages in this log.
func (this *Log) Count() uint {
	if this.ptr == nil {
		return 0
	}
	return uint(C.libvlc_log_count(this.ptr))
}

// Clear removes all messages from the log.
// This should be called regularly to avoid clogging up the pipes.
func (this *Log) Clear() {
	if this.ptr != nil {
		C.libvlc_log_clear(this.ptr)
	}
}

// Reads and returns the next log message if one is available. nil otherwise.
func (this *Log) Next() (m *Message) {
	if C.libvlc_log_iterator_has_next(this.iter) != 0 {
		m = NewMessage()
		m.ptr = C.libvlc_log_iterator_next(this.iter, m.ptr)
	}
	return
}

// A single log message.
type Message struct {
	ptr *C.libvlc_log_message_t
}

// NewMessage initializes a new log message.
func NewMessage() *Message {
	var p C.libvlc_log_message_t
	p.sizeof_msg = 0
	p.i_severity = 0
	p.psz_type = nil
	p.psz_name = nil
	p.psz_header = nil
	p.psz_message = nil
	return &Message{&p}
}

// Priority yields the message severity: Info, Error, Warning, Debug.
func (this *Message) Priority() LogPriority { return LogPriority(this.ptr.i_severity) }

// SetPriority sets the message severity: Info, Error, Warning, Debug.
func (this *Message) SetPriority(v LogPriority) { this.ptr.i_severity = C.int(v) }

// Type returns the module type.
func (this *Message) Type() string { return C.GoString(this.ptr.psz_type) }

// SetType sets the module type.
func (this *Message) SetType(v string) { this.ptr.psz_type = C.CString(v) }

// Name returns the module name.
func (this *Message) Name() string { return C.GoString(this.ptr.psz_name) }

// SetName sets the module name.
func (this *Message) SetName(v string) { this.ptr.psz_name = C.CString(v) }

// Header returns the optional header.
func (this *Message) Header() string { return C.GoString(this.ptr.psz_header) }

// SetHeader sets the optional header.
func (this *Message) SetHeader(v string) { this.ptr.psz_header = C.CString(v) }

// Message returns the actual log message content.
func (this *Message) Message() string { return C.GoString(this.ptr.psz_message) }

// SetMessage sets the actual log message content.
func (this *Message) SetMessage(v string) {
	this.ptr.psz_message = C.CString(v)
	this.ptr.sizeof_msg = C.uint(len(v))
}
