//go:build !prod
// +build !prod

package buildinfo

const (
	devMode   = true
	version   = "DEVELOPMENT Version"
	build     = "DEVELOPMENT Build"
	buildTime = ""
)
