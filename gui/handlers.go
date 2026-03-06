package gui

import (
	"context"
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/dafraer/sentence-gen-grpc-client/anki"
	"github.com/dafraer/sentence-gen-grpc-client/core"
	"github.com/dafraer/sentence-gen-grpc-client/rpc"
)

var ErrUnknownLanguage = errors.New("unknown language")

func (gui *GUI) showErrorNotification(err error, word string) {
	gui.logger.Errorw("showing error notification", "word", word, "err", err)
	errText := ""
	switch {
	case errors.Is(err, ErrUnknownLanguage):
		errText = gui.text.TextErrUnknownLanguage()
	case errors.Is(err, anki.ErrAnkiError):
		errText = "Anki error"
	case errors.Is(err, rpc.ErrInvalidArgument):
		errText = "Invalid argument"
	case errors.Is(err, rpc.ErrDeadlineExceeded):
		errText = "Server deadline exceeded"
	case errors.Is(err, rpc.ErrInternalServer):
		errText = "Internal server error"
	case errors.Is(err, rpc.ErrResourceExhausted):
		errText = "Quota limit reached, try again tomorrow"
	case errors.Is(err, rpc.ErrUnavailable):
		errText = "Server is unavailable"
	case errors.Is(err, rpc.ErrUnknown):
		errText = "Unknown error"
	}
	fyne.Do(func() {
		gui.app.SendNotification(&fyne.Notification{
			Title:   fmt.Sprintf("Error adding word '%s'", word),
			Content: errText,
		})
	})
}

func (gui *GUI) showSuccessNotification(content string) {
	gui.logger.Infow("showing success notification", "content", content)
	fyne.Do(func() {
		gui.app.SendNotification(&fyne.Notification{
			Title:   gui.text.TextSuccess(),
			Content: content,
		})
	})
}

func (gui *GUI) handleGenerateSentences(params *generateSentencesParams) {
	gui.logger.Infow("handle generate sentences", "word", params.word, "wordLang", params.wordLang, "translationLang", params.translationLang, "deck", params.deck, "includeAudio", params.includeAudio)
	wordLang, ok := gui.text.GetLanguageCode(params.wordLang)
	if !ok {
		gui.logger.Errorw("unknown word language", "lang", params.wordLang)
		gui.showErrorNotification(ErrUnknownLanguage, params.word)
		return
	}
	translationLang, ok := gui.text.GetLanguageCode(params.translationLang)
	if !ok {
		gui.logger.Errorw("unknown translation language", "lang", params.translationLang)
		gui.showErrorNotification(ErrUnknownLanguage, params.word)
		return
	}

	_, err := gui.core.GenerateSentence(context.Background(), &core.GenerateSentenceRequest{
		Word:                params.word,
		WordLanguage:        wordLang,
		TranslationLanguage: translationLang,
		TranslationHint:     params.translationHint,
		IncludeAudio:        params.includeAudio,
		AudioGender:         params.audioGender,
		DeckName:            params.deck,
	})
	if err != nil {
		gui.logger.Errorw("generate sentence failed", "word", params.word, "err", err)
		gui.showErrorNotification(err, params.word)
		return
	}
	gui.logger.Infow("generate sentence succeeded", "word", params.word)
	gui.showSuccessNotification(fmt.Sprintf(gui.text.TextSentenceGeneratedSuccessfully(), params.word))
}

func (gui *GUI) handleTranslation(params *translateParams) {
	gui.logger.Infow("handle translation", "word", params.Word, "wordLang", params.WordLang, "translationLang", params.TranslationLang, "deck", params.Deck, "includeAudio", params.IncludeAudio)
	wordLang, ok := gui.text.GetLanguageCode(params.WordLang)
	if !ok {
		gui.logger.Errorw("unknown word language", "lang", params.WordLang)
		gui.showErrorNotification(ErrUnknownLanguage, params.Word)
		return
	}
	translationLang, ok := gui.text.GetLanguageCode(params.TranslationLang)
	if !ok {
		gui.logger.Errorw("unknown translation language", "lang", params.TranslationLang)
		gui.showErrorNotification(ErrUnknownLanguage, params.Word)
		return
	}

	_, err := gui.core.Translate(context.Background(), &core.TranslateRequest{
		Word:            params.Word,
		WordLanguage:    wordLang,
		TranslationLang: translationLang,
		TranslationHint: params.TranslationHint,
		IncludeAudio:    params.IncludeAudio,
		AudioGender:     params.AudioGender,
		DeckName:        params.Deck,
	})
	if err != nil {
		gui.logger.Errorw("translation failed", "word", params.Word, "err", err)
		gui.showErrorNotification(err, params.Word)
		return
	}
	gui.logger.Infow("translation succeeded", "word", params.Word)
	gui.showSuccessNotification(fmt.Sprintf(gui.text.TextTranslationAddedSuccessfully(), params.Word))
}

func (gui *GUI) handleGenerateDefinition(params *generateDefinitionParams) {
	gui.logger.Infow("handle generate definition", "word", params.Word, "wordLang", params.WordLang, "deck", params.Deck, "includeAudio", params.IncludeAudio)
	wordLang, ok := gui.text.GetLanguageCode(params.WordLang)
	if !ok {
		gui.logger.Errorw("unknown word language", "lang", params.WordLang)
		gui.showErrorNotification(ErrUnknownLanguage, params.Word)
		return
	}

	_, err := gui.core.GenerateDefinition(context.Background(), &core.GenerateDefinitionRequest{
		Word:           params.Word,
		Language:       wordLang,
		DefinitionHint: params.DefinitionHint,
		IncludeAudio:   params.IncludeAudio,
		AudioGender:    params.AudioGender,
		DeckName:       params.Deck,
	})
	if err != nil {
		gui.logger.Errorw("generate definition failed", "word", params.Word, "err", err)
		gui.showErrorNotification(err, params.Word)
		return
	}
	gui.logger.Infow("generate definition succeeded", "word", params.Word)
	gui.showSuccessNotification(fmt.Sprintf(gui.text.TextDefinitionAddedSuccessfully(), params.Word))
}
