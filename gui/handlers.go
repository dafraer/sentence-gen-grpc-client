package gui

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2/dialog"
	"github.com/dafraer/sentence-gen-grpc-client/config"
	"github.com/dafraer/sentence-gen-grpc-client/core"
)

func (gui *GUI) handleGenerateSentences(params *generateSentencesParams) {
	wordLang, ok := config.GetLanguageCode(params.wordLang)
	if !ok {
		dialog.ShowError(fmt.Errorf("unknown language: %s", params.wordLang), gui.window)
		return
	}
	translationLang, ok := config.GetLanguageCode(params.translationLang)
	if !ok {
		dialog.ShowError(fmt.Errorf("unknown language: %s", params.translationLang), gui.window)
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
		dialog.ShowError(err, gui.window)
	}
}

func (gui *GUI) handleTranslation(params *translateParams) {
	wordLang, ok := config.GetLanguageCode(params.WordLang)
	if !ok {
		dialog.ShowError(fmt.Errorf("unknown language: %s", params.WordLang), gui.window)
		return
	}
	translationLang, ok := config.GetLanguageCode(params.TranslationLang)
	if !ok {
		dialog.ShowError(fmt.Errorf("unknown language: %s", params.TranslationLang), gui.window)
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
		dialog.ShowError(err, gui.window)
	}
}

func (gui *GUI) handleGenerateDefinition(params *generateDefinitionParams) {
	wordLang, ok := config.GetLanguageCode(params.WordLang)
	if !ok {
		dialog.ShowError(fmt.Errorf("unknown language: %s", params.WordLang), gui.window)
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
		dialog.ShowError(err, gui.window)
	}
}
