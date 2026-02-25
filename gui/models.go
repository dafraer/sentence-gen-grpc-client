package gui

import "fyne.io/fyne/v2/widget"

type GenerateSentencesParams struct {
	Word            string
	TranslationHint string
	WordLang        string
	TranslationLang string
	IncludeAudio    bool
	AudioGender     string
}

type generateSentenceFormParams struct {
	word            *widget.Entry
	translationHint *widget.Entry
	wordLang        *widget.Select
	translationLang *widget.Select
	voice           *widget.Select
	audio           *widget.Check
}

type onGenerateSentenceSubmitParams struct {
	form            *widget.Form
	word            *widget.Entry
	translationHint *widget.Entry
	wordLang        *widget.Select
	translationLang *widget.Select
	voice           *widget.Select
	audio           *widget.Check
}
