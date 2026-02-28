package core

import (
	"context"
	"fmt"
	"time"

	"github.com/dafraer/sentence-gen-grpc-client/anki"
	"github.com/dafraer/sentence-gen-grpc-client/appdata"
	"github.com/dafraer/sentence-gen-grpc-client/rpc"
	"go.uber.org/zap"
)

const audioFormat = "%s_%d.wav"

type Core struct {
	logger     *zap.SugaredLogger
	grpcClient *rpc.Client
	ankiClient *anki.Client
	appData    *appdata.AppData
}

func New(logger *zap.SugaredLogger, grpcClient *rpc.Client, ankiClient *anki.Client, appData *appdata.AppData) *Core {
	return &Core{
		logger:     logger,
		grpcClient: grpcClient,
		ankiClient: ankiClient,
		appData:    appData,
	}
}

func (c *Core) GenerateSentence(ctx context.Context, req *GenerateSentenceRequest) (*GenerateSentenceResponse, error) {
	resp, err := c.grpcClient.GenerateSentence(ctx, &rpc.GenerateSentenceRequest{
		Word:                req.Word,
		WordLanguage:        req.WordLanguage,
		TranslationLanguage: req.TranslationLanguage,
		TranslationHint:     req.TranslationHint,
		IncludeAudio:        req.IncludeAudio,
		VoiceGender:         toRPCGender(req.AudioGender),
	})
	if err != nil {
		return nil, err
	}

	var ankiAudio *anki.Audio
	if req.IncludeAudio && len(resp.Audio) > 0 {
		filename := fmt.Sprintf(audioFormat, req.Word, time.Now().Unix())
		if err := c.appData.SaveAudio(resp.Audio, filename); err != nil {
			c.logger.Errorw("failed to save audio", "err", err)
		}
		ankiAudio = &anki.Audio{Data: resp.Audio, Filename: filename, Fields: []string{"Front"}}
	}

	if err := c.ankiClient.AddCard(ctx, anki.Note{
		DeckName: req.DeckName,
		CardType: anki.BasicAndReverse,
		Front:    resp.OriginalSentence,
		Back:     resp.TranslatedSentence,
		Audio:    ankiAudio,
	}); err != nil {
		return nil, err
	}

	return &GenerateSentenceResponse{
		OriginalSentence:   resp.OriginalSentence,
		TranslatedSentence: resp.TranslatedSentence,
		Audio:              resp.Audio,
	}, nil
}

func (c *Core) Translate(ctx context.Context, req *TranslateRequest) (*TranslateResponse, error) {
	resp, err := c.grpcClient.Translate(ctx, &rpc.TranslateRequest{
		Word:            req.Word,
		FromLanguage:    req.WordLanguage,
		ToLanguage:      req.TranslationLang,
		TranslationHint: req.TranslationHint,
		IncludeAudio:    req.IncludeAudio,
		VoiceGender:     toRPCGender(req.AudioGender),
	})
	if err != nil {
		return nil, err
	}

	var ankiAudio *anki.Audio
	if req.IncludeAudio && len(resp.Audio) > 0 {
		filename := fmt.Sprintf(audioFormat, req.Word, time.Now().Unix())
		if err := c.appData.SaveAudio(resp.Audio, filename); err != nil {
			c.logger.Errorw("failed to save audio", "err", err)
		}
		ankiAudio = &anki.Audio{Data: resp.Audio, Filename: filename, Fields: []string{"Front"}}
	}

	if err := c.ankiClient.AddCard(ctx, anki.Note{
		DeckName: req.DeckName,
		CardType: anki.BasicAndReverse,
		Front:    req.Word,
		Back:     resp.Translation,
		Audio:    ankiAudio,
	}); err != nil {
		return nil, err
	}

	return &TranslateResponse{
		Translation: resp.Translation,
		Audio:       resp.Audio,
	}, nil
}

func (c *Core) GenerateDefinition(ctx context.Context, req *GenerateDefinitionRequest) (*GenerateDefinitionResponse, error) {
	resp, err := c.grpcClient.GenerateDefinition(ctx, &rpc.GenerateDefinitionRequest{
		Word:           req.Word,
		Language:       req.Language,
		DefinitionHint: req.DefinitionHint,
		IncludeAudio:   req.IncludeAudio,
		VoiceGender:    toRPCGender(req.AudioGender),
	})
	if err != nil {
		return nil, err
	}

	var ankiAudio *anki.Audio
	if req.IncludeAudio && len(resp.Audio) > 0 {
		filename := fmt.Sprintf(audioFormat, req.Word, time.Now().Unix())
		if err := c.appData.SaveAudio(resp.Audio, filename); err != nil {
			c.logger.Errorw("failed to save audio", "err", err)
		}
		ankiAudio = &anki.Audio{Data: resp.Audio, Filename: filename, Fields: []string{"Front"}}
	}

	if err := c.ankiClient.AddCard(ctx, anki.Note{
		DeckName: req.DeckName,
		CardType: anki.Basic,
		Front:    req.Word,
		Back:     resp.Definition,
		Audio:    ankiAudio,
	}); err != nil {
		return nil, err
	}

	return &GenerateDefinitionResponse{
		Definition: resp.Definition,
		Audio:      resp.Audio,
	}, nil
}

func toRPCGender(gender string) rpc.Gender {
	if gender == "Male" {
		return rpc.Male
	}
	return rpc.Female
}
