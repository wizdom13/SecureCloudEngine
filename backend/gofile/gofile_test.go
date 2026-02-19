// Test Gofile filesystem interface
package gofile_test

import (
	"testing"

	"github.com/wizdom13/SecureCloudEngine/backend/gofile"
	"github.com/wizdom13/SecureCloudEngine/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestGoFile:",
		NilObject:  (*gofile.Object)(nil),
	})
}
