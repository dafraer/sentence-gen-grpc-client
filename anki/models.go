package anki

type CardType string

const (
	Basic           CardType = "Basic"
	BasicAndReverse CardType = "Basic (and reversed card)"
)

type Note struct {
	DeckName string
	CardType CardType
	Front    string
	Back     string
	Audio    *Audio
}

type Audio struct {
	Data     []byte
	Filename string
	Fields   []string
}

// AnkiConnect API types

type ankiRequest struct {
	Action  string      `json:"action"`
	Version int         `json:"version"`
	Params  interface{} `json:"params,omitempty"`
}

type addNoteParams struct {
	Note noteBody `json:"note"`
}

type noteBody struct {
	DeckName  string            `json:"deckName"`
	ModelName string            `json:"modelName"`
	Fields    map[string]string `json:"fields"`
	Audio     []audioBody       `json:"audio,omitempty"`
}

type audioBody struct {
	Data     string   `json:"data"`
	Filename string   `json:"filename"`
	Fields   []string `json:"fields"`
}

type ankiResponse struct {
	Result interface{} `json:"result"`
	Error  *string     `json:"error"`
}

type deckNamesResponse struct {
	Result []string `json:"result"`
	Error  *string  `json:"error"`
}
