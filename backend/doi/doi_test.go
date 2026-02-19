// Test DOI filesystem interface
package doi

import (
	"testing"

	"github.com/wizdom13/SecureCloudEngine/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestDoi:",
		NilObject:  (*Object)(nil),
	})
}
