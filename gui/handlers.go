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
	fyne.Do(func() {
		gui.app.SendNotification(&fyne.Notification{
			Title:   gui.text.TextSuccess(),
			Content: content,
		})
	})
}

func (gui *GUI) handleGenerateSentences(params *generateSentencesParams) {
	wordLang, ok := gui.text.GetLanguageCode(params.wordLang)
	if !ok {
		gui.showErrorNotification(ErrUnknownLanguage, params.word)
		return
	}
	translationLang, ok := gui.text.GetLanguageCode(params.translationLang)
	if !ok {
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
		gui.showErrorNotification(err, params.word)
		return
	}
	gui.showSuccessNotification(fmt.Sprintf(gui.text.TextSentenceGeneratedSuccessfully(), params.word))
}

func (gui *GUI) handleTranslation(params *translateParams) {
	wordLang, ok := gui.text.GetLanguageCode(params.WordLang)
	if !ok {
		gui.showErrorNotification(ErrUnknownLanguage, params.Word)
		return
	}
	translationLang, ok := gui.text.GetLanguageCode(params.TranslationLang)
	if !ok {
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
		gui.showErrorNotification(err, params.Word)
		return
	}
	gui.showSuccessNotification(fmt.Sprintf(gui.text.TextTranslationAddedSuccessfully(), params.Word))
}

func (gui *GUI) handleGenerateDefinition(params *generateDefinitionParams) {
	wordLang, ok := gui.text.GetLanguageCode(params.WordLang)
	if !ok {
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
		gui.showErrorNotification(err, params.Word)
		return
	}
	gui.showSuccessNotification(fmt.Sprintf(gui.text.TextDefinitionAddedSuccessfully(), params.Word))
}
