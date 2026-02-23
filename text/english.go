package text

type En struct{}

func (t *En) GetLanguageList() []string {
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

func (t *En) GetPageList() []string {
	return []string{"Home", "Generate sentence", "Translate", "Generate definition", "Settings"}
}

func (t *En) Male() string {
	return "Male"
}

func (t *En) Female() string {
	return "Female"
}

func (t *En) Word() string {
	return "Word"
}

func (t *En) TranslationHint() string {
	return "Translation Hint"
}

func (t *En) PickWordLang() string {
	return "Pick word language"
}

func (t *En) PickTranslationLang() string {
	return "Pick translation language"
}

func (t *En) IncludeAudio() string {
	return "Include audio"
}

func (t *En) VoiceGender() string {
	return "Voice Gender"
}
