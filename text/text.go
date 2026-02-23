package text

type Text interface {
	GetLanguageList() []string
	GetHomePage() string
	GetGenders() []string
	GetPageList() []string
}

func NewText(language string) Text {
	if language == "ru" {
		return &Ru{}
	}
	return &En{}
}
