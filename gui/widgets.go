package gui

import (
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (gui *GUI) createForm(params *formParams) *widget.Form {
	form := widget.NewForm(
		&widget.FormItem{Text: gui.text.TextWord(), Widget: params.word},
		&widget.FormItem{Text: gui.text.TextTranslationHint(), Widget: params.translationHint},
		&widget.FormItem{Text: gui.text.TextPickWordLanguage(), Widget: params.wordLang},
		&widget.FormItem{Text: gui.text.TextPickTranslationLanguage(), Widget: params.translationLang},
		&widget.FormItem{Text: gui.text.TextAudio(), Widget: params.audio},
		&widget.FormItem{Text: gui.text.TextVoice(), Widget: params.voice},
	)
	form.OnSubmit = func() {
		params.onSubmit(&onSubmitParams{
			form:            form,
			word:            params.word,
			translationHint: params.translationHint,
			wordLang:        params.wordLang,
			translationLang: params.translationLang,
			audio:           params.audio,
			voice:           params.voice,
		})
		form.Disable()
	}
	form.SubmitText = gui.text.TextGenerate()
	form.Refresh()
	form.Disable()
	return form
}

func (gui *GUI) onGenerateSentenceSubmit(params *onSubmitParams) {
	//Extra check to make sure invalid data is not passed to the handler
	if err := gui.validateForm(params.word.Text,
		params.translationHint.Text,
		params.wordLang.Selected,
		params.translationLang.Selected); err != nil {
		params.form.Disable()
		dialog.NewError(err, gui.window).Show()
		return
	}
	go gui.handleGenerateSentences(&generateSentencesParams{
		word:            params.word.Text,
		translationHint: params.translationHint.Text,
		wordLang:        params.wordLang.Selected,
		translationLang: params.translationLang.Selected,
		includeAudio:    params.audio.Checked,
		audioGender:     params.voice.Selected,
	})
	params.word.SetText("")
	params.translationHint.SetText("")
}

func (gui *GUI) onTranslateSubmit(params *onSubmitParams) {
	//Extra check to make sure invalid data is not passed to the handler
	if err := gui.validateForm(params.word.Text,
		params.translationHint.Text,
		params.wordLang.Selected,
		params.translationLang.Selected); err != nil {
		params.form.Disable()
		dialog.NewError(err, gui.window).Show()
		return
	}
	go gui.handleTranslation(&translateParams{
		Word:            params.word.Text,
		TranslationHint: params.translationHint.Text,
		WordLang:        params.wordLang.Selected,
		TranslationLang: params.translationLang.Selected,
		IncludeAudio:    params.audio.Checked,
		AudioGender:     params.voice.Selected,
	})
	params.word.SetText("")
	params.translationHint.SetText("")
}

func (gui *GUI) createDefinitionForm(params *definitionFormParams) *widget.Form {
	form := widget.NewForm(
		&widget.FormItem{Text: gui.text.TextWord(), Widget: params.word},
		&widget.FormItem{Text: gui.text.TextDefinitionHint(), Widget: params.definitionHint},
		&widget.FormItem{Text: gui.text.TextPickWordLanguage(), Widget: params.wordLang},
		&widget.FormItem{Text: gui.text.TextAudio(), Widget: params.audio},
		&widget.FormItem{Text: gui.text.TextVoice(), Widget: params.voice},
	)
	form.OnSubmit = func() {
		//Extra check to make sure invalid data is not passed to the handler
		if err := gui.validateDefinitionForm(params.word.Text,
			params.definitionHint.Text,
			params.wordLang.Selected); err != nil {
			form.Disable()
			dialog.NewError(err, gui.window).Show()
			return
		}
		go gui.handleGenerateDefinition(&generateDefinitionParams{
			Word:           params.word.Text,
			DefinitionHint: params.definitionHint.Text,
			WordLang:       params.wordLang.Selected,
			IncludeAudio:   params.audio.Checked,
			AudioGender:    params.voice.Selected,
		})
		form.Disable()
		params.word.SetText("")
		params.definitionHint.SetText("")
	}
	form.SubmitText = gui.text.TextGenerate()
	form.Refresh()
	form.Disable()
	return form
}

func (gui *GUI) createVoicePicker() *widget.Select {
	genders := gui.text.GetGenders()
	voiceGenderSelect := widget.NewSelect(genders, nil)
	voiceGenderSelect.SetSelected(gui.text.TextFemale())
	voiceGenderSelect.Disable()
	return voiceGenderSelect
}
