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
	ErrUnavailable       = errors.New("the server is unavailable")
	ErrUnknown           = errors.New("unknown error")
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.SentenceGenClient
	logger *zap.SugaredLogger
	addr   string
}

func NewClient(addr string, logger *zap.SugaredLogger) (*Client, error) {
	logger.Infow("creating gRPC client", "addr", addr)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		logger.Errorw("failed to create gRPC connection", "addr", addr, "err", err)
		return nil, err
	}

	client := pb.NewSentenceGenClient(conn)

	logger.Infow("gRPC client created", "addr", addr)
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
	c.logger.Debugw("calling GenerateSentence", "word", req.Word, "wordLang", req.WordLanguage, "translationLang", req.TranslationLanguage, "includeAudio", req.IncludeAudio)
	resp, err := c.client.GenerateSentence(ctx, &pb.GenerateSentenceRequest{
		Word:                req.Word,
		WordLanguage:        req.WordLanguage,
		TranslationLanguage: req.TranslationLanguage,
		TranslationHint:     req.TranslationHint,
		IncludeAudio:        req.IncludeAudio,
		VoiceGender:         pb.Gender(req.VoiceGender),
	})

	if err != nil {
		mapped := handleErr(err)
		c.logger.Errorw("GenerateSentence RPC failed", "word", req.Word, "err", err, "mappedErr", mapped)
		return nil, mapped
	}
	c.logger.Debugw("GenerateSentence RPC succeeded", "word", req.Word)
	return &GenerateSentenceResponse{
		OriginalSentence:   resp.OriginalSentence,
		TranslatedSentence: resp.TranslatedSentence,
		Audio:              resp.Audio.Data,
	}, nil
}

func (c *Client) Translate(ctx context.Context, req *TranslateRequest) (*TranslateResponse, error) {
	c.logger.Debugw("calling Translate", "word", req.Word, "fromLang", req.FromLanguage, "toLang", req.ToLanguage, "includeAudio", req.IncludeAudio)
	resp, err := c.client.Translate(ctx, &pb.TranslateRequest{
		Word:            req.Word,
		FromLanguage:    req.FromLanguage,
		ToLanguage:      req.ToLanguage,
		TranslationHint: req.TranslationHint,
		IncludeAudio:    req.IncludeAudio,
		VoiceGender:     pb.Gender(req.VoiceGender),
	})

	if err != nil {
		mapped := handleErr(err)
		c.logger.Errorw("Translate RPC failed", "word", req.Word, "err", err, "mappedErr", mapped)
		return nil, mapped
	}

	c.logger.Debugw("Translate RPC succeeded", "word", req.Word)
	return &TranslateResponse{
		Translation: resp.Translation,
		Audio:       resp.Audio.Data,
	}, nil
}

func (c *Client) GenerateDefinition(ctx context.Context, req *GenerateDefinitionRequest) (*GenerateDefinitionResponse, error) {
	c.logger.Debugw("calling GenerateDefinition", "word", req.Word, "lang", req.Language, "includeAudio", req.IncludeAudio)
	resp, err := c.client.GenerateDefinition(ctx, &pb.GenerateDefinitionRequest{
		Word:           req.Word,
		Language:       req.Language,
		DefinitionHint: req.DefinitionHint,
		IncludeAudio:   req.IncludeAudio,
		VoiceGender:    pb.Gender(req.VoiceGender),
	})
	if err != nil {
		mapped := handleErr(err)
		c.logger.Errorw("GenerateDefinition RPC failed", "word", req.Word, "err", err, "mappedErr", mapped)
		return nil, mapped
	}

	c.logger.Debugw("GenerateDefinition RPC succeeded", "word", req.Word)
	return &GenerateDefinitionResponse{
		Definition: resp.Definition,
		Audio:      resp.Audio.Data,
	}, nil
}

func handleErr(err error) error {
	statusCode := status.Convert(err)
	switch statusCode.Code() {
	case codes.InvalidArgument:
		return ErrInvalidArgument
	case codes.DeadlineExceeded:
		return ErrDeadlineExceeded
	case codes.ResourceExhausted:
		return ErrResourceExhausted
	case codes.Internal:
		return ErrInternalServer
	case codes.Unavailable:
		return ErrUnavailable
	default:
		return ErrUnknown
	}
}
