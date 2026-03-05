package config

import (
	"github.com/dafraer/sentence-gen-grpc-client/platform"
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
	logPath, err := platform.GetLogFilePath()
	if err != nil {
		return nil, err
	}

	if debug {
		logPath = "stderr"
	}
	return &Config{
		LogPath:         logPath,
		AnkiConnectAddr: ankiConnectAddr,
		ServerAddr:      serverAddr,
	}, nil
}
