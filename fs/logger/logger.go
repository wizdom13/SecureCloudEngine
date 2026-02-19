// Package logger implements testing for the sync (and bisync) logger
package logger

import (
	_ "/backend/all" // import all backends
	"/cmd"
	_ "/cmd/all"    // import all commands
	_ "/lib/plugin" // import plugins
)

// Main enables the testscript package. See:
// https://bitfieldconsulting.com/golang/cli-testing
// https://pkg.go.dev/github.com/rogpeppe/go-internal@v1.11.0/testscript
func Main() {
	cmd.Main()
}
