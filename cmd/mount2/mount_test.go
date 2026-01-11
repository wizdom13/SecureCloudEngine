//go:build linux

package mount2

import (
	"os"
	"os/exec"
	"testing"

	"github.com/rclone/rclone/vfs/vfscommon"
	"github.com/rclone/rclone/vfs/vfstest"
)

func requireFUSE(t *testing.T) {
	t.Helper()
	if _, err := os.Stat("/dev/fuse"); err != nil {
		t.Skip("skipping mount2 tests; FUSE device not available")
	}
	if _, err := exec.LookPath("fusermount3"); err != nil {
		if _, err := exec.LookPath("fusermount"); err != nil {
			t.Skip("skipping mount2 tests; fusermount not found")
		}
	}
}

func TestMount(t *testing.T) {
	requireFUSE(t)
	vfstest.RunTests(t, false, vfscommon.CacheModeOff, true, mount)
}
