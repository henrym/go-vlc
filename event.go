// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0

package vlc

// #include "glue.h"
import "C"
import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

// Generic event type. Use a switch on Event.Type to determine which method
// to call for its data. The method names are the same as the event type.
// It's a bit of an odd way to handle this, but since event data comes in the
// form of a C union, there is some binary magic that has to be performed based
// on the event type.
type Event struct {
	Type EventType
	b    bytes.Buffer
}

func (this *Event) MediaMetaChanged() MetaProperty                    { return MetaProperty(this.readU32()) }
func (this *Event) MediaSubItemAdded() *Media                         { return this.readMedia() }
func (this *Event) MediaDurationChanged() int64                       { return this.readI64() }
func (this *Event) MediaParsedChanged() int                           { return int(this.readI32()) }
func (this *Event) MediaFreed() *Media                                { return this.readMedia() }
func (this *Event) MediaStateChanged() MediaState                     { return MediaState(this.readU32()) }
func (this *Event) MediaPlayerTimeChanged() int64                     { return this.readI64() }
func (this *Event) MediaPlayerPositionChanged() float32               { return this.readF32() }
func (this *Event) MediaPlayerSeekableChanged() bool                  { return this.readB() }
func (this *Event) MediaPlayerTitleChanged() bool                     { return this.readB() }
func (this *Event) MediaPlayerLengthChanged() bool                    { return this.readB() }
func (this *Event) MediaPlayerSnapshotTaken() string                  { return this.readS() }
func (this *Event) MediaListItemAdded() *Media                        { return this.readMedia() }
func (this *Event) MediaListWillAddItem() *Media                      { return this.readMedia() }
func (this *Event) MediaListItemDeleted() *Media                      { return this.readMedia() }
func (this *Event) MediaListWillDeleteItem() *Media                   { return this.readMedia() }
func (this *Event) MediaListPlayerNextItemSet() *Media                { return this.readMedia() }
func (this *Event) VlmMediaAdded() (string, string)                   { return this.readS2() }
func (this *Event) VlmMediaRemoved() (string, string)                 { return this.readS2() }
func (this *Event) VlmMediaChanged() (string, string)                 { return this.readS2() }
func (this *Event) VlmMediaInstanceStarted() (string, string)         { return this.readS2() }
func (this *Event) VlmMediaInstanceStopped() (string, string)         { return this.readS2() }
func (this *Event) VlmMediaInstanceStatusInitAdded() (string, string) { return this.readS2() }
func (this *Event) VlmMediaInstanceStatusOpening() (string, string)   { return this.readS2() }
func (this *Event) VlmMediaInstanceStatusPlaying() (string, string)   { return this.readS2() }
func (this *Event) VlmMediaInstanceStatusPause() (string, string)     { return this.readS2() }
func (this *Event) VlmMediaInstanceStatusEnd() (string, string)       { return this.readS2() }
func (this *Event) VlmMediaInstanceStatusError() (string, string)     { return this.readS2() }
func (this *Event) MediaPlayerMediaChanged() *Media                   { return this.readMedia() }
func (this *Event) MediaPlayerNothingSpecial() MediaState             { return MediaState(this.readU32()) }
func (this *Event) MediaPlayerOpening() MediaState                    { return MediaState(this.readU32()) }
func (this *Event) MediaPlayerBuffering() MediaState                  { return MediaState(this.readU32()) }
func (this *Event) MediaPlayerPlaying() MediaState                    { return MediaState(this.readU32()) }
func (this *Event) MediaPlayerPaused() MediaState                     { return MediaState(this.readU32()) }
func (this *Event) MediaPlayerStopped() MediaState                    { return MediaState(this.readU32()) }
func (this *Event) MediaPlayerForward() MediaState                    { return MediaState(this.readU32()) }
func (this *Event) MediaPlayerBackward() MediaState                   { return MediaState(this.readU32()) }
func (this *Event) MediaPlayerEndReached() MediaState                 { return MediaState(this.readU32()) }
func (this *Event) MediaPlayerEncounteredError() MediaState           { return MediaState(this.readU32()) }
func (this *Event) MediaListPlayerPlayed() MediaState                 { return MediaState(this.readU32()) }
func (this *Event) MediaListPlayerStopped() MediaState                { return MediaState(this.readU32()) }
func (this *Event) MediaListViewItemAdded() *Media                    { return this.readMedia() }
func (this *Event) MediaListViewWillAddItem() *Media                  { return this.readMedia() }
func (this *Event) MediaListViewItemDeleted() *Media                  { return this.readMedia() }
func (this *Event) MediaListViewWillDeleteItem() *Media               { return this.readMedia() }
func (this *Event) MediaDiscovererStarted() *Discoverer               { return this.readDiscoverer() }
func (this *Event) MediaDiscovererEnded() *Discoverer                 { return this.readDiscoverer() }

func (this *Event) readI8() int8 {
	var i int8
	binary.Read(&this.b, binary.LittleEndian, &i)
	return i
}

func (this *Event) readI16() int16 {
	var i int16
	binary.Read(&this.b, binary.LittleEndian, &i)
	return i
}

func (this *Event) readI32() int32 {
	var i int32
	binary.Read(&this.b, binary.LittleEndian, &i)
	return i
}

func (this *Event) readI64() int64 {
	var i int64
	binary.Read(&this.b, binary.LittleEndian, &i)
	return i
}

func (this *Event) readU8() uint8 {
	var i uint8
	binary.Read(&this.b, binary.LittleEndian, &i)
	return i
}
func (this *Event) readU16() uint16 {
	var i uint16
	binary.Read(&this.b, binary.LittleEndian, &i)
	return i
}

func (this *Event) readU32() uint32 {
	var i uint32
	binary.Read(&this.b, binary.LittleEndian, &i)
	return i
}

func (this *Event) readU64() uint64 {
	var i uint64
	binary.Read(&this.b, binary.LittleEndian, &i)
	return i
}

func (this *Event) readF32() float32 {
	var i float32
	binary.Read(&this.b, binary.LittleEndian, &i)
	return i
}

func (this *Event) readF64() float64 {
	var i float64
	binary.Read(&this.b, binary.LittleEndian, &i)
	return i
}

func (this *Event) readS() string {
	var i uint64
	binary.Read(&this.b, binary.LittleEndian, &i)
	up := unsafe.Pointer(uintptr(i))
	s := C.GoString((*C.char)(up))
	C.free(up)
	return s
}

func (this *Event) readS2() (string, string) {
	var a, b uint64

	binary.Read(&this.b, binary.LittleEndian, &a)
	binary.Read(&this.b, binary.LittleEndian, &b)

	pa := unsafe.Pointer(uintptr(a))
	pb := unsafe.Pointer(uintptr(b))

	sa := C.GoString((*C.char)(pa))
	sb := C.GoString((*C.char)(pb))

	C.free(pa)
	C.free(pb)
	return sa, sb
}

func (this *Event) readB() bool {
	var i int64
	binary.Read(&this.b, binary.LittleEndian, &i)
	return i == 0
}

func (this *Event) readMedia() *Media {
	var i uint64
	binary.Read(&this.b, binary.LittleEndian, &i)
	return &Media{(*C.libvlc_media_t)(unsafe.Pointer(uintptr(i)))}
}

func (this *Event) readDiscoverer() *Discoverer {
	var i uint64
	binary.Read(&this.b, binary.LittleEndian, &i)
	return &Discoverer{(*C.libvlc_media_discoverer_t)(unsafe.Pointer(uintptr(i)))}
}
