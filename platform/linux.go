package platform

import (
	"os"
	"path/filepath"
)

type linux struct{}

func (l *linux) GetLogFilePath() (string, error) {
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

func (l *linux) GetStateFilePath() (string, error) {
	// XDG Base Directory Specification:
	// State lives in $XDG_STATE_HOME (default ~/.local/state)
	stateHome := os.Getenv("XDG_STATE_HOME")
	if stateHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		stateHome = filepath.Join(home, ".local", "state")
	}

	// Ensure directory exists
	base := filepath.Join(stateHome, appName)
	if err := ensurePath(base); err != nil {
		return "", err
	}
	return filepath.Join(base, stateFileName), nil
}
