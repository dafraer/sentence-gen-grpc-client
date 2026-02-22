package rpc

import "go.uber.org/zap"

type Client struct {
	logger *zap.SugaredLogger
	addr   string
}

func NewClient(addr string, logger *zap.SugaredLogger) *Client {
	return &Client{
		logger: logger,
		addr:   addr,
	}
}
