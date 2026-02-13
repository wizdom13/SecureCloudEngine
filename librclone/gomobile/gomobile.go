// Package gomobile exports shims for gomobile use
package gomobile

import (
	"github.com/rclone/rclone/librclone/librclone"

	_ "github.com/rclone/rclone/backend/all" // import all backends
	_ "github.com/rclone/rclone/lib/plugin"  // import plugins

	_ "golang.org/x/mobile/event/key" // make go.mod add this as a dependency
)

// SecureCloudEngineInitialize initializes rclone as a library
func SecureCloudEngineInitialize() {
	librclone.Initialize()
}

// SecureCloudEngineFinalize finalizes the library
func SecureCloudEngineFinalize() {
	librclone.Finalize()
}

// SecureCloudEngineRPCResult is returned from SecureCloudEngineRPC
//
//	Output will be returned as a serialized JSON object
//	Status is a HTTP status return (200=OK anything else fail)
type SecureCloudEngineRPCResult struct {
	Output string
	Status int
}

// SecureCloudEngineRPC has an interface optimised for gomobile, in particular
// the function signature is valid under gobind rules.
//
// https://pkg.go.dev/golang.org/x/mobile/cmd/gobind#hdr-Type_restrictions
func SecureCloudEngineRPC(method string, input string) (result *SecureCloudEngineRPCResult) { //nolint:deadcode
	output, status := librclone.RPC(method, input)
	return &SecureCloudEngineRPCResult{
		Output: output,
		Status: status,
	}
}
