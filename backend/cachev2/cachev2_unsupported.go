// Build for cachev2 for unsupported platforms to stop go complaining
// about "no buildable Go source files "

//go:build plan9 || js

// Package cachev2 implements a virtual provider to cache existing remotes.
package cachev2
