package platform

import (
	"os"
	"runtime"
)

const (
	logFileName   = "sengen.log"
	stateFileName = "state.json"
	appName       = "Sengen"
)

type Platform interface {
	GetStateFilePath() (string, error)
	GetLogFilePath() (string, error)
}

// New returns platform or panics if platform is not supported
func New() Platform {
	switch runtime.GOOS {
	case "darwin":
		return &darwin{}
	case "linux":
		return &linux{}
	case "windows":
		return &windows{}
	}
	panic("platform does not support os " + runtime.GOOS)
}

// ensurePath create directory if it does not exist
func ensurePath(dir string) error {
	return os.MkdirAll(dir, 0o755)
}
