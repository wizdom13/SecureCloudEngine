// Package cmdtest creates a testable interface to rclone main
//
// The interface is used to perform end-to-end test of
// commands, flags, environment variables etc.
package cmdtest

// The rest of this file is a 1:1 copy from rclone.go

import (
	_ "/backend/all" // import all backends
	"/cmd"
	_ "/cmd/all"    // import all commands
	_ "/lib/plugin" // import plugins
)

func main() {
	cmd.Main()
}
