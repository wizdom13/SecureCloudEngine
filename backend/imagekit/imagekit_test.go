package imagekit

import (
	"testing"

	"github.com/wizdom13/SecureCloudEngine/fstest"
	"github.com/wizdom13/SecureCloudEngine/fstest/fstests"
)

func TestIntegration(t *testing.T) {
	debug := true
	fstest.Verbose = &debug
	fstests.Run(t, &fstests.Opt{
		RemoteName:      "TestImageKit:",
		NilObject:       (*Object)(nil),
		SkipFsCheckWrap: true,
	})
}
