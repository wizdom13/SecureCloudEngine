// Test Sia filesystem interface
package sia_test

import (
	"os"
	"testing"

	"github.com/rclone/rclone/backend/sia"
	"github.com/rclone/rclone/fstest/fstests"
)

func requireDockerTestServer(t *testing.T) {
	t.Helper()
	if _, err := os.Stat("fstest/testserver/init.d/docker.bash"); err != nil {
		t.Skip("skipping Sia integration tests that require docker-based test servers")
	}
}

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	requireDockerTestServer(t)
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestSia:",
		NilObject:  (*sia.Object)(nil),
	})
}
