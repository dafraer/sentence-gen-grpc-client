package gui

import (
	"github.com/dafraer/sentence-gen-grpc-client/core"
	"go.uber.org/zap"
)

type GUI struct {
	logger *zap.SugaredLogger
	core   *core.Core
}

func NewGUI(logger *zap.SugaredLogger, core *core.Core) *GUI {
	return &GUI{
		logger: logger,
		core:   core,
	}
}
