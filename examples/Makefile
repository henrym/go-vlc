# This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
# license. Its contents can be found at:
# http://creativecommons.org/publicdomain/zero/1.0

include $(GOROOT)/src/Make.inc

APPS = tut

all: $(APPS)

tut:
	$(GC) -o tut.$(O) tut.go
	$(LD) -s -o tut tut.$(O)

clean:
	rm -rf *.o *.a *.[$(OS)] [$(OS)].out $(APPS)

