package gui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dafraer/sentence-gen-grpc-client/core"
	"github.com/dafraer/sentence-gen-grpc-client/text"
	"go.uber.org/zap"
)

type GUI struct {
	logger *zap.SugaredLogger
	core   *core.Core
	text   text.Text
}

func New(logger *zap.SugaredLogger, core *core.Core, text text.Text) *GUI {
	return &GUI{
		logger: logger,
		core:   core,
		text:   text,
	}
}

func (gui *GUI) Run() {
	//Create new instance of the GUI app
	a := app.New()

	//Create new window
	w := a.NewWindow("Sengen")

	//Items in the list
	items := gui.text.GetPageList()

	//Create pages
	homePage := gui.CreateHomePage()
	generateSentencePage := gui.CreateGenerateSentencePage()
	translatePage := gui.CreateTranslatePage()
	generateDefinitionPage := gui.CreateGenerateDefinitionPage()
	settingsPage := gui.CreateSettingsPage()

	//Create variable for contents of the page
	content := container.NewStack(homePage) // right side

	//Create map of all pages
	pages := map[string]fyne.CanvasObject{
		items[0]: homePage,
		items[1]: generateSentencePage,
		items[2]: translatePage,
		items[3]: generateDefinitionPage,
		items[4]: settingsPage,
	}

	//Create the actual fyne list
	list := widget.NewList(
		func() int { return len(items) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) { obj.(*widget.Label).SetText(items[id]) },
	)

	//Add action when the list is selected
	list.OnSelected = func(id widget.ListItemID) {
		name := items[id]
		content.Objects = []fyne.CanvasObject{pages[name]}
		content.Refresh()
	}

	//Create horizontal split between the list and the contents of the page
	split := container.NewHSplit(list, content)
	split.Offset = 0.15

	//Set window content and size
	w.SetContent(split)
	w.Resize(fyne.NewSize(1000, 1000))

	//Run the gui
	w.ShowAndRun()
}

func (gui *GUI) CreateHomePage() fyne.CanvasObject {
	return widget.NewRichTextFromMarkdown(gui.text.GetHomePage())
}

func (gui *GUI) CreateGenerateSentencePage() fyne.CanvasObject {
	title := widget.NewRichTextFromMarkdown("# Generate sentences")
	word := widget.NewEntry()
	word.Validator = func(s string) error {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("word is required")
		}
		return nil
	}
	translationHint := widget.NewEntry()

	languages := gui.text.GetLanguageList()
	wordLangSelector := widget.NewSelect(languages, nil)
	translationLangSelector := widget.NewSelect(languages, nil)

	genders := gui.text.GetGenders()
	voiceGenderSelector := widget.NewSelect(genders, nil)
	voiceGenderSelector.SetSelected(gui.text.Female())
	voiceGenderSelector.Disable()

	includeAudio := widget.NewCheck("", func(b bool) {
		if b {
			voiceGenderSelector.Enable()
		} else {
			voiceGenderSelector.Disable()
		}
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: gui.text.Word(), Widget: word},
			{Text: gui.text.TranslationHint(), Widget: translationHint},
			{Text: gui.text.PickWordLang, Widget: wordLangSelector},
			{Text: gui.text.PickTranslationLang, Widget: translationLangSelector},
			{Text: gui.text.IncludeAudio, Widget: includeAudio},
			{Text: gui.text.VoiceGender, Widget: voiceGenderSelector},
		},

		OnSubmit: func() {
			if err := word.Validate(); err != nil {
				return
			}
			go gui.handleGenerateSentences(&GenerateSentencesParams{
				Word:            word.Text,
				TranslationHint: translationHint.Text,
				WordLang:        wordLangSelector.Selected,
				TranslationLang: translationLangSelector.Selected,
				IncludeAudio:    includeAudio.Checked,
			})
			word.SetText("")
			word.SetValidationError(nil)

			translationHint.SetText("")
		},
	}

	return container.NewVBox(title, form)
}

func (gui *GUI) CreateTranslatePage() fyne.CanvasObject {
	return widget.NewLabel("Translate")
}

func (gui *GUI) CreateGenerateDefinitionPage() fyne.CanvasObject {
	return widget.NewLabel("Definition")
}

func (gui *GUI) CreateSettingsPage() fyne.CanvasObject {
	return widget.NewLabel("Settings")
}
