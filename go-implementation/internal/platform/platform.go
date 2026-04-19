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
	return "", fmt.Errorf("platform does not support os " + runtime.GOOS)
}

// getWindowsLogFilePath returns log file path for windows
// On Windows logs are stored in %LOCALAPPDATA%\<AppName>
// It also checks that the directory exists and creates it if it doesn't
func getWindowsLogFilePath() (string, error) {
	//Logs go under %LOCALAPPDATA%\<AppName>
	local := os.Getenv("LOCALAPPDATA")

	//Check that env variable is set
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

// getMacLogFilePath returns log file path for macOS
// On macOS logs are stored in ~/Library/Logs/<AppName>
// It also checks that the directory exists and creates it if it doesn't
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

// getLinuxLogFilePath returns log file path for linux
// On Linux logs are usually stored in $XDG_STATE_HOME (default ~/.local/state)
// It also checks that the directory exists and creates it if it doesn't
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
