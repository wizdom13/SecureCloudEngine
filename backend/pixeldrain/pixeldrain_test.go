// Test pixeldrain filesystem interface
package pixeldrain_test

import (
	"testing"

	"github.com/wizdom13/SecureCloudEngine/backend/pixeldrain"
	"github.com/wizdom13/SecureCloudEngine/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName:      "TestPixeldrain:",
		NilObject:       (*pixeldrain.Object)(nil),
		SkipInvalidUTF8: true, // Pixeldrain throws an error on invalid utf-8
	})
}
