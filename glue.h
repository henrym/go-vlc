// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0

#include <stdlib.h>
#include <vlc/vlc.h>

extern void  goEventCB(const struct libvlc_event_t*, void*);
extern void* goLockCB(void*, void**);
extern void  goUnlockCB(void*, void*, void* const*);
extern void  goDisplayCB(void*, void*);
