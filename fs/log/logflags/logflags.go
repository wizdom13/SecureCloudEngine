// Package logflags implements command line flags to set up the log
package logflags

import (
	"github.com/wizdom13/SecureCloudEngine/fs/config/flags"
	"github.com/wizdom13/SecureCloudEngine/fs/log"
	"github.com/spf13/pflag"
)

// AddFlags adds the log flags to the flagSet
func AddFlags(flagSet *pflag.FlagSet) {
	flags.AddFlagsFromOptions(flagSet, "", log.OptionsInfo)
}
