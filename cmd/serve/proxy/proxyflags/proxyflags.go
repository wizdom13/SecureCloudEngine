// Package proxyflags implements command line flags to set up a proxy
package proxyflags

import (
	"github.com/spf13/pflag"

	"github.com/wizdom13/SecureCloudEngine/cmd/serve/proxy"
	"github.com/wizdom13/SecureCloudEngine/fs/config/flags"
)

// AddFlags adds the non filing system specific flags to the command
func AddFlags(flagSet *pflag.FlagSet) {
	flags.AddFlagsFromOptions(flagSet, "", proxy.OptionsInfo)
}
