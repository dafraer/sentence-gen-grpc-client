package text

type En struct{}
type Language struct {
	DisplayName string
	Code        string
}

// languages has all languages that have Chirp3-HD voice
var languages = []Language{
	{DisplayName: "Arabic", Code: "ar-XA"},
	{DisplayName: "Bengali (India)", Code: "bn-IN"},
	{DisplayName: "Bulgarian", Code: "bg-BG"},
	{DisplayName: "Chinese (Hong Kong)", Code: "yue-HK"},
	{DisplayName: "Chinese (Mandarin)", Code: "cmn-CN"},
	{DisplayName: "Croatian", Code: "hr-HR"},
	{DisplayName: "Czech", Code: "cs-CZ"},
	{DisplayName: "Danish", Code: "da-DK"},
	{DisplayName: "Dutch (Belgium)", Code: "nl-BE"},
	{DisplayName: "Dutch (Netherlands)", Code: "nl-NL"},
	{DisplayName: "English (US)", Code: "en-US"},
	{DisplayName: "English (Australia)", Code: "en-AU"},
	{DisplayName: "English (India)", Code: "en-IN"},
	{DisplayName: "English (UK)", Code: "en-GB"},
	{DisplayName: "Estonian", Code: "et-EE"},
	{DisplayName: "Finnish", Code: "fi-FI"},
	{DisplayName: "French", Code: "fr-FR"},
	{DisplayName: "French (Canada)", Code: "fr-CA"},
	{DisplayName: "German", Code: "de-DE"},
	{DisplayName: "Greek", Code: "el-GR"},
	{DisplayName: "Gujarati (India)", Code: "gu-IN"},
	{DisplayName: "Hebrew", Code: "he-IL"},
	{DisplayName: "Hindi (India)", Code: "hi-IN"},
	{DisplayName: "Hungarian", Code: "hu-HU"},
	{DisplayName: "Indonesian", Code: "id-ID"},
	{DisplayName: "Italian", Code: "it-IT"},
	{DisplayName: "Japanese", Code: "ja-JP"},
	{DisplayName: "Kannada (India)", Code: "kn-IN"},
	{DisplayName: "Korean", Code: "ko-KR"},
	{DisplayName: "Latvian", Code: "lv-LV"},
	{DisplayName: "Lithuanian", Code: "lt-LT"},
	{DisplayName: "Malayalam (India)", Code: "ml-IN"},
	{DisplayName: "Marathi (India)", Code: "mr-IN"},
	{DisplayName: "Norwegian Bokmål", Code: "nb-NO"},
	{DisplayName: "Polish", Code: "pl-PL"},
	{DisplayName: "Portuguese (Brazil)", Code: "pt-BR"},
	{DisplayName: "Punjabi (India)", Code: "pa-IN"},
	{DisplayName: "Romanian", Code: "ro-RO"},
	{DisplayName: "Russian", Code: "ru-RU"},
	{DisplayName: "Serbian", Code: "sr-RS"},
	{DisplayName: "Slovak", Code: "sk-SK"},
	{DisplayName: "Slovenian", Code: "sl-SI"},
	{DisplayName: "Spanish", Code: "es-ES"},
	{DisplayName: "Spanish (US)", Code: "es-US"},
	{DisplayName: "Swahili", Code: "sw-KE"},
	{DisplayName: "Swedish", Code: "sv-SE"},
	{DisplayName: "Tamil (India)", Code: "ta-IN"},
	{DisplayName: "Telugu (India)", Code: "te-IN"},
	{DisplayName: "Thai", Code: "th-TH"},
	{DisplayName: "Turkish", Code: "tr-TR"},
	{DisplayName: "Ukrainian", Code: "uk-UA"},
	{DisplayName: "Urdu (India)", Code: "ur-IN"},
	{DisplayName: "Vietnamese", Code: "vi-VN"},
}

func (t *En) GetLanguages() []string {
	names := make([]string, len(languages))
	for i, l := range languages {
		names[i] = l.DisplayName
	}
	return names
}

func (t *En) GetGenders() []string {
	return []string{"Male", "Female"}
}

func (t *En) GetPageNames() []string {
	return []string{"Generate sentence", "Translate", "Generate definition", "Tutorial"}
}

func (t *En) TextMale() string {
	return "Male"
}

func (t *En) TextFemale() string {
	return "Female"
}

func (t *En) TextWord() string {
	return "Word"
}

func (t *En) TextTranslationHint() string {
	return "Translation Hint"
}

func (t *En) TextPickWordLanguage() string {
	return "Pick word language"
}

func (t *En) TextPickTranslationLanguage() string {
	return "Pick translation language"
}

func (t *En) TextAudio() string {
	return "Audio"
}

func (t *En) TextVoice() string {
	return "Voice"
}

func (t *En) TextGenerateSentenceTitle() string {
	return "# Generate sentences"
}

func (t *En) TextErrWordRequired() string {
	return "word required"
}

func (t *En) TextGenerate() string {
	return "Generate"
}

func (t *En) TextErrHintTooLong() string {
	return "hint too long"
}

func (t *En) TextErrWordTooLong() string {
	return "word too long"
}

func (t *En) TextGenerateTranslationTitle() string {
	return "# Generate translation"
}

func (t *En) TextGenerateDefinitionTitle() string {
	return "# Generate definition"
}
func (t *En) TextDefinitionHint() string {
	return "Definition hint"
}

func (t *En) TextSettingsTitle() string {
	return "# Settings"
}

func (t *En) TextDeck() string {
	return "Deck"
}

func (t *En) GetLanguageCode(displayName string) (string, bool) {
	for _, l := range languages {
		if l.DisplayName == displayName {
			return l.Code, true
		}
	}
	return "", false
}

func (t *En) TextError() string {
	return "Error"
}

func (t *En) TextErrUnknownLanguage() string {
	return "Unknown language"
}

func (t *En) TextErrAnkiError() string {
	return "Anki error"
}

func (t *En) TextErrInvalidArgument() string {
	return "Invalid argument"
}

func (t *En) TextErrDeadlineExceeded() string {
	return "Server deadline exceeded"
}

func (t *En) TextErrInternalServer() string {
	return "Internal server error"
}

func (t *En) TextErrResourceExhausted() string {
	return "Quota limit reached, try again tomorrow"
}

func (t *En) TextErrUnavailable() string {
	return "Server is unavailable"
}

func (t *En) TextErrUnknown() string {
	return "Unknown error"
}

func (t *En) TextErrAddingWord() string {
	return "Error adding word '%s'"
}

func (t *En) TextSuccess() string {
	return "Success"
}

func (t *En) TextSentenceGeneratedSuccessfully() string {
	return "Sentence with the word '%s' has been added successfully"
}

func (t *En) TextTranslationAddedSuccessfully() string {
	return "Translation with the word '%s' has been added successfully"
}

func (t *En) TextDefinitionAddedSuccessfully() string {
	return "Definition of the word '%s' has been added successfully"
}

func (t *En) TextTutorialTitle() string {
	return "# Tutorial"
}

func (t *En) TextTutorialDescription() string {
	return "You can check out the tutorial on the project's"
}

func (t *En) TextTutorialLink() string {
	return "github page"
}
