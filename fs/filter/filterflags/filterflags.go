// Package filterflags implements command line flags to set up a filter
package filterflags

import (
	"github.com/wizdom13/SecureCloudEngine/fs/config/flags"
	"github.com/wizdom13/SecureCloudEngine/fs/filter"
	"github.com/spf13/pflag"
)

// AddFlags adds the non filing system specific flags to the command
func AddFlags(flagSet *pflag.FlagSet) {
	flags.AddFlagsFromOptions(flagSet, "", filter.OptionsInfo)
}
