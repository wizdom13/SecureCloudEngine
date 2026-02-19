//go:build linux

package mount

import (
	"os/exec"
	"testing"

	"github.com/wizdom13/SecureCloudEngine/vfs/vfscommon"
	"github.com/wizdom13/SecureCloudEngine/vfs/vfstest"
)

func TestMount(t *testing.T) {
	if _, err := exec.LookPath("fusermount3"); err != nil {
		t.Skip("fusermount3 not found")
	}
	vfstest.RunTests(t, false, vfscommon.CacheModeOff, true, mount)
}
