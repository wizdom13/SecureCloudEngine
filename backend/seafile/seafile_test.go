// Test Seafile filesystem interface
package seafile_test

import (
	"os"
	"testing"

	"github.com/rclone/rclone/backend/seafile"
	"github.com/rclone/rclone/fstest/fstests"
)

func requireDockerTestServer(t *testing.T) {
	t.Helper()
	if _, err := os.Stat("fstest/testserver/init.d/docker.bash"); err != nil {
		t.Skip("skipping Seafile integration tests that require docker-based test servers")
	}
}

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	requireDockerTestServer(t)
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestSeafile:",
		NilObject:  (*seafile.Object)(nil),
	})
}
