package text

type Text interface {
	GetLanguageList() []string
	GetHomePage() string
	GetGenders() []string
	GetPageList() []string
	Male() string
	Female() string
	Word() string
	TranslationHint() string
	PickWordLang() string
	PickTranslationLang() string
	IncludeAudio() string
	VoiceGender() string
	GenerateSentenceTitle() string
	ErrWordRequired() string
	GenerateButton() string
}

func NewText(language string) Text {
	if language == "ru" {
		return &Ru{}
	}
	return &En{}
}
