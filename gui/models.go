package gui

import "fyne.io/fyne/v2/widget"

type generateSentencesParams struct {
	word            string
	translationHint string
	wordLang        string
	translationLang string
	includeAudio    bool
	audioGender     string
}

type formParams struct {
	word            *widget.Entry
	translationHint *widget.Entry
	wordLang        *widget.Select
	translationLang *widget.Select
	voice           *widget.Select
	audio           *widget.Check
	onSubmit        func(params *onSubmitParams)
}

type onSubmitParams struct {
	form            *widget.Form
	word            *widget.Entry
	translationHint *widget.Entry
	wordLang        *widget.Select
	translationLang *widget.Select
	voice           *widget.Select
	audio           *widget.Check
}

type translateParams struct {
	Word            string
	TranslationHint string
	WordLang        string
	TranslationLang string
	IncludeAudio    bool
	AudioGender     string
}

type generateDefinitionParams struct {
	Word           string
	DefinitionHint string
	WordLang       string
	IncludeAudio   bool
	AudioGender    string
}

type definitionFormParams struct {
	word           *widget.Entry
	definitionHint *widget.Entry
	wordLang       *widget.Select
	voice          *widget.Select
	audio          *widget.Check
}
