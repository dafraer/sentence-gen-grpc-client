package text

type Text interface {
	GetLanguages() []string
	GetHomePage() string
	GetGenders() []string
	GetPageNames() []string
	TextMale() string
	TextFemale() string
	TextWord() string
	TextTranslationHint() string
	TextPickWordLanguage() string
	TextPickTranslationLanguage() string
	TextAudio() string
	TextVoice() string
	TextGenerateSentenceTitle() string
	TextErrWordRequired() string
	TextGenerate() string
	TextErrHintTooLong() string
	TextErrWordTooLong() string
	TextGenerateTranslationTitle() string
	TextGenerateDefinitionTitle() string
	TextDefinitionHint() string
	TextSettingsTitle() string
}

func NewText(language string) Text {
	//TODO: add other languages
	return &En{}
}
