package gui

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/dafraer/sentence-gen-grpc-client/config"
	"github.com/dafraer/sentence-gen-grpc-client/core"
)

func (gui *GUI) showError(err error) {
	fyne.Do(func() {
		gui.app.SendNotification(&fyne.Notification{
			Title:   "Error",
			Content: err.Error(),
		})
		dialog.ShowError(err, gui.window)
	})
}

func (gui *GUI) handleGenerateSentences(params *generateSentencesParams) {
	wordLang, ok := config.GetLanguageCode(params.wordLang)
	if !ok {
		gui.showError(fmt.Errorf("unknown language: %s", params.wordLang))
		return
	}
	translationLang, ok := config.GetLanguageCode(params.translationLang)
	if !ok {
		gui.showError(fmt.Errorf("unknown language: %s", params.translationLang))
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
		gui.showError(err)
	}
}

func (gui *GUI) handleTranslation(params *translateParams) {
	wordLang, ok := config.GetLanguageCode(params.WordLang)
	if !ok {
		gui.showError(fmt.Errorf("unknown language: %s", params.WordLang))
		return
	}
	translationLang, ok := config.GetLanguageCode(params.TranslationLang)
	if !ok {
		gui.showError(fmt.Errorf("unknown language: %s", params.TranslationLang))
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
		gui.showError(err)
	}
}

func (gui *GUI) handleGenerateDefinition(params *generateDefinitionParams) {
	wordLang, ok := config.GetLanguageCode(params.WordLang)
	if !ok {
		gui.showError(fmt.Errorf("unknown language: %s", params.WordLang))
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
		gui.showError(err)
	}
}
