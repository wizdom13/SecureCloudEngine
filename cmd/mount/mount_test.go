//go:build linux

package mount

import (
	"testing"

	"github.com/wizdom13/SecureCloudEngine/vfs/vfscommon"
	"github.com/wizdom13/SecureCloudEngine/vfs/vfstest"
)

func TestMount(t *testing.T) {
	vfstest.RunTests(t, false, vfscommon.CacheModeOff, true, mount)
}
