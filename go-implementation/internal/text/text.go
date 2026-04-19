package text

type Text interface {
	GetLanguages() []string
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
	TextDeck() string
	GetLanguageCode(string) (string, bool)
	TextErrUnknownLanguage() string
	TextErrAnkiError() string
	TextErrInvalidArgument() string
	TextErrDeadlineExceeded() string
	TextErrInternalServer() string
	TextErrResourceExhausted() string
	TextErrUnavailable() string
	TextErrUnknown() string
	TextErrAddingWord() string
	TextSuccess() string
	TextSentenceGeneratedSuccessfully() string
	TextTranslationAddedSuccessfully() string
	TextDefinitionAddedSuccessfully() string
	TextTutorialTitle() string
	TextTutorialDescription() string
	TextTutorialLink() string
	TextErrPickLanguages() string
	TextErrLanguagesSame() string
	TextErrPickLanguage() string
	TextErrPickDeck() string
}

func NewText() Text {
	return &En{}
}
