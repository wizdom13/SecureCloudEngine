//go:build noselfupdate

package selfupdate

import (
	"/lib/buildinfo"
)

func init() {
	buildinfo.Tags = append(buildinfo.Tags, "noselfupdate")
}
