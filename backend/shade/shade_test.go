package shade_test

import (
	"testing"

	"github.com/wizdom13/SecureCloudEngine/backend/shade"
	"github.com/wizdom13/SecureCloudEngine/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	name := "TestShade"
	fstests.Run(t, &fstests.Opt{
		RemoteName:      name + ":",
		NilObject:       (*shade.Object)(nil),
		SkipInvalidUTF8: true,
		ExtraConfig: []fstests.ExtraConfigItem{
			{Name: name, Key: "eventually_consistent_delay", Value: "7"},
		},
	})
}
