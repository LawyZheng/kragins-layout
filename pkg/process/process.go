package process

import (
	"os"
	"path/filepath"
)

type IProcessPath interface {
	ExecutablePath() string
}

var Path IProcessPath = &defaultImpl{}

type defaultImpl struct {
}

func (d *defaultImpl) ExecutablePath() string {
	if ex, err := os.Executable(); err == nil {
		return filepath.Dir(ex)
	}
	return "."
}

// ExecutablePath to get executable binary path
func ExecutablePath() string {
	return Path.ExecutablePath()
}

// GetPID to get process pid
func GetPID() int {
	return os.Getpid()
}
