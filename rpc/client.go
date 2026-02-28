package rpc

import (
	"context"
	"errors"

	pb "github.com/dafraer/sentence-gen-grpc-client/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidArgument   = errors.New("invalid argument")
	ErrInternalServer    = errors.New("internal server error")
	ErrResourceExhausted = errors.New("resource exhausted")
	ErrDeadlineExceeded  = errors.New("deadline exceeded")
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.SentenceGenClient
	logger *zap.SugaredLogger
	addr   string
}

func NewClient(addr string, logger *zap.SugaredLogger) (*Client, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, err
	}

	client := pb.NewSentenceGenClient(conn)

	return &Client{
		conn:   conn,
		client: client,
		logger: logger,
		addr:   addr,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GenerateSentence(ctx context.Context, req *GenerateSentenceRequest) (*GenerateSentenceResponse, error) {
	resp, err := c.client.GenerateSentence(ctx, &pb.GenerateSentenceRequest{
		Word:                req.Word,
		WordLanguage:        req.WordLanguage,
		TranslationLanguage: req.TranslationLanguage,
		TranslationHint:     req.TranslationHint,
		IncludeAudio:        req.IncludeAudio,
		VoiceGender:         pb.Gender(req.VoiceGender),
	})
	if err != nil {
		statusCode := status.Convert(err)
		switch statusCode.Code() {
		case codes.InvalidArgument:
			return nil, ErrInvalidArgument
		case codes.DeadlineExceeded:
			return nil, ErrDeadlineExceeded
		case codes.ResourceExhausted:
			return nil, ErrResourceExhausted
		case codes.Internal:
			return nil, ErrInternalServer
		}
	}
	c.logger.Infow("Successfully generated sentence", "result", GenerateSentenceResponse{
		OriginalSentence:   resp.OriginalSentence,
		TranslatedSentence: resp.TranslatedSentence,
		Audio:              resp.Audio.Data,
	})
	return &GenerateSentenceResponse{
		OriginalSentence:   resp.OriginalSentence,
		TranslatedSentence: resp.TranslatedSentence,
		Audio:              resp.Audio.Data,
	}, nil
}

func (c *Client) Translate(ctx context.Context, req *TranslateRequest) (*TranslateResponse, error) {
	return nil, nil
}

func (c *Client) GenerateDefinition(ctx context.Context, req *GenerateDefinitionRequest) (*GenerateDefinitionResponse, error) {
	return nil, nil
}
