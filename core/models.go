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

type GenerateSentenceResponse struct {
	OriginalSentence   string
	TranslatedSentence string
	Audio              []byte
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

type TranslateResponse struct {
	Translation string
	Audio       []byte
}

type GenerateDefinitionRequest struct {
	Word           string
	Language       string
	DefinitionHint string
	IncludeAudio   bool
	AudioGender    string
	DeckName       string
}

type GenerateDefinitionResponse struct {
	Definition string
	Audio      []byte
}
