package filelu_test

import (
	"testing"

	"github.com/wizdom13/SecureCloudEngine/fstest/fstests"
)

// TestIntegration runs integration tests for the FileLu backend
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName:      "TestFileLu:",
		NilObject:       nil,
		SkipInvalidUTF8: true,
	})
}
