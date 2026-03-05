package text

type En struct{}
type Language struct {
	DisplayName string
	Code        string
}

var languages = []Language{
	{DisplayName: "English", Code: "en-US"},
	{DisplayName: "Russian", Code: "ru-RU"},
	{DisplayName: "Turkish", Code: "tr-TR"},
	{DisplayName: "French", Code: "fr-FR"},
	{DisplayName: "Italian", Code: "it-IT"},
	{DisplayName: "Spanish", Code: "es-ES"},
	{DisplayName: "German", Code: "de-DE"},
	{DisplayName: "Ukrainian", Code: "uk-UA"},
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

func (t *En) GetHomePage() string {
	return `
# Home 

## This app:
- Generates sentences 
- Translates
- Generates definitions
- Adds all of the above to your anki card
`
}

func (t *En) GetPageNames() []string {
	return []string{"Home", "Generate sentence", "Translate", "Generate definition"}
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
	return "Hint too long"
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
