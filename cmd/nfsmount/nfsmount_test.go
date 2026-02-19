//go:build unix

package nfsmount

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/wizdom13/SecureCloudEngine/cmd/mountlib"
	"github.com/wizdom13/SecureCloudEngine/cmd/serve/nfs"
	"github.com/wizdom13/SecureCloudEngine/fs/object"
	"github.com/wizdom13/SecureCloudEngine/fstest/testy"
	"github.com/wizdom13/SecureCloudEngine/vfs"
	"github.com/wizdom13/SecureCloudEngine/vfs/vfscommon"
	"github.com/wizdom13/SecureCloudEngine/vfs/vfstest"
)

// Return true if the command ran without error
func commandOK(name string, arg ...string) bool {
	cmd := exec.Command(name, arg...)
	_, err := cmd.CombinedOutput()
	return err == nil
}

func checkNFS(t *testing.T) {
	if runtime.GOOS != "linux" {
		return
	}
	data, err := os.ReadFile("/proc/filesystems")
	if err != nil {
		t.Logf("Failed to read /proc/filesystems: %v", err)
		return // Ignore error
	}
	if !strings.Contains(string(data), "\tnfs\n") {
		t.Logf("NFS not found in /proc/filesystems:\n%s", string(data))
		t.Skip("Skipping test because NFS is not supported by kernel (not listed in /proc/filesystems)")
	}
}

func TestMount(t *testing.T) {
	testy.SkipUnreliable(t)
	checkNFS(t)
	if runtime.GOOS != "darwin" {
		if !commandOK("sudo", "-n", "mount", "--help") {
			t.Skip("Can't run sudo mount without a password")
		}
		if !commandOK("sudo", "-n", "umount", "--help") {
			t.Skip("Can't run sudo umount without a password")
		}
		sudo = true
	}
	for _, cacheType := range []string{"memory", "disk", "symlink"} {
		t.Run(cacheType, func(t *testing.T) {
			nfs.Opt.HandleCacheDir = t.TempDir()
			require.NoError(t, nfs.Opt.HandleCache.Set(cacheType))
			// Check we can create a handler
			_, err := nfs.NewHandler(context.Background(), vfs.New(object.MemoryFs, nil), &nfs.Opt)
			if errors.Is(err, nfs.ErrorSymlinkCacheNotSupported) || errors.Is(err, nfs.ErrorSymlinkCacheNoPermission) {
				t.Skip(err.Error() + ": run with: go test -c && sudo setcap cap_dac_read_search+ep ./nfsmount.test && ./nfsmount.test -test.v")
			}
			require.NoError(t, err)
			// Configure rclone via environment var since the mount gets run in a subprocess
			_ = os.Setenv("SCE_NFS_CACHE_DIR", nfs.Opt.HandleCacheDir)
			_ = os.Setenv("SCE_NFS_CACHE_TYPE", cacheType)
			t.Cleanup(func() {
				_ = os.Unsetenv("SCE_NFS_CACHE_DIR")
				_ = os.Unsetenv("SCE_NFS_CACHE_TYPE")
			})
			mountWrapper := func(VFS *vfs.VFS, mountpoint string, opt *mountlib.Options) (<-chan error, func() error, error) {
				asyncerrors, unmount, err := mount(VFS, mountpoint, opt)
				if err != nil && strings.Contains(err.Error(), "unknown filesystem type 'nfs'") {
					t.Skip("Skipping test because NFS is not supported: " + err.Error())
				}
				return asyncerrors, unmount, err
			}
			vfstest.RunTests(t, false, vfscommon.CacheModeWrites, false, mountWrapper)
		})
	}
}
