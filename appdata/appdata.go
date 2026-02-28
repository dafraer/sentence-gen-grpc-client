package appdata

import (
	"encoding/json"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

const (
	stateFile    = "state.json"
	settingsFile = "settings.json"
)

type AppData struct {
	logger   *zap.SugaredLogger
	stateDir string
	audioDir string
}

func NewAppData(logger *zap.SugaredLogger, stateDir, audioDir string) *AppData {
	return &AppData{
		logger:   logger,
		stateDir: stateDir,
		audioDir: audioDir,
	}
}

func (a *AppData) SaveState(state State) error {
	return writeJSON(filepath.Join(a.stateDir, stateFile), state)
}

func (a *AppData) SaveSettings(settings Settings) error {
	return writeJSON(filepath.Join(a.stateDir, settingsFile), settings)
}

func (a *AppData) SaveAudio(data []byte, filename string) error {
	return os.WriteFile(filepath.Join(a.audioDir, filename), data, 0644)
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
