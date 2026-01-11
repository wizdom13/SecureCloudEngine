// Test Webdav filesystem interface
package webdav

import (
	"os"
	"os/exec"
	"testing"

	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fstest"
	"github.com/rclone/rclone/fstest/fstests"
)

func requireDockerTestServer(t *testing.T) {
	t.Helper()
	if _, err := os.Stat("fstest/testserver/init.d/docker.bash"); err != nil {
		t.Skip("skipping WebDAV integration tests that require docker-based test servers")
	}
}

func requireRcloneBinary(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("rclone"); err != nil {
		t.Skip("skipping WebDAV rclone integration tests; rclone binary not found in PATH")
	}
}

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	requireDockerTestServer(t)
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestWebdavNextcloud:",
		NilObject:  (*Object)(nil),
		ChunkedUpload: fstests.ChunkedUploadConfig{
			MinChunkSize: 1 * fs.Mebi,
		},
	})
}

// TestIntegration runs integration tests against the remote
func TestIntegration2(t *testing.T) {
	if *fstest.RemoteName != "" {
		t.Skip("skipping as -remote is set")
	}
	requireDockerTestServer(t)
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestWebdavOwncloud:",
		NilObject:  (*Object)(nil),
		ChunkedUpload: fstests.ChunkedUploadConfig{
			Skip: true,
		},
	})
}

// TestIntegration runs integration tests against the remote
func TestIntegration3(t *testing.T) {
	if *fstest.RemoteName != "" {
		t.Skip("skipping as -remote is set")
	}
	requireRcloneBinary(t)
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestWebdavRclone:",
		NilObject:  (*Object)(nil),
		ChunkedUpload: fstests.ChunkedUploadConfig{
			Skip: true,
		},
	})
}

// TestIntegration runs integration tests against the remote
func TestIntegration4(t *testing.T) {
	if *fstest.RemoteName != "" {
		t.Skip("skipping as -remote is set")
	}
	requireDockerTestServer(t)
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestWebdavNTLM:",
		NilObject:  (*Object)(nil),
	})
}

func (f *Fs) SetUploadChunkSize(cs fs.SizeSuffix) (fs.SizeSuffix, error) {
	return f.setUploadChunkSize(cs)
}
