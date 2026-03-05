package appdata

import (
	"encoding/json"
	"os"

	"go.uber.org/zap"
)

const ()

type AppData struct {
	logger        *zap.SugaredLogger
	stateFilePath string
}

func NewAppData(logger *zap.SugaredLogger, stateFilePath string) *AppData {
	return &AppData{
		logger:        logger,
		stateFilePath: stateFilePath,
	}
}

func (a *AppData) SaveState(state State) error {
	return writeJSON(a.stateFilePath, state)
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
