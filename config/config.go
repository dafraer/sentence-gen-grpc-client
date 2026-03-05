package config

import (
	"github.com/dafraer/sentence-gen-grpc-client/platform"
)

const (
	ankiConnectAddr = "localhost:8765"
	serverAddr      = "localhost:50051"
)

type Config struct {
	StatePath       string
	LogPath         string
	ServerAddr      string
	AnkiConnectAddr string
}

// Load loads config and panics if it encounters an error
func Load(debug bool) (*Config, error) {
	pl := platform.New()
	statePath, err := pl.GetStateFilePath()
	if err != nil {
		return nil, err
	}

	logPath, err := pl.GetLogFilePath()
	if err != nil {
		return nil, err
	}

	if debug {
		statePath = "./state.json"
		logPath = "stderr"
	}
	return &Config{
		StatePath:       statePath,
		LogPath:         logPath,
		AnkiConnectAddr: ankiConnectAddr,
		ServerAddr:      serverAddr,
	}, nil
}
