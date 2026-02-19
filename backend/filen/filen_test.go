package filen

import (
	"testing"

	"github.com/wizdom13/SecureCloudEngine/fstest/fstests"
)

func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestFilen:",
		NilObject:  (*Object)(nil),
	})
}
