package core

type GenerateSentenceRequest struct {
	Word                string
	WordLanguage        string
	TranslationLanguage string
	TranslationHint     string
	IncludeAudio        bool
	AudioGender         string
	DeckName            string
}

type TranslateRequest struct {
	Word            string
	WordLanguage    string
	TranslationLang string
	TranslationHint string
	IncludeAudio    bool
	AudioGender     string
	DeckName        string
}

type GenerateDefinitionRequest struct {
	Word           string
	Language       string
	DefinitionHint string
	IncludeAudio   bool
	AudioGender    string
	DeckName       string
}
