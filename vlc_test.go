// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0

package vlc

import (
	"testing"
	"os"
	"fmt"
)

func Test(t *testing.T) {
	fmt.Printf("VLC Version: %s\n", VersionString())

	var err error
	var inst *Instance

	if inst, err = New(os.Args); err != nil {
		t.Error(err.Error())
		return
	}

	defer inst.Release()

	inst.SetLogVerbosity(3)
	fmt.Printf("Log level: %d\n", inst.LogVerbosity())

	inst.StartUI("")
	inst.Wait()
}
