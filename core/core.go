package core

import (
	"github.com/dafraer/sentence-gen-grpc-client/anki"
	"github.com/dafraer/sentence-gen-grpc-client/rpc"
	"go.uber.org/zap"
)

type Core struct {
	logger     *zap.SugaredLogger
	grpcClient *rpc.Client
	ankiClient *anki.Client
}

func New(logger *zap.SugaredLogger, grpcClient *rpc.Client, ankiClient *anki.Client) *Core {
	return &Core{
		logger:     logger,
		grpcClient: grpcClient,
		ankiClient: ankiClient,
	}
}
