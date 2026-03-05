package platform

import (
	"fmt"
	"os"
	"path/filepath"
)

type windows struct{}

func (w *windows) GetLogFilePath() (string, error) {
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

func (w *windows) GetStateFilePath() (string, error) {
	//State goes under %LOCALAPPDATA%\<AppName>
	local := os.Getenv("LOCALAPPDATA")

	if local == "" {
		return "", fmt.Errorf("LOCALAPPDATA environment variable not set")
	}

	// Ensure directory exists
	base := filepath.Join(local, appName, "Data")
	if err := ensurePath(base); err != nil {
		return "", err
	}
	return filepath.Join(base, stateFileName), nil
}
