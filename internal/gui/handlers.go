package gui

import (
	"context"
	"errors"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"github.com/dafraer/sentence-gen-grpc-client/internal/anki"
	"github.com/dafraer/sentence-gen-grpc-client/internal/core"
	"github.com/dafraer/sentence-gen-grpc-client/internal/rpc"
)

const timeOut = time.Second * 50

var ErrUnknownLanguage = errors.New("unknown language")

// showErrorNotification shows system notification for an error and a word for which error occurred
func (gui *GUI) showErrorNotification(err error, word string) {
	gui.logger.Errorw("showing error notification", "word", word, "err", err)
	errText := ""
	switch {
	case errors.Is(err, ErrUnknownLanguage):
		errText = gui.text.TextErrUnknownLanguage()
	case errors.Is(err, anki.ErrAnkiError):
		errText = gui.text.TextErrAnkiError()
	case errors.Is(err, rpc.ErrInvalidArgument):
		errText = gui.text.TextErrInvalidArgument()
	case errors.Is(err, rpc.ErrDeadlineExceeded):
		errText = gui.text.TextErrDeadlineExceeded()
	case errors.Is(err, rpc.ErrInternalServer):
		errText = gui.text.TextErrInternalServer()
	case errors.Is(err, rpc.ErrResourceExhausted):
		errText = gui.text.TextErrResourceExhausted()
	case errors.Is(err, rpc.ErrUnavailable):
		errText = gui.text.TextErrUnavailable()
	default:
		errText = gui.text.TextErrUnknown()
	}

	//Send notification using fyne.Do because we are in another goroutine
	fyne.Do(func() {
		gui.app.SendNotification(&fyne.Notification{
			Title:   fmt.Sprintf(gui.text.TextErrAddingWord(), word),
			Content: errText,
		})
	})
}

// showSuccessNotification shows notification when anki card is successfully added
func (gui *GUI) showSuccessNotification(content string) {
	gui.logger.Infow("showing success notification", "content", content)

	//Send notification using fyne.Do because we are in another goroutine
	fyne.Do(func() {
		gui.app.SendNotification(&fyne.Notification{
			Title:   gui.text.TextSuccess(),
			Content: content,
		})
	})
}

// handleGenerateSentences calls core method to generate sentences and add them to anki
func (gui *GUI) handleGenerateSentences(params *generateSentencesParams) {
	gui.logger.Infow("handle generate sentences", "word", params.word, "wordLang", params.wordLang, "translationLang", params.translationLang, "deck", params.deck, "includeAudio", params.includeAudio)

	//Get language code since server expects it
	wordLang, ok := gui.text.GetLanguageCode(params.wordLang)
	if !ok {
		gui.logger.Errorw("unknown word language", "lang", params.wordLang)
		gui.showErrorNotification(ErrUnknownLanguage, params.word)
		return
	}

	//Get language code for translation language
	translationLang, ok := gui.text.GetLanguageCode(params.translationLang)
	if !ok {
		gui.logger.Errorw("unknown translation language", "lang", params.translationLang)
		gui.showErrorNotification(ErrUnknownLanguage, params.word)
		return
	}

	//Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	//Generate sentence and add it to anki
	if err := gui.core.GenerateSentence(ctx, &core.GenerateSentenceRequest{
		Word:                params.word,
		WordLanguage:        wordLang,
		TranslationLanguage: translationLang,
		TranslationHint:     params.translationHint,
		IncludeAudio:        params.includeAudio,
		AudioGender:         params.audioGender,
		DeckName:            params.deck,
	}); err != nil {
		gui.logger.Errorw("generate sentence failed", "word", params.word, "err", err)
		gui.showErrorNotification(err, params.word)
		return
	}
	gui.logger.Infow("generate sentence succeeded", "word", params.word)
	gui.showSuccessNotification(fmt.Sprintf(gui.text.TextSentenceGeneratedSuccessfully(), params.word))
}

// handleTranslation calls core method to generate translation and add it to anki
func (gui *GUI) handleTranslation(params *translateParams) {
	gui.logger.Infow("handle translation", "word", params.word, "wordLang", params.wordLang, "translationLang", params.translationLang, "deck", params.deck, "includeAudio", params.includeAudio)

	//Get language code since server expects language code
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

	//Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	//Call translate method to generate the translation and add it to anki
	if err := gui.core.Translate(ctx, &core.TranslateRequest{
		Word:            params.word,
		WordLanguage:    wordLang,
		TranslationLang: translationLang,
		TranslationHint: params.translationHint,
		IncludeAudio:    params.includeAudio,
		AudioGender:     params.audioGender,
		DeckName:        params.deck,
	}); err != nil {
		gui.logger.Errorw("translation failed", "word", params.word, "err", err)
		gui.showErrorNotification(err, params.word)
		return
	}

	gui.logger.Infow("translation succeeded", "word", params.word)
	gui.showSuccessNotification(fmt.Sprintf(gui.text.TextTranslationAddedSuccessfully(), params.word))
}

// handleGenerateDefinition calls core method to generate definition and add it to anki
func (gui *GUI) handleGenerateDefinition(params *generateDefinitionParams) {
	gui.logger.Infow("handle generate definition", "word", params.word, "wordLang", params.wordLang, "deck", params.deck, "includeAudio", params.includeAudio)

	//Get language code since server expects it
	wordLang, ok := gui.text.GetLanguageCode(params.wordLang)
	if !ok {
		gui.logger.Errorw("unknown word language", "lang", params.wordLang)
		gui.showErrorNotification(ErrUnknownLanguage, params.word)
		return
	}

	//Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	//Call core method to generate definition and add it to anki
	if err := gui.core.GenerateDefinition(ctx, &core.GenerateDefinitionRequest{
		Word:           params.word,
		Language:       wordLang,
		DefinitionHint: params.definitionHint,
		IncludeAudio:   params.includeAudio,
		AudioGender:    params.audioGender,
		DeckName:       params.deck,
	}); err != nil {
		gui.logger.Errorw("generate definition failed", "word", params.word, "err", err)
		gui.showErrorNotification(err, params.word)
		return
	}

	gui.logger.Infow("generate definition succeeded", "word", params.word)
	gui.showSuccessNotification(fmt.Sprintf(gui.text.TextDefinitionAddedSuccessfully(), params.word))
}
