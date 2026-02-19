package protondrive_test

import (
	"testing"

	"github.com/wizdom13/SecureCloudEngine/backend/protondrive"
	"github.com/wizdom13/SecureCloudEngine/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestProtonDrive:",
		NilObject:  (*protondrive.Object)(nil),
	})
}
