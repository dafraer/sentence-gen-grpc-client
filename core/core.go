package core

import (
	"context"
	"fmt"
	"time"

	"github.com/dafraer/sentence-gen-grpc-client/anki"
	"github.com/dafraer/sentence-gen-grpc-client/rpc"
	"go.uber.org/zap"
)

const audioFormat = "%s_%d.wav" //Format for .wav audio

type Core struct {
	logger     *zap.SugaredLogger
	grpcClient *rpc.Client
	ankiClient *anki.Client
}

// New creates new app core
func New(logger *zap.SugaredLogger, grpcClient *rpc.Client, ankiClient *anki.Client) *Core {
	return &Core{
		logger:     logger,
		grpcClient: grpcClient,
		ankiClient: ankiClient,
	}
}

// GetDeckNames get deck names from anki package
func (c *Core) GetDeckNames(ctx context.Context) ([]string, error) {
	return c.ankiClient.GetDeckNames(ctx)
}

// GenerateSentence generates sentences and audio to them
func (c *Core) GenerateSentence(ctx context.Context, req *GenerateSentenceRequest) error {
	c.logger.Infow("generating sentence", "word", req.Word, "wordLang", req.WordLanguage, "translationLang", req.TranslationLanguage, "deck", req.DeckName, "includeAudio", req.IncludeAudio)

	//Generate sentences
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
		return err
	}

	//Build audio file for anki
	var ankiAudio *anki.Audio
	if req.IncludeAudio && len(resp.Audio) > 0 {
		filename := fmt.Sprintf(audioFormat, req.Word, time.Now().Unix())
		c.logger.Debugw("audio received for sentence", "word", req.Word, "bytes", len(resp.Audio), "filename", filename)
		ankiAudio = &anki.Audio{Data: resp.Audio, Filename: filename, Fields: []string{"Front"}}
	}

	//Add anki card
	if err := c.ankiClient.AddCard(ctx, anki.Note{
		DeckName: req.DeckName,
		CardType: anki.BasicAndReverse,
		Front:    resp.OriginalSentence,
		Back:     resp.TranslatedSentence,
		Audio:    ankiAudio,
	}); err != nil {
		c.logger.Errorw("failed to add sentence card to Anki", "word", req.Word, "deck", req.DeckName, "err", err)
		return err
	}

	c.logger.Infow("sentence card added successfully", "word", req.Word, "deck", req.DeckName)
	return nil
}

// Translate generates translation and adds it to anki
func (c *Core) Translate(ctx context.Context, req *TranslateRequest) error {
	c.logger.Infow("translating word", "word", req.Word, "wordLang", req.WordLanguage, "translationLang", req.TranslationLang, "deck", req.DeckName, "includeAudio", req.IncludeAudio)
	//Translate the word
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
		return err
	}

	//Build audio file for anki
	var ankiAudio *anki.Audio
	if req.IncludeAudio && len(resp.Audio) > 0 {
		filename := fmt.Sprintf(audioFormat, req.Word, time.Now().Unix())
		c.logger.Debugw("audio received for translation", "word", req.Word, "bytes", len(resp.Audio), "filename", filename)
		ankiAudio = &anki.Audio{Data: resp.Audio, Filename: filename, Fields: []string{"Front"}}
	}

	//Add card to anki
	if err := c.ankiClient.AddCard(ctx, anki.Note{
		DeckName: req.DeckName,
		CardType: anki.BasicAndReverse,
		Front:    req.Word,
		Back:     resp.Translation,
		Audio:    ankiAudio,
	}); err != nil {
		c.logger.Errorw("failed to add translation card to Anki", "word", req.Word, "deck", req.DeckName, "err", err)
		return err
	}

	c.logger.Infow("translation card added successfully", "word", req.Word, "deck", req.DeckName)
	return nil
}

// GenerateDefinition generates definition for a word and adds it to the anki deck
func (c *Core) GenerateDefinition(ctx context.Context, req *GenerateDefinitionRequest) error {
	c.logger.Infow("generating definition", "word", req.Word, "lang", req.Language, "deck", req.DeckName, "includeAudio", req.IncludeAudio)
	//Generate definition
	resp, err := c.grpcClient.GenerateDefinition(ctx, &rpc.GenerateDefinitionRequest{
		Word:           req.Word,
		Language:       req.Language,
		DefinitionHint: req.DefinitionHint,
		IncludeAudio:   req.IncludeAudio,
		VoiceGender:    toRPCGender(req.AudioGender),
	})
	if err != nil {
		c.logger.Errorw("failed to generate definition", "word", req.Word, "err", err)
		return err
	}

	//Build audio file for anki
	var ankiAudio *anki.Audio
	if req.IncludeAudio && len(resp.Audio) > 0 {
		filename := fmt.Sprintf(audioFormat, req.Word, time.Now().Unix())
		c.logger.Debugw("audio received for definition", "word", req.Word, "bytes", len(resp.Audio), "filename", filename)
		ankiAudio = &anki.Audio{Data: resp.Audio, Filename: filename, Fields: []string{"Front"}}
	}

	//Add card to anki
	if err := c.ankiClient.AddCard(ctx, anki.Note{
		DeckName: req.DeckName,
		CardType: anki.Basic,
		Front:    req.Word,
		Back:     resp.Definition,
		Audio:    ankiAudio,
	}); err != nil {
		c.logger.Errorw("failed to add definition card to Anki", "word", req.Word, "deck", req.DeckName, "err", err)
		return err
	}

	c.logger.Infow("definition card added successfully", "word", req.Word, "deck", req.DeckName)
	return nil
}

func toRPCGender(gender string) rpc.Gender {
	if gender == "Male" {
		return rpc.Male
	}
	return rpc.Female
}
