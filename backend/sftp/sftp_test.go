// Test Sftp filesystem interface

//go:build !plan9

package sftp_test

import (
	"testing"

	"github.com/wizdom13/SecureCloudEngine/backend/sftp"
	"github.com/wizdom13/SecureCloudEngine/fstest"
	"github.com/wizdom13/SecureCloudEngine/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestSFTPOpenssh:",
		NilObject:  (*sftp.Object)(nil),
	})
}

func TestIntegration2(t *testing.T) {
	if *fstest.RemoteName != "" {
		t.Skip("skipping as -remote is set")
	}
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestSFTPSecureCloudEngine:",
		NilObject:  (*sftp.Object)(nil),
	})
}

func TestIntegration3(t *testing.T) {
	if *fstest.RemoteName != "" {
		t.Skip("skipping as -remote is set")
	}
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestSFTPSecureCloudEngineSSH:",
		NilObject:  (*sftp.Object)(nil),
	})
}
