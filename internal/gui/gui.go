package gui

import (
	"context"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dafraer/sentence-gen-grpc-client/internal/core"
	"github.com/dafraer/sentence-gen-grpc-client/internal/text"
	"go.uber.org/zap"
)

const (
	projectLink  = "https://github.com/dafraer/sentence-gen-grpc-client"
	screenWidth  = 720
	screenHeight = 480
)

type GUI struct {
	app    fyne.App
	window fyne.Window
	logger *zap.SugaredLogger
	core   *core.Core
	text   text.Text
}

// New creates new GUI
func New(logger *zap.SugaredLogger, core *core.Core, text text.Text) *GUI {
	return &GUI{
		logger: logger,
		core:   core,
		text:   text,
	}
}

// Run runs the app
func (gui *GUI) Run() {
	gui.logger.Infow("starting GUI")
	//Create new instance of the GUI app
	gui.app = app.New()

	//Create new window
	gui.window = gui.app.NewWindow("Sengen")

	//Items in the list
	items := gui.text.GetPageNames()

	//Create pages
	generateSentencePage := gui.createGenerateSentencePage()
	translatePage := gui.createTranslatePage()
	generateDefinitionPage := gui.createGenerateDefinitionPage()
	tutorialPage := gui.createTutorialPage()
	//Create variable for contents of the page
	content := container.NewStack(generateSentencePage) // right side

	//Create map of all pages
	pages := map[string]fyne.CanvasObject{
		items[0]: generateSentencePage,
		items[1]: translatePage,
		items[2]: generateDefinitionPage,
		items[3]: tutorialPage,
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

	//Wrap list in a container whose min width fits the longest item
	listWrapper := container.New(&minWidthLayout{width: listMinWidth(items)}, list)

	//Create horizontal split between the list and the contents of the page
	split := container.NewHSplit(listWrapper, content)
	split.Offset = 0

	//Set window content and size
	gui.window.SetContent(split)
	gui.window.Resize(fyne.NewSize(screenWidth, screenHeight))

	//Run the gui
	gui.window.ShowAndRun()
}

// createTutorialPage creates tutorial page
func (gui *GUI) createTutorialPage() fyne.CanvasObject {
	//Create page header
	header := widget.NewRichTextFromMarkdown(gui.text.TextTutorialTitle())

	//Parse project url
	u, err := url.Parse(projectLink)
	if err != nil {
		panic(err)
	}

	// Create the body of the page which is text + embedded link to the Github page
	body := widget.NewRichText(
		&widget.TextSegment{Text: gui.text.TextTutorialDescription() + " ", Style: widget.RichTextStyle{Inline: true}},
		&widget.HyperlinkSegment{Text: gui.text.TextTutorialLink(), URL: u},
	)

	return container.NewVBox(header, body)
}

// createGenerateSentencePage creates generate sentence page
func (gui *GUI) createGenerateSentencePage() fyne.CanvasObject {
	return gui.createTemplateTranslationPage(gui.text.TextGenerateSentenceTitle(), gui.onGenerateSentenceSubmit)
}

// createTranslatePage creates translate page
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
		gui.logger.Errorw("failed to fetch deck names", "err", err)
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

// createGenerateDefinitionPage creates generate definition page
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
		gui.logger.Errorw("failed to fetch deck names", "err", err)
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
