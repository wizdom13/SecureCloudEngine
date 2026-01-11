// Test Sftp filesystem interface

//go:build !plan9

package sftp_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/rclone/rclone/backend/sftp"
	"github.com/rclone/rclone/fstest"
	"github.com/rclone/rclone/fstest/fstests"
)

func requireRcloneBinary(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("rclone"); err != nil {
		t.Skip("skipping SFTP integration tests; rclone binary not found in PATH")
	}
}

func requireTestServerScript(t *testing.T, scriptName string) {
	t.Helper()
	scriptPath := filepath.Join("fstest", "testserver", "init.d", scriptName)
	if _, err := os.Stat(scriptPath); err != nil {
		t.Skip("skipping SFTP integration tests; test server script not available")
	}
}

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	requireTestServerScript(t, "TestSFTPOpenssh")
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestSFTPOpenssh:",
		NilObject:  (*sftp.Object)(nil),
	})
}

func TestIntegration2(t *testing.T) {
	if *fstest.RemoteName != "" {
		t.Skip("skipping as -remote is set")
	}
	requireRcloneBinary(t)
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestSFTPRclone:",
		NilObject:  (*sftp.Object)(nil),
	})
}

func TestIntegration3(t *testing.T) {
	if *fstest.RemoteName != "" {
		t.Skip("skipping as -remote is set")
	}
	requireRcloneBinary(t)
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestSFTPRcloneSSH:",
		NilObject:  (*sftp.Object)(nil),
	})
}
