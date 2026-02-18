// Test Files filesystem interface
package filescom_test

import (
	"testing"

	"github.com/wizdom13/SecureCloudEngine/backend/filescom"
	"github.com/wizdom13/SecureCloudEngine/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestFilesCom:",
		NilObject:  (*filescom.Object)(nil),
	})
}
