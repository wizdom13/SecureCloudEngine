// Package librclone exports shims for C library use
//
// This directory contains code to build rclone as a C library and the
// shims for accessing rclone from C.
//
// The shims are a thin wrapper over the rclone RPC.
//
// Build a shared library like this:
//
//	go build --buildmode=c-shared -o librclone.so github.com/wizdom13/SecureCloudEngine/librclone
//
// Build a static library like this:
//
//	go build --buildmode=c-archive -o librclone.a github.com/wizdom13/SecureCloudEngine/librclone
//
// Both the above commands will also generate `librclone.h` which should
// be `#include`d in `C` programs wishing to use the library.
//
// The library will depend on `libdl` and `libpthread`.
package main

/*
#include <stdlib.h>

struct SecureCloudEngineRPCResult {
	char*	Output;
	int	Status;
};
*/
import "C"

import (
	"unsafe"

	"github.com/wizdom13/SecureCloudEngine/librclone/librclone"

	_ "github.com/wizdom13/SecureCloudEngine/backend/all"   // import all backends
	_ "github.com/wizdom13/SecureCloudEngine/cmd/cmount"    // import cmount
	_ "github.com/wizdom13/SecureCloudEngine/cmd/mount"     // import mount
	_ "github.com/wizdom13/SecureCloudEngine/cmd/mount2"    // import mount2
	_ "github.com/wizdom13/SecureCloudEngine/fs/operations" // import operations/* rc commands
	_ "github.com/wizdom13/SecureCloudEngine/fs/sync"       // import sync/*
	_ "github.com/wizdom13/SecureCloudEngine/lib/plugin"    // import plugins
)

// SecureCloudEngineInitialize initializes rclone as a library
//
//export SecureCloudEngineInitialize
func SecureCloudEngineInitialize() {
	librclone.Initialize()
}

// SecureCloudEngineFinalize finalizes the library
//
//export SecureCloudEngineFinalize
func SecureCloudEngineFinalize() {
	librclone.Finalize()
}

// SecureCloudEngineRPCResult is returned from SecureCloudEngineRPC
//
//	Output will be returned as a serialized JSON object
//	Status is a HTTP status return (200=OK anything else fail)
type SecureCloudEngineRPCResult struct { //nolint:deadcode
	Output *C.char
	Status C.int
}

// SecureCloudEngineRPC does a single RPC call. The inputs are (method, input)
// and the output is (output, status). This is an exported interface
// to the rclone API as described in https://rclone.org/rc/
//
//	method is a string, eg "operations/list"
//	input should be a string with a serialized JSON object
//	result.Output will be returned as a string with a serialized JSON object
//	result.Status is a HTTP status return (200=OK anything else fail)
//
// All strings are UTF-8 encoded, on all platforms.
//
// Caller is responsible for freeing the memory for result.Output
// (see SecureCloudEngineFreeString), result itself is passed on the stack.
//
//export SecureCloudEngineRPC
func SecureCloudEngineRPC(method *C.char, input *C.char) (result C.struct_SecureCloudEngineRPCResult) { //nolint:golint
	output, status := librclone.RPC(C.GoString(method), C.GoString(input))
	result.Output = C.CString(output)
	result.Status = C.int(status)
	return result
}

// SecureCloudEngineFreeString may be used to free the string returned by SecureCloudEngineRPC
//
// If the caller has access to the C standard library, the free function can
// normally be called directly instead. In some cases the caller uses a
// runtime library which is not compatible, and then this function can be
// used to release the memory with the same library that allocated it.
//
//export SecureCloudEngineFreeString
func SecureCloudEngineFreeString(str *C.char) {
	C.free(unsafe.Pointer(str))
}

// do nothing here - necessary for building into a C library
func main() {}
