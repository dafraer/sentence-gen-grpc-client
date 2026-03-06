package rpc

import (
	"context"
	"errors"

	"github.com/dafraer/sentence-gen-grpc-client/internal/proto"
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
	client proto.SentenceGenClient
	logger *zap.SugaredLogger
}

// NewClient creates new grpc Client
func NewClient(addr string, logger *zap.SugaredLogger) (*Client, error) {
	logger.Infow("creating gRPC client", "addr", addr)
	opts := []grpc.DialOption{
		//TODO:Change to secure
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	//Create new grpc client
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		logger.Errorw("failed to create gRPC connection", "addr", addr, "err", err)
		return nil, err
	}

	//Create new SentenceGen client
	client := proto.NewSentenceGenClient(conn)

	logger.Infow("gRPC client created", "addr", addr)
	return &Client{
		conn:   conn,
		client: client,
		logger: logger,
	}, nil
}

// Close closes the connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// GenerateSentence generates sentences and audio for specified word
func (c *Client) GenerateSentence(ctx context.Context, req *GenerateSentenceRequest) (*GenerateSentenceResponse, error) {
	c.logger.Debugw("calling GenerateSentence", "word", req.Word, "wordLang", req.WordLanguage, "translationLang", req.TranslationLanguage, "includeAudio", req.IncludeAudio)

	//Call grpc function
	resp, err := c.client.GenerateSentence(ctx, &proto.GenerateSentenceRequest{
		Word:                req.Word,
		WordLanguage:        req.WordLanguage,
		TranslationLanguage: req.TranslationLanguage,
		TranslationHint:     req.TranslationHint,
		IncludeAudio:        req.IncludeAudio,
		VoiceGender:         proto.Gender(req.VoiceGender),
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

// Translate generates translation and audio for specified word
func (c *Client) Translate(ctx context.Context, req *TranslateRequest) (*TranslateResponse, error) {
	c.logger.Debugw("calling Translate", "word", req.Word, "fromLang", req.FromLanguage, "toLang", req.ToLanguage, "includeAudio", req.IncludeAudio)
	resp, err := c.client.Translate(ctx, &proto.TranslateRequest{
		Word:            req.Word,
		FromLanguage:    req.FromLanguage,
		ToLanguage:      req.ToLanguage,
		TranslationHint: req.TranslationHint,
		IncludeAudio:    req.IncludeAudio,
		VoiceGender:     proto.Gender(req.VoiceGender),
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

// GenerateDefinition generates definition and audio for specified word
func (c *Client) GenerateDefinition(ctx context.Context, req *GenerateDefinitionRequest) (*GenerateDefinitionResponse, error) {
	c.logger.Debugw("calling GenerateDefinition", "word", req.Word, "lang", req.Language, "includeAudio", req.IncludeAudio)
	resp, err := c.client.GenerateDefinition(ctx, &proto.GenerateDefinitionRequest{
		Word:           req.Word,
		Language:       req.Language,
		DefinitionHint: req.DefinitionHint,
		IncludeAudio:   req.IncludeAudio,
		VoiceGender:    proto.Gender(req.VoiceGender),
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

// handleErr maps gRPC status code to a local error
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
