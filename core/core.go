package core

import (
	"context"
	"fmt"
	"time"

	"github.com/dafraer/sentence-gen-grpc-client/anki"
	"github.com/dafraer/sentence-gen-grpc-client/rpc"
	"go.uber.org/zap"
)

const audioFormat = "%s_%d.wav"

type Core struct {
	logger     *zap.SugaredLogger
	grpcClient *rpc.Client
	ankiClient *anki.Client
}

func New(logger *zap.SugaredLogger, grpcClient *rpc.Client, ankiClient *anki.Client) *Core {
	return &Core{
		logger:     logger,
		grpcClient: grpcClient,
		ankiClient: ankiClient,
	}
}

func (c *Core) GetDeckNames(ctx context.Context) ([]string, error) {
	return c.ankiClient.GetDeckNames(ctx)
}

func (c *Core) GenerateSentence(ctx context.Context, req *GenerateSentenceRequest) (*GenerateSentenceResponse, error) {
	c.logger.Infow("generating sentence", "word", req.Word, "wordLang", req.WordLanguage, "translationLang", req.TranslationLanguage, "deck", req.DeckName, "includeAudio", req.IncludeAudio)
	resp, err := c.grpcClient.GenerateSentence(ctx, &rpc.GenerateSentenceRequest{
		Word:                req.Word,
		WordLanguage:        req.WordLanguage,
		TranslationLanguage: req.TranslationLanguage,
		TranslationHint:     req.TranslationHint,
		IncludeAudio:        req.IncludeAudio,
		VoiceGender:         toRPCGender(req.AudioGender),
	})
	if err != nil {
		c.logger.Errorw("failed to generate sentence", "word", req.Word, "err", err)
		return nil, err
	}

	var ankiAudio *anki.Audio
	if req.IncludeAudio && len(resp.Audio) > 0 {
		filename := fmt.Sprintf(audioFormat, req.Word, time.Now().Unix())
		c.logger.Debugw("audio received for sentence", "word", req.Word, "bytes", len(resp.Audio), "filename", filename)
		//if err := c.appData.SaveAudio(resp.Audio, filename); err != nil {
		//	c.logger.Errorw("failed to save audio", "err", err)
		//}
		ankiAudio = &anki.Audio{Data: resp.Audio, Filename: filename, Fields: []string{"Front"}}
	}

	if err := c.ankiClient.AddCard(ctx, anki.Note{
		DeckName: req.DeckName,
		CardType: anki.BasicAndReverse,
		Front:    resp.OriginalSentence,
		Back:     resp.TranslatedSentence,
		Audio:    ankiAudio,
	}); err != nil {
		c.logger.Errorw("failed to add sentence card to Anki", "word", req.Word, "deck", req.DeckName, "err", err)
		return nil, err
	}

	c.logger.Infow("sentence card added successfully", "word", req.Word, "deck", req.DeckName)
	return &GenerateSentenceResponse{
		OriginalSentence:   resp.OriginalSentence,
		TranslatedSentence: resp.TranslatedSentence,
		Audio:              resp.Audio,
	}, nil
}

func (c *Core) Translate(ctx context.Context, req *TranslateRequest) (*TranslateResponse, error) {
	c.logger.Infow("translating word", "word", req.Word, "wordLang", req.WordLanguage, "translationLang", req.TranslationLang, "deck", req.DeckName, "includeAudio", req.IncludeAudio)
	resp, err := c.grpcClient.Translate(ctx, &rpc.TranslateRequest{
		Word:            req.Word,
		FromLanguage:    req.WordLanguage,
		ToLanguage:      req.TranslationLang,
		TranslationHint: req.TranslationHint,
		IncludeAudio:    req.IncludeAudio,
		VoiceGender:     toRPCGender(req.AudioGender),
	})
	if err != nil {
		c.logger.Errorw("failed to translate word", "word", req.Word, "err", err)
		return nil, err
	}

	var ankiAudio *anki.Audio
	if req.IncludeAudio && len(resp.Audio) > 0 {
		filename := fmt.Sprintf(audioFormat, req.Word, time.Now().Unix())
		c.logger.Debugw("audio received for translation", "word", req.Word, "bytes", len(resp.Audio), "filename", filename)
		ankiAudio = &anki.Audio{Data: resp.Audio, Filename: filename, Fields: []string{"Front"}}
	}

	if err := c.ankiClient.AddCard(ctx, anki.Note{
		DeckName: req.DeckName,
		CardType: anki.BasicAndReverse,
		Front:    req.Word,
		Back:     resp.Translation,
		Audio:    ankiAudio,
	}); err != nil {
		c.logger.Errorw("failed to add translation card to Anki", "word", req.Word, "deck", req.DeckName, "err", err)
		return nil, err
	}

	c.logger.Infow("translation card added successfully", "word", req.Word, "deck", req.DeckName)
	return &TranslateResponse{
		Translation: resp.Translation,
		Audio:       resp.Audio,
	}, nil
}

func (c *Core) GenerateDefinition(ctx context.Context, req *GenerateDefinitionRequest) (*GenerateDefinitionResponse, error) {
	c.logger.Infow("generating definition", "word", req.Word, "lang", req.Language, "deck", req.DeckName, "includeAudio", req.IncludeAudio)
	resp, err := c.grpcClient.GenerateDefinition(ctx, &rpc.GenerateDefinitionRequest{
		Word:           req.Word,
		Language:       req.Language,
		DefinitionHint: req.DefinitionHint,
		IncludeAudio:   req.IncludeAudio,
		VoiceGender:    toRPCGender(req.AudioGender),
	})
	if err != nil {
		c.logger.Errorw("failed to generate definition", "word", req.Word, "err", err)
		return nil, err
	}

	var ankiAudio *anki.Audio
	if req.IncludeAudio && len(resp.Audio) > 0 {
		filename := fmt.Sprintf(audioFormat, req.Word, time.Now().Unix())
		c.logger.Debugw("audio received for definition", "word", req.Word, "bytes", len(resp.Audio), "filename", filename)
		ankiAudio = &anki.Audio{Data: resp.Audio, Filename: filename, Fields: []string{"Front"}}
	}

	if err := c.ankiClient.AddCard(ctx, anki.Note{
		DeckName: req.DeckName,
		CardType: anki.Basic,
		Front:    req.Word,
		Back:     resp.Definition,
		Audio:    ankiAudio,
	}); err != nil {
		c.logger.Errorw("failed to add definition card to Anki", "word", req.Word, "deck", req.DeckName, "err", err)
		return nil, err
	}

	c.logger.Infow("definition card added successfully", "word", req.Word, "deck", req.DeckName)
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
