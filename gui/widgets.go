package gui

import (
	"fyne.io/fyne/v2/widget"
)

func (gui *GUI) createGenerateSentenceForm(params *generateSentenceFormParams) *widget.Form {
	form := widget.NewForm(
		&widget.FormItem{Text: gui.text.TextWord(), Widget: params.word},
		&widget.FormItem{Text: gui.text.TextTranslationHint(), Widget: params.translationHint},
		&widget.FormItem{Text: gui.text.TextPickWordLanguage(), Widget: params.wordLang},
		&widget.FormItem{Text: gui.text.TextPickTranslationLanguage(), Widget: params.translationLang},
		&widget.FormItem{Text: gui.text.TextAudio(), Widget: params.audio},
		&widget.FormItem{Text: gui.text.TextVoice(), Widget: params.voice},
	)
	form.OnSubmit = func() {
		gui.onGenerateSentenceSubmit(&onGenerateSentenceSubmitParams{
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

func (gui *GUI) onGenerateSentenceSubmit(params *onGenerateSentenceSubmitParams) {
	go gui.handleGenerateSentences(&GenerateSentencesParams{
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
