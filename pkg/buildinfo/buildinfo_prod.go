//go:build prod
// +build prod

package buildinfo

const (
	devMode = false
)

var (
	version   = ""
	build     = ""
	buildTime = ""
)
