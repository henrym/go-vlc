# This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
# license. Its contents can be found at:
# http://creativecommons.org/publicdomain/zero/1.0

include $(GOROOT)/src/Make.inc

TARG = github.com/jteeuwen/go-vlc
GOFILES = const.go
CGOFILES = vlc.go instance.go eventmanager.go event.go callback.go media.go \
	log.go stats.go trackinfo.go player.go misc.go medialist.go listplayer.go \
	library.go

include $(GOROOT)/src/Make.pkg
