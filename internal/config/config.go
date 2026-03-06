package config

import (
	"github.com/dafraer/sentence-gen-grpc-client/internal/platform"
)

const (
	ankiConnectAddr = "localhost:8765"
	serverAddr      = "localhost:50051"
)

type Config struct {
	LogPath         string
	ServerAddr      string
	AnkiConnectAddr string
}

// Load loads config and panics if it encounters an error
func Load(debug bool) (*Config, error) {
	//Get path for the log file on the user's OS
	logPath, err := platform.GetLogFilePath()
	if err != nil {
		return nil, err
	}

	//If debug mode is on output logs to stderr
	if debug {
		logPath = "stderr"
	}

	return &Config{
		LogPath:         logPath,
		AnkiConnectAddr: ankiConnectAddr,
		ServerAddr:      serverAddr,
	}, nil
}
