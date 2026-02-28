package main

import (
	"github.com/dafraer/sentence-gen-grpc-client/anki"
	"github.com/dafraer/sentence-gen-grpc-client/appdata"
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

	grpcClient, err := rpc.NewClient("localhost:50051", sugar)
	if err != nil {
		panic(err)
	}

	appData := appdata.NewAppData(sugar, ".", ".")

	txt := text.NewText("en")

	appCore := core.New(sugar, grpcClient, ankiClient, appData)

	appGUI := gui.New(sugar, appCore, txt)

	appGUI.Run()

}
