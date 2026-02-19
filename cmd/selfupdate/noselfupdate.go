//go:build noselfupdate

package selfupdate

import (
	"github.com/wizdom13/SecureCloudEngine/lib/buildinfo"
)

func init() {
	buildinfo.Tags = append(buildinfo.Tags, "noselfupdate")
}
