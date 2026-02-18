// Test Yandex filesystem interface
package yandex_test

import (
	"testing"

	"github.com/wizdom13/SecureCloudEngine/backend/yandex"
	"github.com/wizdom13/SecureCloudEngine/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestYandex:",
		NilObject:  (*yandex.Object)(nil),
	})
}
