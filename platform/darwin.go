package platform

import (
	"os"
	"path/filepath"
)

type darwin struct{}

func (d *darwin) GetLogFilePath() (string, error) {
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

func (d *darwin) GetStateFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// State data goes in ~/Library/Application Support/<AppName>.
	base := filepath.Join(home, "Library", "Application Support", appName)
	if err := ensurePath(base); err != nil {
		return "", err
	}
	return filepath.Join(base, stateFileName), nil
}
