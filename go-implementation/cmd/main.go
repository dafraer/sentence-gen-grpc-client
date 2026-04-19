package main

import (
	"os"

	"github.com/dafraer/sentence-gen-grpc-client/internal/anki"
	"github.com/dafraer/sentence-gen-grpc-client/internal/config"
	"github.com/dafraer/sentence-gen-grpc-client/internal/core"
	"github.com/dafraer/sentence-gen-grpc-client/internal/gui"
	"github.com/dafraer/sentence-gen-grpc-client/internal/rpc"
	"github.com/dafraer/sentence-gen-grpc-client/internal/text"
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

	defer func(grpcClient *rpc.Client) {
		if err := grpcClient.Close(); err != nil {
			panic(err)
		}
	}(grpcClient)

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
