// Package vfsflags implements command line flags to set up a vfs
package vfsflags

import (
	"github.com/spf13/pflag"

	"github.com/wizdom13/SecureCloudEngine/fs/config/flags"
	"github.com/wizdom13/SecureCloudEngine/vfs/vfscommon"
)

// AddFlags adds the non filing system specific flags to the command
func AddFlags(flagSet *pflag.FlagSet) {
	flags.AddFlagsFromOptions(flagSet, "", vfscommon.OptionsInfo)
}
