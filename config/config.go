package config

type Language struct {
	DisplayName string
	Code        string
}

var Languages = []Language{
	{DisplayName: "English", Code: "en-US"},
	{DisplayName: "Russian", Code: "ru-RU"},
	{DisplayName: "Turkish", Code: "tr-TR"},
	{DisplayName: "French", Code: "fr-FR"},
	{DisplayName: "Italian", Code: "it-IT"},
}

func GetLanguageCode(displayName string) (string, bool) {
	for _, l := range Languages {
		if l.DisplayName == displayName {
			return l.Code, true
		}
	}
	return "", false
}

func GetDisplayNames() []string {
	names := make([]string, len(Languages))
	for i, l := range Languages {
		names[i] = l.DisplayName
	}
	return names
}
