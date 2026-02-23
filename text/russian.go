package text

type Ru struct{}

func (t *Ru) GetLanguageList() []string {
	return []string{
		"-",
		"Английский",
		"Русский",
		"Турецкий",
		"Французский",
		"Итальянский",
	}
}

func (t *Ru) GetGenders() []string {
	return []string{"Мужской", "Женский"}
}

func (t *Ru) GetHomePage() string {
	return `
# Главная

## Это приложение:
- Генерирует предложения
- Переводит
- Генерирует определения
- Добавляет всё вышеперечисленное в вашу карточку Anki
`
}

func (t *Ru) GetPageList() []string {
	return []string{"Главная", "Сгенерировать предложение", "Перевести", "Сгенерировать определение", "Настройки"}
}

func (t *Ru) Male() string {
	return "Мужской"
}

func (t *Ru) Female() string {
	return "Женский"
}

func (t *Ru) Word() string {
	return "Слово"
}

func (t *Ru) TranslationHint() string {
	return "Подсказка для перевода"
}

func (t *Ru) PickWordLang() string {
	return "Выберите язык слова"
}

func (t *Ru) PickTranslationLang() string {
	return "Выберите язык перевода"
}

func (t *Ru) IncludeAudio() string {
	return "Аудио"
}

func (t *Ru) VoiceGender() string {
	return "Голос"
}

func (t *Ru) GenerateSentenceTitle() string {
	return "# Сгенерировать предложение"
}

func (t *Ru) ErrWordRequired() string {
	return "Введите слово"
}

func (t *Ru) GenerateButton() string {
	return "Сгенерировать"
}
