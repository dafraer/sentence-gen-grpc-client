package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	logFileName = "sengen.log"
	appName     = "Sengen"
)

// GetLogFilePath returns path to which logs will be saved depending on the OS
func GetLogFilePath() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return getMacLogFilePath()
	case "linux":
		return getLinuxLogFilePath()
	case "windows":
		return getWindowsLogFilePath()
	}
	panic("platform does not support os " + runtime.GOOS)
}
func getWindowsLogFilePath() (string, error) {
	//Logs go under %LOCALAPPDATA%\<AppName>
	local := os.Getenv("LOCALAPPDATA")

	if local == "" {
		return "", fmt.Errorf("LOCALAPPDATA environment variable not set")
	}

	// Ensure directory exists
	base := filepath.Join(local, appName, "Logs")
	if err := ensurePath(base); err != nil {
		return "", err
	}
	return filepath.Join(base, logFileName), nil
}

func getMacLogFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// Logs go in ~/Library/Logs/<AppName>

	//Ensure directory exists
	base := filepath.Join(home, "Library", "Logs", appName)
	if err := ensurePath(base); err != nil {
		return "", err
	}
	return filepath.Join(base, logFileName), nil
}

func getLinuxLogFilePath() (string, error) {
	// XDG Base Directory Specification:
	// Logs live in $XDG_STATE_HOME (default ~/.local/state)
	stateHome := os.Getenv("XDG_STATE_HOME")
	if stateHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		stateHome = filepath.Join(home, ".local", "state")
	}

	// Ensure directory exists
	base := filepath.Join(stateHome, appName, "logs")
	if err := ensurePath(base); err != nil {
		return "", err
	}
	return filepath.Join(base, logFileName), nil
}

// ensurePath create directory if it does not exist
func ensurePath(dir string) error {
	return os.MkdirAll(dir, 0o755)
}
