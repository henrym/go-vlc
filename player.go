// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0

package vlc

// #include "glue.h"
//
// extern void* goLockCB(void*, void**);
// extern void  goUnlockCB(void*, void*, void* const*); 
// extern void  goDisplayCB(void*, void*);
//
// static void goSetCallbacks(libvlc_media_player_t* mp, void* userdata) {
//    libvlc_video_set_callbacks(mp, goLockCB, goUnlockCB, goDisplayCB, userdata);
// }
import "C"
import (
	"os"
	"unsafe"
)

type Player struct {
	ptr *C.libvlc_media_player_t
}

// Retain increments the reference count of this player.
func (this *Player) Retain() (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_player_retain(this.ptr)
	return
}

// Release decreases the reference count of the instance and destroys it
// when it reaches zero.
func (this *Player) Release() (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_player_release(this.ptr)
	return
}

// Media returns the media currently associated with this player.
func (this *Player) Media() (*Media, error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_media_player_get_media(this.ptr); c != nil {
		return &Media{c}, nil
	}

	return nil, checkError()
}

// SetMedia sets the new media to be used by this player. If existing media is
// loaded, it will be destroyed.
func (this *Player) SetMedia(m *Media) error {
	if this.ptr == nil || m.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_player_set_media(this.ptr, m.ptr)
	return checkError()
}

// Events returns an event manager for this player.
func (this *Player) Events() (*EventManager, error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_media_player_event_manager(this.ptr); c != nil {
		return NewEventManager(c), nil
	}

	return nil, checkError()
}

// IsPlaying returns whether or not this player is currently playing.
func (this *Player) IsPlaying() bool {
	if this.ptr == nil {
		return false
	}
	return C.libvlc_media_player_is_playing(this.ptr) != 0
}

// Play begins playback.
func (this *Player) Play() (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	if C.libvlc_media_player_play(this.ptr) < 0 {
		err = checkError()
	}

	return
}

// TogglePause toggles between pause and resume.
// Has no effect if no media is loaded.
func (this *Player) TogglePause(pause bool) (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	if pause {
		C.libvlc_media_player_set_pause(this.ptr, 1)
	} else {
		C.libvlc_media_player_set_pause(this.ptr, 0)
	}
	return
}

// Pause pauses playback.
// Has no effect if no media is loaded.
func (this *Player) Pause() (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_player_pause(this.ptr)
	return
}

// Stop stops playback.
// Has no effect if no media is loaded.
func (this *Player) Stop() (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_player_stop(this.ptr)
	return
}

// SetCallbacks set callbacks and private data to render decoded video to a
// custom area in memory. Use libvlc_video_set_format() to configure the decoded
// format.
//
// Whenever a new video frame needs to be decoded, the lock callback is
// invoked. Depending on the video chroma, one or three pixel planes of
// adequate dimensions must be returned. Those planes must be aligned on
// 32-bytes boundaries.
//
// When the video frame is decoded, the unlock callback is invoked. The
// second parameter to the callback is the return value of the lock callback.
// The third parameter conveys the pixel planes for convenience.
//
// When the video frame needs to be shown, as determined by the media playback
// clock, the display callback is invoked. The second parameter also conveys
// the return value from the lock callback.
func (this *Player) SetCallbacks(lh LockHandler, uh UnlockHandler, dh DisplayHandler, userdata interface{}) (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.goSetCallbacks(this.ptr, unsafe.Pointer(&memRenderReq{lh, uh, dh, userdata}))
	return
}

// SetFormat specifies the decoded video chroma and dimensions. This only works
// in combination with Player.SetCallbacks().
//
// The chroma parameter should be a four-character string identifying the chroma
// (e.g. "RV32" or "I420").
//
// Width and height indicate the pixel dimensions and pitch is the line pitch
// in bytes.
func (this *Player) SetFormat(chroma string, width, height, pitch uint) error {
	if this.ptr == nil {
		return os.EINVAL
	}

	c := C.CString(chroma)
	C.libvlc_video_set_format(this.ptr, c, C.uint(width), C.uint(height), C.uint(pitch))
	C.free(unsafe.Pointer(c))
	return nil
}

// SetNSObject sets the NSView handler where the media player should render its
// video output.
//
// Use the vout called "macosx".
//
// The drawable is an NSObject that follow the VLCOpenGLVideoViewEmbedding
// protocol:
//
//     @protocol VLCOpenGLVideoViewEmbedding <NSObject>
//     - (void)addVoutSubview:(NSView *)view;
//     - (void)removeVoutSubview:(NSView *)view;
//     @end
//
// Or it can be an NSView object.
//
// If you want to use it along with Qt4 see the QMacCocoaViewContainer. Then
// the following code should work:
// 
//     NSView *video = [[NSView alloc] init];
//     QMacCocoaViewContainer *container = new QMacCocoaViewContainer(video, parent);
//     libvlc_media_player_set_nsobject(mp, video);
//     [video release];
//
// You can find a live example in VLCVideoView in VLCKit.framework.
func (this *Player) SetNSObject(drawable uintptr) (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_player_set_nsobject(this.ptr, unsafe.Pointer(drawable))
	return
}

// NSObject returns the NSView handler previously set with Player.SetNSObject().
func (this *Player) NSObject() (drawable uintptr, err error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}

	return uintptr(C.libvlc_media_player_get_nsobject(this.ptr)), checkError()
}

// SetAGL set the agl handler where the media player should render its video output.
func (this *Player) SetAGL(drawable uint32) (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_player_set_agl(this.ptr, C.uint32_t(drawable))
	return
}

// AGL returns the agl handler where the media player should render its video output.
func (this *Player) AGL() (drawable uint32, err error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}

	return uint32(C.libvlc_media_player_get_agl(this.ptr)), checkError()
}

// Set an X Window System drawable where the media player should render its
// video output. If LibVLC was built without X11 output support, this has
// no effects.
//
// The specified identifier must correspond to an existing Input/Output class X11
// window. Pixmaps are /not/ supported. The caller shall ensure that the X11
// server is the same as the one the VLC instance has been configured with.
func (this *Player) SetXWindow(drawable uint32) (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_player_set_xwindow(this.ptr, C.uint32_t(drawable))
	return
}

// XWindow returns the X Window System window identifier previously set with
// Player.SetXWindow(). Note that this will return the identifier even if VLC is
// not currently using it (for instance if it is playing an audio-only input).
func (this *Player) XWindow() (drawable uint32, err error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}

	return uint32(C.libvlc_media_player_get_xwindow(this.ptr)), checkError()
}

// SetHwnd sets a Win32/Win64 API window handle (HWND) where the media player
// should render its video output. If LibVLC was built without Win32/Win64 API
// output support, then this has no effects.
func (this *Player) SetHwnd(drawable uintptr) (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_media_player_set_hwnd(this.ptr, unsafe.Pointer(drawable))
	return
}

// Hwnd returns the Windows API window handle (HWND) previously set with
// Player.SetHwnd(). The handle will be returned even if LibVLC is not currently
// outputting any video to it.
func (this *Player) Hwnd() (drawable uintptr, err error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}

	return uintptr(C.libvlc_media_player_get_hwnd(this.ptr)), checkError()
}

// Length returns the current movie length in milliseconds.
func (this *Player) Length() (int64, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int64(C.libvlc_media_player_get_length(this.ptr)), checkError()
}

// Time returns the current movie time in milliseconds.
func (this *Player) Time() (int64, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int64(C.libvlc_media_player_get_time(this.ptr)), checkError()
}

// SetTime sets the movie time in milliseconds. This has no effect if no media
// is being played. Not all formats and protocols support this.
func (this *Player) SetTime(v int64) {
	if this.ptr != nil {
		C.libvlc_media_player_set_time(this.ptr, C.libvlc_time_t(v))
	}
}

// Position returns the current movie position.
func (this *Player) Position() (float32, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return float32(C.libvlc_media_player_get_position(this.ptr)), checkError()
}

// SetPosition sets the movie position. This has no effect if playback is not
// enabled. This might not work depending on the underlying input format and protocol.
func (this *Player) SetPosition(v float32) {
	if this.ptr != nil {
		C.libvlc_media_player_set_position(this.ptr, C.float(v))
	}
}

// ChapterCount returns the number of available movie chapters.
func (this *Player) ChapterCount() (int, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int(C.libvlc_media_player_get_chapter_count(this.ptr)), checkError()
}

// Chapter returns the current movie chapter.
func (this *Player) Chapter() (int, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int(C.libvlc_media_player_get_chapter(this.ptr)), checkError()
}

// SetChapter sets the current movie chapter. This has no effect if playback is not
// enabled. This might not work depending on the underlying input format and protocol.
func (this *Player) SetChapter(v int) {
	if this.ptr != nil {
		C.libvlc_media_player_set_chapter(this.ptr, C.int(v))
	}
}

// WillPlay returns true if the player is able to play.
func (this *Player) WillPlay() bool {
	if this.ptr != nil {
		return C.libvlc_media_player_will_play(this.ptr) != 0
	}
	return false
}

// TitleChapterCount returns the number of available movie chapters for the given title.
func (this *Player) TitleChapterCount(title int) (int, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int(C.libvlc_media_player_get_chapter_count_for_title(this.ptr, C.int(title))), checkError()
}

// TitleCount returns the number of available movie titles.
func (this *Player) TitleCount(title int) (int, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int(C.libvlc_media_player_get_title_count(this.ptr)), checkError()
}

// Title returns the current movie title.
func (this *Player) Title() (int, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int(C.libvlc_media_player_get_title(this.ptr)), checkError()
}

// SetTitle sets the current movie title.
func (this *Player) SetTitle(v int) {
	if this.ptr != nil {
		C.libvlc_media_player_set_title(this.ptr, C.int(v))
	}
}

// PreviousChapter sets the previous chapter if applicable.
func (this *Player) PreviousChapter() (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}
	C.libvlc_media_player_previous_chapter(this.ptr)
	return checkError()
}

// NextChapter sets the next chapter if applicable.
func (this *Player) NextChapter() (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}
	C.libvlc_media_player_next_chapter(this.ptr)
	return checkError()
}

// Rate returns the current movie playback rate.
// Note: Depending on the underlying media, the requested rate may be
// different from the real playback rate.
func (this *Player) Rate() (float32, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return float32(C.libvlc_media_player_get_rate(this.ptr)), checkError()
}

// SetRate sets the requested movie playback rate.
func (this *Player) SetRate(v float32) error {
	if this.ptr == nil {
		return os.EINVAL
	}
	C.libvlc_media_player_set_rate(this.ptr, C.float(v))
	return checkError()
}

// State returns the current movie state.
func (this *Player) State() (MediaState, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return MediaState(C.libvlc_media_player_get_state(this.ptr)), checkError()
}

// Fps returns the current movie frame rate.
func (this *Player) Fps() (float32, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return float32(C.libvlc_media_player_get_fps(this.ptr)), checkError()
}

// OutputCount returns the number of outputs the current media has.
func (this *Player) OutputCount() (int, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int(C.libvlc_media_player_has_vout(this.ptr)), checkError()
}

// CanSeek returns whether or not seeking is allowed for the current media.
func (this *Player) CanSeek() (bool, error) {
	if this.ptr == nil {
		return false, os.EINVAL
	}
	return C.libvlc_media_player_is_seekable(this.ptr) != 0, checkError()
}

// CanPause returns whether or not pause/resume is allowed for the current media.
func (this *Player) CanPause() (bool, error) {
	if this.ptr == nil {
		return false, os.EINVAL
	}
	return C.libvlc_media_player_can_pause(this.ptr) != 0, checkError()
}

// NextFrame jumps to the next frame if applicable.
func (this *Player) NextFrame() error {
	if this.ptr == nil {
		return os.EINVAL
	}
	C.libvlc_media_player_next_frame(this.ptr)
	return checkError()
}

// ToggleFullscreen switches between fullscreen and windowed modes on
// non-embedded video outputs.
//
// Note: The same limitations apply to this as to Player.SetFullscreen()
func (this *Player) ToggleFullscreen() error {
	if this.ptr == nil {
		return os.EINVAL
	}
	C.libvlc_toggle_fullscreen(this.ptr)
	return checkError()
}

// SetFullscreen switches from fullscreen to windowed mode or vice-versa.
// This applies only to non-embedded video outputs.
//
// Note: With most window managers, only a top-level windows can be in
// full-screen mode. Hence, this function will not operate properly if
// Player.SetXWindow() was used to embed the video in a non-top-level window.
// In that case, the embedding window must be reparented to the root window
// /before/ fullscreen mode is enabled. You will want to reparent it back to its
// normal parent when disabling fullscreen.
func (this *Player) SetFullscreen(toggle bool) error {
	if this.ptr == nil {
		return os.EINVAL
	}

	if toggle {
		C.libvlc_set_fullscreen(this.ptr, 1)
	} else {
		C.libvlc_set_fullscreen(this.ptr, 0)
	}

	return checkError()
}

// Fullscreen returns wether or not we are currently in fullscreen mode.
func (this *Player) Fullscreen() (bool, error) {
	if this.ptr == nil {
		return false, os.EINVAL
	}
	return C.libvlc_get_fullscreen(this.ptr) != 0, checkError()
}

// SetKeyInput enables or disables key press events handling, according to the
// LibVLC hotkeys configuration. By default and for historical reasons, keyboard
// events are handled by the LibVLC video widget.
//
// Note: On X11, there can be only one subscriber for key press and mouse
// click events per window. If your application has subscribed to those events
// for the X window ID of the video widget, then LibVLC will not be able to
// handle key presses and mouse clicks in any case.
//
// Note: This function is only implemented for X11 and Win32 at the moment.
func (this *Player) SetKeyInput(toggle bool) error {
	if this.ptr == nil {
		return os.EINVAL
	}

	if toggle {
		C.libvlc_video_set_key_input(this.ptr, 1)
	} else {
		C.libvlc_video_set_key_input(this.ptr, 0)
	}

	return checkError()
}

// SetMouseInput enables or disables mouse click events handling. By default,
// those events are handled. This is needed for DVD menus to work, as well as a
// few video filters such as "puzzle".
//
// Note: On X11, there can be only one subscriber for key press and mouse
// click events per window. If your application has subscribed to those events
// for the X window ID of the video widget, then LibVLC will not be able to
// handle key presses and mouse clicks in any case.
//
// Note: This function is only implemented for X11 and Win32 at the moment.
func (this *Player) SetMouseInput(toggle bool) error {
	if this.ptr == nil {
		return os.EINVAL
	}

	if toggle {
		C.libvlc_video_set_mouse_input(this.ptr, 1)
	} else {
		C.libvlc_video_set_mouse_input(this.ptr, 0)
	}

	return checkError()
}

// Size returns the pixel dimensions of a video.
// vidnum is the number of the target video. Most commonly starts at 0.
func (this *Player) Size(vidnum uint) (width, height uint, err error) {
	if this.ptr == nil {
		return 0, 0, os.EINVAL
	}

	var w, h C.uint
	C.libvlc_video_get_size(this.ptr, C.uint(vidnum), &w, &h)
	return uint(w), uint(h), nil
}

// Get the mouse pointer coordinates over a video.
// Coordinates are expressed in terms of the decoded video resolution,
// /not/ in terms of pixels on the screen/viewport (to get the latter,
// you can query your windowing system directly).
//
// Either of the coordinates may be negative or larger than the corresponding
// dimension of the video, if the cursor is outside the rendering area.
//
// Note: The coordinates may be out-of-date if the pointer is not located
// on the video rendering area. LibVLC does not track the pointer if it is
// outside of the video widget.
//
// Note: LibVLC does not support multiple pointers (it does of course support
// multiple input devices sharing the same pointer) at the moment.
// 
// vidnum is the number of the target video. Most commonly starts at 0.
func (this *Player) Cursor(vidnum uint) (cx, cy int, err error) {
	if this.ptr == nil {
		return 0, 0, os.EINVAL
	}

	var x, y C.int
	C.libvlc_video_get_cursor(this.ptr, C.uint(vidnum), &x, &y)
	return int(x), int(y), nil
}

// Scale returns the video scaling factor. That is the ratio of the number of
// pixels on screen to the number of pixels in the original decoded video in each
// dimension.
//
// Note: Not all video outputs support scaling.
func (this *Player) Scale() (float32, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return float32(C.libvlc_video_get_scale(this.ptr)), checkError()
}

// SetScale sets the video scaling factor. That is the ratio of the number of
// pixels on screen to the number of pixels in the original decoded video in each
// dimension. Zero is a special value; it will adjust the video to the output
// window/drawable (in windowed mode) or the entire screen.
//
// Note: Not all video outputs support scaling.
func (this *Player) SetScale(v float32) error {
	if this.ptr == nil {
		return os.EINVAL
	}
	C.libvlc_video_set_scale(this.ptr, C.float(v))
	return checkError()
}

// Aspect returns the current aspect ratio.
func (this *Player) Aspect() (s string, err error) {
	if this.ptr == nil {
		return "", os.EINVAL
	}

	if c := C.libvlc_video_get_aspect_ratio(this.ptr); c != nil {
		s = C.GoString(c)
		C.free(unsafe.Pointer(c))
		return
	}

	return "", checkError()
}

// SetAspect sets the current aspect ratio.
func (this *Player) SetAspect(v string) error {
	if this.ptr == nil {
		return os.EINVAL
	}
	c := C.CString(v)
	C.libvlc_video_set_aspect_ratio(this.ptr, c)
	C.free(unsafe.Pointer(c))
	return checkError()
}

// SubTile returns the current video subtitle or -1 if none is set.
func (this *Player) SubTile() (int, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int(C.libvlc_video_get_spu(this.ptr)), checkError()
}

// SubTileCount returns the number of available subtitles.
func (this *Player) SubTileCount() (int, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int(C.libvlc_video_get_spu_count(this.ptr)), checkError()
}

// SubTileDescription returns descriptions for the current subtitle track.
//
// Note: make sure to call TrackDescriptionList.Release() when you are done with it.
func (this *Player) SubTileDescription() (TrackDescriptionList, error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_video_get_spu_description(this.ptr); c != nil {
		var l TrackDescriptionList
		l.fromC(c)
		return l, nil
	}

	return nil, checkError()
}

// SetSubtitle sets the current subtitle track.
func (this *Player) SetSubtitle(s uint) (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	if C.libvlc_video_set_spu(this.ptr, C.uint(s)) != 0 {
		err = checkError()
	}

	return
}

// SetSubtitle sets the current subtitle from a file.
func (this *Player) SetSubtitleFile(path string) (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	c := C.CString(path)
	C.libvlc_video_set_subtitle_file(this.ptr, c)
	C.free(unsafe.Pointer(c))
	return checkError()
}

// ChapterDescription returns descriptions of available chapters for a specific title.
//
// Note: make sure to call TrackDescriptionList.Release() when you are done with it.
func (this *Player) ChapterDescription(title int) (TrackDescriptionList, error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_video_get_chapter_description(this.ptr, C.int(title)); c != nil {
		var l TrackDescriptionList
		l.fromC(c)
		return l, nil
	}

	return nil, checkError()
}

// CropGeometry returns the current crop filter geometry.
func (this *Player) CropGeometry() (s string, err error) {
	if this.ptr == nil {
		return "", os.EINVAL
	}

	if c := C.libvlc_video_get_crop_geometry(this.ptr); c != nil {
		s = C.GoString(c)
		C.free(unsafe.Pointer(c))
	}

	return "", checkError()
}

// SetCropGeometry sets the current crop filter geometry. Specify an empty
// string to clear the filter.
func (this *Player) SetCropGeometry(s string) (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	var c *C.char

	if len(s) > 0 {
		c = C.CString(s)
		defer C.free(unsafe.Pointer(c))
	}

	C.libvlc_video_set_crop_geometry(this.ptr, c)
	return checkError()
}

// Teletext returns the current requested teletext page.
func (this *Player) Teletext() (int, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int(C.libvlc_video_get_teletext(this.ptr)), checkError()
}

// SetTeletext sets a new teletext page to retrieve.
func (this *Player) SetTeletext(page int) error {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_video_set_teletext(this.ptr, C.int(page))
	return checkError()
}

// ToggleTeletext toggles transparent teletext status on video output.
func (this *Player) ToggleTeletext() error {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_toggle_teletext(this.ptr)
	return checkError()
}

// VideoTrackCount returns the number of video tracks in the current media.
func (this *Player) VideoTrackCount() (int, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int(C.libvlc_video_get_track_count(this.ptr)), checkError()
}

// VideoDescription returns descriptions for the current video tracks.
//
// Note: make sure to call TrackDescriptionList.Release() when you are done with it.
func (this *Player) VideoDescription() (TrackDescriptionList, error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_video_get_track_description(this.ptr); c != nil {
		var l TrackDescriptionList
		l.fromC(c)
		return l, nil
	}

	return nil, checkError()
}

// VideoTrack returns the current video track.
func (this *Player) VideoTrack() (int, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int(C.libvlc_video_get_track(this.ptr)), checkError()
}

// SetVideoTrack sets the current video track.
func (this *Player) SetVideoTrack(track int) (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	if C.libvlc_video_set_track(this.ptr, C.int(track)) != 0 {
		err = checkError()
	}

	return
}

// TakeSnapshot takes a snapshot of the selected video output ans saves it
// to the specified file.
//
// If width AND height are both 0, the original size is used.
// If width OR height is 0, the original aspect ratio is preserved.
//
// Vidnum is the number of the video output (typically 0 for the first/only one)
func (this *Player) TakeSnapshot(path string, vidnum, width, height uint) (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	c := C.CString(path)

	if C.libvlc_video_take_snapshot(this.ptr, C.uint(vidnum), c, C.uint(width), C.uint(height)) != 0 {
		err = checkError()
	}

	C.free(unsafe.Pointer(c))
	return
}

// SetDeinterlace sets the deinterlace filter. Supply an empty string to disable
// the filter.
func (this *Player) SetDeinterlace(f string) (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	var c *C.char
	if len(f) > 0 {
		c := C.CString(f)
		defer C.free(unsafe.Pointer(c))
	}

	C.libvlc_video_set_deinterlace(this.ptr, c)
	return
}

// MarqueeOption returns an integer marquee option value.
func (this *Player) MarqueeOption(option MarqueeOption) (int, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int(C.libvlc_video_get_marquee_int(this.ptr, C.uint(option))), checkError()
}

// SetMarqueeOption sets an integer marquee option value.
func (this *Player) SetMarqueeOption(option MarqueeOption, v int) error {
	if this.ptr == nil {
		return os.EINVAL
	}
	C.libvlc_video_set_marquee_int(this.ptr, C.uint(option), C.int(v))
	return checkError()
}

// MarqueeOptionString returns a string marquee option value.
func (this *Player) MarqueeOptionString(option MarqueeOption) (s string, err error) {
	if this.ptr == nil {
		return "", os.EINVAL
	}

	if c := C.libvlc_video_get_marquee_string(this.ptr, C.uint(option)); c != nil {
		s = C.GoString(c)
		C.free(unsafe.Pointer(c))
	}

	return s, checkError()
}

// SetMarqueeOptionString sets a string marquee option value.
func (this *Player) SetMarqueeOptionString(option MarqueeOption, s string) error {
	if this.ptr == nil {
		return os.EINVAL
	}
	c := C.CString(s)
	C.libvlc_video_set_marquee_string(this.ptr, C.uint(option), c)
	C.free(unsafe.Pointer(c))
	return checkError()
}

// LogoOption returns an integer logo option.
func (this *Player) LogoOption(option LogoOption) (int, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int(C.libvlc_video_get_logo_int(this.ptr, C.uint(option))), checkError()
}

// SetLogoOption sets an integer logo option value.
// Options that take a different type value are ignored.
// Passing LOEnable as option value has the side effect of starting (arg !0) or
// stopping (arg 0) the logo filter.
func (this *Player) SetLogoOption(option MarqueeOption, v int) error {
	if this.ptr == nil {
		return os.EINVAL
	}
	C.libvlc_video_set_logo_int(this.ptr, C.uint(option), C.int(v))
	return checkError()
}

// SetLogoOptionString sets a string logo option value.
func (this *Player) SetLogoOptionString(option LogoOption, s string) error {
	if this.ptr == nil {
		return os.EINVAL
	}
	c := C.CString(s)
	C.libvlc_video_set_logo_string(this.ptr, C.uint(option), c)
	C.free(unsafe.Pointer(c))
	return checkError()
}

// AdjustOption returns an integer adjustment option.
func (this *Player) AdjustOption(option AdjustOption) (int, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int(C.libvlc_video_get_adjust_int(this.ptr, C.uint(option))), checkError()
}

// SetAdjustOption sets an integer adjustment option value.
// Options that take a different type value are ignored.
// Passing AOEnable as option value has the side effect of starting (arg !0) or
// stopping (arg 0) the logo filter.
func (this *Player) SetAdjustOption(option AdjustOption, v int) error {
	if this.ptr == nil {
		return os.EINVAL
	}
	C.libvlc_video_set_logo_int(this.ptr, C.uint(option), C.int(v))
	return checkError()
}

// AdjustOptionFloat returns a float adjustment option.
func (this *Player) AdjustOptionFloat(option AdjustOption) (float32, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return float32(C.libvlc_video_get_adjust_float(this.ptr, C.uint(option))), checkError()
}

// SetAdjustOptionFloat sets a float adjustment option value.
// Options that take a different type value are ignored.
func (this *Player) SetAdjustOptionFloat(option AdjustOption, v float32) error {
	if this.ptr == nil {
		return os.EINVAL
	}
	C.libvlc_video_set_adjust_float(this.ptr, C.uint(option), C.float(v))
	return checkError()
}

// AudioOutput returns a list of available audio outputs.
//
// Note: Be sure to call AudioOutputList.Release() after you are done with the list.
func (this *Player) AudioOutput() (AudioOutputList, error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_audio_output_list_get(this.ptr); c != nil {
		var l AudioOutputList
		l.fromC(c)
		return l, nil
	}

	return nil, checkError()
}

// SetAudioOutput sets the current audio output. Changes will be applied after
// stop and play.
func (this *Player) SetAudioOutput(output string) (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	c := C.CString(output)

	if C.libvlc_audio_output_set(this.ptr, c) == 0 {
		err = checkError()
	}

	C.free(unsafe.Pointer(c))
	return
}

// AudioDeviceCount returns the number of devices for audio output. These devices
// are hardware oriented like analog or digital output of sound cards.
func (this *Player) AudioDeviceCount(output string) (int, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}

	c := C.CString(output)
	defer C.free(unsafe.Pointer(c))
	return int(C.libvlc_audio_output_device_count(this.ptr, c)), checkError()
}

// AudioDeviceName returns the long name of an audio device.
// If it is not available, the short name is given.
func (this *Player) AudioDeviceName(output string, device int) (s string, err error) {
	if this.ptr == nil {
		return "", os.EINVAL
	}

	c := C.CString(output)
	defer C.free(unsafe.Pointer(c))

	if r := C.libvlc_audio_output_device_longname(this.ptr, c, C.int(device)); r != nil {
		s = C.GoString(r)
		C.free(unsafe.Pointer(r))
		return
	}

	return "", checkError()
}

// AudioDeviceId returns the id of an audio device.
func (this *Player) AudioDeviceId(output string, device int) (s string, err error) {
	if this.ptr == nil {
		return "", os.EINVAL
	}

	c := C.CString(output)
	defer C.free(unsafe.Pointer(c))

	if r := C.libvlc_audio_output_device_id(this.ptr, c, C.int(device)); r != nil {
		s = C.GoString(r)
		C.free(unsafe.Pointer(r))
		return
	}

	return "", checkError()
}

// SetAudioDevice sets the current audio output device. Changes will be applied after
// stop and play.
func (this *Player) SetAudioDevice(output, deviceid string) (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	a := C.CString(output)
	b := C.CString(deviceid)
	C.libvlc_audio_output_device_set(this.ptr, a, b)
	C.free(unsafe.Pointer(a))
	C.free(unsafe.Pointer(b))
	return
}

// AudioDeviceType return the current audio device type.
// Device type describes something like character of output sound - stereo
// sound, 2.1, 5.1 etc
func (this *Player) AudioDeviceType() (AudioDevice, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}

	return AudioDevice(C.libvlc_audio_output_get_device_type(this.ptr)), checkError()
}

// SetAudioDeviceType sets the current audio device type.
// Device type describes something like character of output sound - stereo
// sound, 2.1, 5.1 etc
func (this *Player) SetAudioDeviceType(ad AudioDevice) error {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_audio_output_set_device_type(this.ptr, C.int(ad))
	return checkError()
}

// ToggleMute toggles the current mute status.
func (this *Player) ToggleMute() error {
	if this.ptr == nil {
		return os.EINVAL
	}
	C.libvlc_audio_toggle_mute(this.ptr)
	return checkError()
}

// IsMute returns whether or not mute is enabled.
func (this *Player) IsMute() (bool, error) {
	if this.ptr == nil {
		return false, os.EINVAL
	}
	return C.libvlc_audio_get_mute(this.ptr) != 0, checkError()
}

// SetMute sets mute mode to the specified value.
func (this *Player) SetMute(toggle bool) error {
	if this.ptr == nil {
		return os.EINVAL
	}

	if toggle {
		C.libvlc_audio_set_mute(this.ptr, 1)
	} else {
		C.libvlc_audio_set_mute(this.ptr, 0)
	}

	return checkError()
}

// Volume returns the current audio level.
func (this *Player) Volume() (int, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int(C.libvlc_audio_get_volume(this.ptr)), checkError()
}

// SetVolume sets the current audio level.
func (this *Player) SetVolume(v int) error {
	if this.ptr == nil {
		return os.EINVAL
	}

	C.libvlc_audio_set_volume(this.ptr, C.int(v))
	return checkError()
}

// AudioTrackCount returns the number of available audio tracks.
func (this *Player) AudioTrackCount() (int, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int(C.libvlc_audio_get_track_count(this.ptr)), checkError()
}

// AudioDescription returns descriptions for the current audio tracks.
//
// Note: make sure to call TrackDescriptionList.Release() when you are done with it.
func (this *Player) AudioDescription() (TrackDescriptionList, error) {
	if this.ptr == nil {
		return nil, os.EINVAL
	}

	if c := C.libvlc_audio_get_track_description(this.ptr); c != nil {
		var l TrackDescriptionList
		l.fromC(c)
		return l, nil
	}

	return nil, checkError()
}

// AudioTrack returns the current audio track.
func (this *Player) AudioTrack() (int, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int(C.libvlc_audio_get_track(this.ptr)), checkError()
}

// SetAudioTrack sets the current audio track.
func (this *Player) SetAudioTrack(track int) (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	if C.libvlc_audio_set_track(this.ptr, C.int(track)) != 0 {
		err = checkError()
	}

	return
}

// AudioChannel returns the current audio channel.
func (this *Player) AudioChannel() (int, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int(C.libvlc_audio_get_channel(this.ptr)), checkError()
}

// SetAudioChannel sets the current audio channel.
func (this *Player) SetAudioChannel(channel int) (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	if C.libvlc_audio_set_channel(this.ptr, C.int(channel)) != 0 {
		err = checkError()
	}

	return
}

// AudioDelay returns the current audio delay.
func (this *Player) AudioDelay() (int64, error) {
	if this.ptr == nil {
		return 0, os.EINVAL
	}
	return int64(C.libvlc_audio_get_delay(this.ptr)), checkError()
}

// SetAudioDelay sets the current audio delay.
func (this *Player) SetAudioDelay(delay int64) (err error) {
	if this.ptr == nil {
		return os.EINVAL
	}

	if C.libvlc_audio_set_delay(this.ptr, C.int64_t(delay)) != 0 {
		err = checkError()
	}

	return
}
