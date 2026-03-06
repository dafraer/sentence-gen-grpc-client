package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// minWidthLayout is a layout that enforces a minimum width on its container.
type minWidthLayout struct {
	width float32
}

func (l *minWidthLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Move(fyne.NewPos(0, 0))
		o.Resize(size)
	}
}

func (l *minWidthLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(l.width, 0)
}

// listMinWidth calculates the minimum width that a list can take
// Minimum width is based on list's longest element
func listMinWidth(items []string) float32 {
	longest := ""
	for _, item := range items {
		if len(item) > len(longest) {
			longest = item
		}
	}
	return widget.NewLabel(longest).MinSize().Width + theme.Padding()
}

// createForm creates form for Translate and Generate sentence pages
func (gui *GUI) createForm(params *formParams) *widget.Form {
	//Create new form
	form := widget.NewForm(
		&widget.FormItem{Text: gui.text.TextWord(), Widget: params.word},
		&widget.FormItem{Text: gui.text.TextTranslationHint(), Widget: params.translationHint},
		&widget.FormItem{Text: gui.text.TextPickWordLanguage(), Widget: params.wordLang},
		&widget.FormItem{Text: gui.text.TextPickTranslationLanguage(), Widget: params.translationLang},
		&widget.FormItem{Text: gui.text.TextAudio(), Widget: params.audio},
		&widget.FormItem{Text: gui.text.TextVoice(), Widget: params.voice},
		&widget.FormItem{Text: gui.text.TextDeck(), Widget: params.deck},
	)
	//Define onSubmit function from the params and disable the form
	form.OnSubmit = func() {
		params.onSubmit(&onSubmitParams{
			form:            form,
			word:            params.word,
			translationHint: params.translationHint,
			wordLang:        params.wordLang,
			translationLang: params.translationLang,
			audio:           params.audio,
			voice:           params.voice,
			deck:            params.deck,
		})
		form.Disable()
	}
	form.SubmitText = gui.text.TextGenerate()
	//Refresh the form so we can disable it
	form.Refresh()
	form.Disable()
	return form
}

// onGenerateSentenceSubmit is called when Generate Sentence form is submitted
func (gui *GUI) onGenerateSentenceSubmit(params *onSubmitParams) {
	//Extra check to make sure invalid data is not passed to the handler
	if err := gui.validateForm(params.word.Text,
		params.translationHint.Text,
		params.wordLang.Selected,
		params.translationLang.Selected,
		params.deck.Selected); err != nil {
		gui.logger.Errorw("generate sentence form validation failed", "word", params.word.Text, "err", err)

		//In case of validation error show error dialog and disable the form
		params.form.Disable()
		dialog.NewError(err, gui.window).Show()
		return
	}
	gui.logger.Debugw("generate sentence form submitted", "word", params.word.Text, "wordLang", params.wordLang.Selected, "translationLang", params.translationLang.Selected, "deck", params.deck.Selected)

	//Run handler in a separate goroutine
	go gui.handleGenerateSentences(&generateSentencesParams{
		word:            params.word.Text,
		translationHint: params.translationHint.Text,
		wordLang:        params.wordLang.Selected,
		translationLang: params.translationLang.Selected,
		includeAudio:    params.audio.Checked,
		audioGender:     params.voice.Selected,
		deck:            params.deck.Selected,
	})

	//Reset the fields that will be re-entered in the next call
	params.word.SetText("")
	params.translationHint.SetText("")
}

// onTranslateSubmit is called when Translate form is submitted
func (gui *GUI) onTranslateSubmit(params *onSubmitParams) {
	//Extra check to make sure invalid data is not passed to the handler
	if err := gui.validateForm(params.word.Text,
		params.translationHint.Text,
		params.wordLang.Selected,
		params.translationLang.Selected,
		params.deck.Selected); err != nil {
		gui.logger.Errorw("translate form validation failed", "word", params.word.Text, "err", err)
		params.form.Disable()
		dialog.NewError(err, gui.window).Show()
		return
	}
	gui.logger.Debugw("translate form submitted", "word", params.word.Text, "wordLang", params.wordLang.Selected, "translationLang", params.translationLang.Selected, "deck", params.deck.Selected)

	//Run handler in a separate goroutine
	go gui.handleTranslation(&translateParams{
		Word:            params.word.Text,
		TranslationHint: params.translationHint.Text,
		WordLang:        params.wordLang.Selected,
		TranslationLang: params.translationLang.Selected,
		IncludeAudio:    params.audio.Checked,
		AudioGender:     params.voice.Selected,
		Deck:            params.deck.Selected,
	})

	//Reset the fields that will be re-entered in the next call
	params.word.SetText("")
	params.translationHint.SetText("")
}

// createDefinitionForm creates form for Generate Definition page
func (gui *GUI) createDefinitionForm(params *definitionFormParams) *widget.Form {
	//Create new form
	form := widget.NewForm(
		&widget.FormItem{Text: gui.text.TextWord(), Widget: params.word},
		&widget.FormItem{Text: gui.text.TextDefinitionHint(), Widget: params.definitionHint},
		&widget.FormItem{Text: gui.text.TextPickWordLanguage(), Widget: params.wordLang},
		&widget.FormItem{Text: gui.text.TextAudio(), Widget: params.audio},
		&widget.FormItem{Text: gui.text.TextVoice(), Widget: params.voice},
		&widget.FormItem{Text: gui.text.TextDeck(), Widget: params.deck},
	)

	//Define onSubmit function
	form.OnSubmit = func() {
		//Extra check to make sure invalid data is not passed to the handler
		if err := gui.validateDefinitionForm(params.word.Text,
			params.definitionHint.Text,
			params.wordLang.Selected,
			params.deck.Selected); err != nil {
			gui.logger.Errorw("generate definition form validation failed", "word", params.word.Text, "err", err)
			form.Disable()
			dialog.NewError(err, gui.window).Show()
			return
		}
		gui.logger.Debugw("generate definition form submitted", "word", params.word.Text, "wordLang", params.wordLang.Selected, "deck", params.deck.Selected)

		//Run handler in a separate goroutine
		go gui.handleGenerateDefinition(&generateDefinitionParams{
			Word:           params.word.Text,
			DefinitionHint: params.definitionHint.Text,
			WordLang:       params.wordLang.Selected,
			IncludeAudio:   params.audio.Checked,
			AudioGender:    params.voice.Selected,
			Deck:           params.deck.Selected,
		})

		//Disable the form and reset the fields that will be re-entered in the next call
		form.Disable()
		params.word.SetText("")
		params.definitionHint.SetText("")
	}
	form.SubmitText = gui.text.TextGenerate()
	//Refresh form right away so we can disable it
	form.Refresh()
	form.Disable()
	return form
}

// createVoicePicker creates a voice picker
// The voices available are Male and Female
func (gui *GUI) createVoicePicker() *widget.Select {
	genders := gui.text.GetGenders()
	voiceGenderSelect := widget.NewSelect(genders, nil)
	voiceGenderSelect.SetSelected(gui.text.TextFemale())
	voiceGenderSelect.Disable()
	return voiceGenderSelect
}
