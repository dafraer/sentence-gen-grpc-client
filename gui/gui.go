package gui

import (
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
	items := gui.text.GetPageNames()

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
	//Create title
	title := widget.NewRichTextFromMarkdown(gui.text.TextGenerateSentence())

	//Create word entry
	wordEntry := widget.NewEntry()

	translationHint := widget.NewEntry()

	//Create language selectors
	languages := gui.text.GetLanguages()
	wordLangSelect := widget.NewSelect(languages, nil)

	translationLangSelect := widget.NewSelect(languages, nil)

	//Create gender
	genders := gui.text.GetGenders()
	voiceGenderSelect := widget.NewSelect(genders, nil)
	voiceGenderSelect.SetSelected(gui.text.TextFemale())
	voiceGenderSelect.Disable()

	//Create check to include audio
	audioCheck := widget.NewCheck("", func(b bool) {
		if b {
			voiceGenderSelect.Enable()
		} else {
			voiceGenderSelect.Disable()
		}
	})

	//Create the form
	form := gui.createGenerateSentenceForm(&generateSentenceFormParams{
		wordEntry,
		translationHint,
		wordLangSelect,
		translationLangSelect,
		voiceGenderSelect,
		audioCheck,
	})

	//Every time something changes, validate form and enable/disable it
	validateForm := func(s string) {
		if err := gui.validateGenerateSentenceForm(wordEntry.Text,
			translationHint.Text,
			wordLangSelect.Selected,
			translationLangSelect.Selected); err != nil {
			form.Disable()
			return
		}
		form.Enable()
	}
	wordLangSelect.OnChanged = validateForm
	translationLangSelect.OnChanged = validateForm
	wordEntry.OnChanged = validateForm
	translationHint.OnChanged = validateForm

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
