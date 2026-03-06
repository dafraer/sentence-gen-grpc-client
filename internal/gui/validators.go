package gui

import (
	"errors"
	"strings"
)

const (
	maxHintLen = 100
	maxWordLen = 100
)

// validateWord checks that a word is not empty and not too long
func (gui *GUI) validateWord(w string) error {
	if strings.TrimSpace(w) == "" {
		return errors.New(gui.text.TextErrWordRequired())
	}
	if len([]rune(w)) > maxWordLen {
		return errors.New(gui.text.TextErrWordTooLong())
	}
	return nil
}

// validateLanguages checks that languages are not empty and not identical
func (gui *GUI) validateLanguages(lang1, lang2 string) error {
	if lang1 == "" || lang2 == "" {
		return errors.New(gui.text.TextErrPickLanguages())
	}
	if lang1 == lang2 {
		return errors.New(gui.text.TextErrLanguagesSame())
	}
	return nil
}

// validateHint checks that the hint is not too long
func (gui *GUI) validateHint(h string) error {
	if len([]rune(h)) > maxHintLen {
		return errors.New(gui.text.TextErrHintTooLong())
	}
	return nil
}

// validateForm validates each field in the form
func (gui *GUI) validateForm(word, hint, lang1, lang2, deck string) error {
	if err := gui.validateWord(word); err != nil {
		return err
	}
	if err := gui.validateLanguages(lang1, lang2); err != nil {
		return err
	}
	if deck == "" {
		return errors.New(gui.text.TextErrPickDeck())
	}
	return gui.validateHint(hint)
}

// validateDefinitionForm checks each field in the definition form
func (gui *GUI) validateDefinitionForm(word, hint, lang, deck string) error {
	if err := gui.validateWord(word); err != nil {
		return err
	}
	if lang == "" {
		return errors.New(gui.text.TextErrPickLanguage())
	}
	if deck == "" {
		return errors.New(gui.text.TextErrPickDeck())
	}
	return gui.validateHint(hint)
}
