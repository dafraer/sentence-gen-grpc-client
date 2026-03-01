package gui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dafraer/sentence-gen-grpc-client/core"
	"github.com/dafraer/sentence-gen-grpc-client/text"
	"go.uber.org/zap"
)

type GUI struct {
	app    fyne.App
	window fyne.Window
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
	gui.app = app.New()

	//Create new window
	gui.window = gui.app.NewWindow("Sengen")

	//Items in the list
	items := gui.text.GetPageNames()

	//Create pages
	homePage := gui.createHomePage()
	generateSentencePage := gui.createGenerateSentencePage()
	translatePage := gui.createTranslatePage()
	generateDefinitionPage := gui.createGenerateDefinitionPage()
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
	gui.window.SetContent(split)
	gui.window.Resize(fyne.NewSize(1000, 1000))

	//Run the gui
	gui.window.ShowAndRun()
}

func (gui *GUI) createHomePage() fyne.CanvasObject {
	return widget.NewRichTextFromMarkdown(gui.text.GetHomePage())
}

func (gui *GUI) createGenerateSentencePage() fyne.CanvasObject {
	return gui.createTemplateTranslationPage(gui.text.TextGenerateSentenceTitle(), gui.onGenerateSentenceSubmit)
}

func (gui *GUI) createTranslatePage() fyne.CanvasObject {
	return gui.createTemplateTranslationPage(gui.text.TextGenerateTranslationTitle(), gui.onTranslateSubmit)
}

// createTemplateTranslationPage creates either translate or generate sentence page, depending on params
func (gui *GUI) createTemplateTranslationPage(title string, onSubmit func(params *onSubmitParams)) fyne.CanvasObject {
	//Create header
	header := widget.NewRichTextFromMarkdown(title)

	//Create word entry
	wordEntry := widget.NewEntry()

	translationHint := widget.NewEntry()

	//Create language selectors
	languages := gui.text.GetLanguages()
	wordLangSelect := widget.NewSelect(languages, nil)

	translationLangSelect := widget.NewSelect(languages, nil)

	//Create voice selector
	voiceSelect := gui.createVoicePicker()

	//Create check to include audio
	audioCheck := widget.NewCheck("", func(b bool) {
		if b {
			voiceSelect.Enable()
		} else {
			voiceSelect.Disable()
		}
	})

	//Create deck selector populated from Anki
	decks, err := gui.core.GetDeckNames(context.Background())
	if err != nil {
		gui.logger.Warnw("failed to fetch deck names", "err", err)
	}
	deckSelect := widget.NewSelect(decks, nil)

	//Create the form
	form := gui.createForm(&formParams{
		word:            wordEntry,
		translationHint: translationHint,
		wordLang:        wordLangSelect,
		translationLang: translationLangSelect,
		voice:           voiceSelect,
		audio:           audioCheck,
		deck:            deckSelect,
		onSubmit:        onSubmit,
	})

	//Every time something changes, validate form and enable/disable it
	validateForm := func(s string) {
		if err := gui.validateForm(wordEntry.Text,
			translationHint.Text,
			wordLangSelect.Selected,
			translationLangSelect.Selected,
			deckSelect.Selected); err != nil {
			form.Disable()
			return
		}
		form.Enable()
	}
	wordLangSelect.OnChanged = validateForm
	translationLangSelect.OnChanged = validateForm
	wordEntry.OnChanged = validateForm
	translationHint.OnChanged = validateForm
	deckSelect.OnChanged = validateForm

	return container.NewVBox(header, form)
}

func (gui *GUI) createGenerateDefinitionPage() fyne.CanvasObject {
	//Create title
	title := widget.NewRichTextFromMarkdown(gui.text.TextGenerateDefinitionTitle())

	//Create word entry
	wordEntry := widget.NewEntry()

	definitionHint := widget.NewEntry()

	//Create language selectors
	languages := gui.text.GetLanguages()
	wordLangSelect := widget.NewSelect(languages, nil)

	//Create voice selector
	voiceSelect := gui.createVoicePicker()

	//Create check to include audio
	audioCheck := widget.NewCheck("", func(b bool) {
		if b {
			voiceSelect.Enable()
		} else {
			voiceSelect.Disable()
		}
	})

	//Create deck selector populated from Anki
	decks, err := gui.core.GetDeckNames(context.Background())
	if err != nil {
		gui.logger.Warnw("failed to fetch deck names", "err", err)
	}
	deckSelect := widget.NewSelect(decks, nil)

	//Create the form
	form := gui.createDefinitionForm(&definitionFormParams{
		word:           wordEntry,
		definitionHint: definitionHint,
		wordLang:       wordLangSelect,
		voice:          voiceSelect,
		audio:          audioCheck,
		deck:           deckSelect,
	})

	//Every time something changes, validate form and enable/disable it
	validateForm := func(s string) {
		if err := gui.validateDefinitionForm(wordEntry.Text,
			definitionHint.Text,
			wordLangSelect.Selected,
			deckSelect.Selected); err != nil {
			form.Disable()
			return
		}
		form.Enable()
	}
	wordLangSelect.OnChanged = validateForm
	wordEntry.OnChanged = validateForm
	definitionHint.OnChanged = validateForm
	deckSelect.OnChanged = validateForm

	return container.NewVBox(title, form)
}

func (gui *GUI) CreateSettingsPage() fyne.CanvasObject {
	//TODO: finish when will add settings package
	//Create header
	header := widget.NewRichTextFromMarkdown(gui.text.TextSettingsTitle())
	languages := gui.text.GetLanguages()
	languageSelector := widget.NewSelect(languages, nil)
	return container.NewVBox(header, languageSelector)
}
