// Sync files and directories to and from local and remote object stores
//
// Nick Craig-Wood <nick@craig-wood.com>
package main

import (
	"os"

	_ "github.com/wizdom13/SecureCloudEngine/backend/all" // import all backends
	"github.com/wizdom13/SecureCloudEngine/cmd"
	_ "github.com/wizdom13/SecureCloudEngine/cmd/all"    // import all commands
	_ "github.com/wizdom13/SecureCloudEngine/lib/plugin" // import plugins
)

func main() {
	if !isLaunchedBySecureCloud() {
		os.Exit(1)
	}

	cmd.Main()
}
