package gui

func (gui *GUI) handleGenerateSentences(params *generateSentencesParams) {

	gui.logger.Infow("handler called", "data", *params)
}

func (gui *GUI) handleTranslation(params *translateParams) {
	gui.logger.Infow("handler called", "data", *params)
}

func (gui *GUI) handleGenerateDefinition(params *generateDefinitionParams) {
	gui.logger.Infow("handler called", "data", *params)
}
