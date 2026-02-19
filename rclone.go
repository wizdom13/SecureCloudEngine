// Sync files and directories to and from local and remote object stores
//
// Nick Craig-Wood <nick@craig-wood.com>
package main

import (
	"os"

	_ "/backend/all" // import all backends
	"/cmd"
	_ "/cmd/all"    // import all commands
	_ "/lib/plugin" // import plugins
)

func main() {
	if !isLaunchedBySecureCloud() {
		os.Exit(1)
	}

	cmd.Main()
}
