package main

import (
	"github.com/dafraer/sentence-gen-grpc-client/anki"
	"github.com/dafraer/sentence-gen-grpc-client/core"
	"github.com/dafraer/sentence-gen-grpc-client/gui"
	"github.com/dafraer/sentence-gen-grpc-client/rpc"
	"github.com/dafraer/sentence-gen-grpc-client/text"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()

	defer func(logger *zap.SugaredLogger) {
		if err := logger.Sync(); err != nil {
			panic(err)
		}
	}(sugar)

	ankiClient := anki.NewClient(sugar, "localhost:8765")

	grpcClient := rpc.NewClient("localhost:50051", sugar)

	txt := text.NewText("ru")

	appCore := core.New(sugar, grpcClient, ankiClient)

	appGUI := gui.New(sugar, appCore, txt)

	appGUI.Run()

}
