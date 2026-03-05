package main

import (
	"os"

	"github.com/dafraer/sentence-gen-grpc-client/anki"
	"github.com/dafraer/sentence-gen-grpc-client/config"
	"github.com/dafraer/sentence-gen-grpc-client/core"
	"github.com/dafraer/sentence-gen-grpc-client/gui"
	"github.com/dafraer/sentence-gen-grpc-client/rpc"
	"github.com/dafraer/sentence-gen-grpc-client/text"
	"go.uber.org/zap"
)

func main() {
	debug := len(os.Args) > 1 && os.Args[1] == "debug"

	cfg, err := config.Load(debug)
	if err != nil {
		panic(err)
	}

	opt := zap.NewDevelopmentConfig()
	opt.OutputPaths = []string{cfg.LogPath}
	logger, err := opt.Build()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()

	sugar.Infow("starting server")

	defer func(logger *zap.SugaredLogger) {
		if err := logger.Sync(); err != nil {
			panic(err)
		}
	}(sugar)

	ankiClient := anki.NewClient(sugar, cfg.AnkiConnectAddr)

	grpcClient, err := rpc.NewClient(cfg.ServerAddr, sugar)
	if err != nil {
		panic(err)
	}

	txt := text.NewText("en")

	appCore := core.New(sugar, grpcClient, ankiClient)

	appGUI := gui.New(sugar, appCore, txt)

	appGUI.Run()

}
