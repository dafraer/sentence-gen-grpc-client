package anki

import "go.uber.org/zap"

type Client struct {
	logger          *zap.SugaredLogger
	ankiConnectAddr string
}

func NewClient(logger *zap.SugaredLogger, ankiConnectAddr string) *Client {
	return &Client{
		logger:          logger,
		ankiConnectAddr: ankiConnectAddr,
	}
}
