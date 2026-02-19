//go:build linux

package mount2

import (
	"testing"

	"/vfs/vfscommon"
	"/vfs/vfstest"
)

func TestMount(t *testing.T) {
	vfstest.RunTests(t, false, vfscommon.CacheModeOff, true, mount)
}
