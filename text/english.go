package text

type En struct{}

func (t *En) GetLanguages() []string {
	return []string{
		"English",
		"Russian",
		"Turkish",
		"French",
		"Italian",
	}
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
	return []string{"Home", "Generate sentence", "Translate", "Generate definition", "Settings"}
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

func (t *En) TextGenerateSentence() string {
	return "# Generate sentences"
}

func (t *En) TextErrWordRequired() string {
	return "Word required"
}

func (t *En) TextGenerate() string {
	return "Generate"
}

func (t *En) TextErrHintTooLong() string {
	return "Hint too long"
}

func (t *En) TextErrWordTooLong() string {
	return "Word too long"
}
