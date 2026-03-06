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
	//Check for the debug CLI argument
	debug := len(os.Args) > 1 && os.Args[1] == "debug"

	//Load the config
	cfg, err := config.Load(debug)
	if err != nil {
		panic(err)
	}

	//Create sugared logger
	opt := zap.NewDevelopmentConfig()
	opt.OutputPaths = []string{cfg.LogPath}
	logger, err := opt.Build()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()

	defer func(logger *zap.SugaredLogger) {
		if err := logger.Sync(); err != nil {
			panic(err)
		}
	}(sugar)

	//Create anki connect client
	ankiClient := anki.NewClient(sugar, cfg.AnkiConnectAddr)

	//Create grpc client
	grpcClient, err := rpc.NewClient(cfg.ServerAddr, sugar)
	if err != nil {
		panic(err)
	}

	//Get the text
	txt := text.NewText()

	//Create domain layer
	appCore := core.New(sugar, grpcClient, ankiClient)

	sugar.Infow("starting Sengen", "serverAddr", cfg.ServerAddr, "ankiAddr", cfg.AnkiConnectAddr, "logPath", cfg.LogPath)

	//Create GUI
	appGUI := gui.New(sugar, appCore, txt)

	//Run the app
	appGUI.Run()
}
