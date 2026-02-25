package gui

import (
	"fmt"
	"strings"
)

const (
	maxHintLen = 100
	maxWordLen = 100
)

func (gui *GUI) validateWord(w string) error {
	if strings.TrimSpace(w) == "" {
		return fmt.Errorf(gui.text.TextErrWordRequired())
	}
	if len([]rune(w)) > maxWordLen {
		return fmt.Errorf(gui.text.TextErrWordTooLong())
	}
	return nil
}

func (gui *GUI) validateLanguages(lang1, lang2 string) error {
	fmt.Println(lang1, lang2)
	if lang1 == "" || lang2 == "" {
		return fmt.Errorf("pick languages")
	}
	if lang1 == lang2 {
		return fmt.Errorf("languages cannot be the same")
	}
	return nil
}

func (gui *GUI) validateHint(h string) error {
	if len([]rune(h)) > maxHintLen {
		return fmt.Errorf(gui.text.TextErrHintTooLong())
	}
	return nil
}

func (gui *GUI) validateGenerateSentenceForm(word, hint, lang1, lang2 string) error {
	if err := gui.validateWord(word); err != nil {
		return err
	}
	if err := gui.validateLanguages(lang1, lang2); err != nil {
		return err
	}
	return gui.validateHint(hint)
}
