package appdata

type State struct {
	GenerateSentence   GenerateSentenceState
	Translate          TranslateState
	GenerateDefinition GenerateDefinitionState
}

type GenerateSentenceState struct {
	Word            string
	TranslationHint string
	WordLang        string
	TranslationLang string
	IncludeAudio    bool
	AudioGender     string
}

type TranslateState struct {
	Word            string
	TranslationHint string
	WordLang        string
	TranslationLang string
	IncludeAudio    bool
	AudioGender     string
}

type GenerateDefinitionState struct {
	Word           string
	DefinitionHint string
	WordLang       string
	IncludeAudio   bool
	AudioGender    string
}

// Settings is a placeholder to be extended later.
type Settings struct{}
