package anki

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Client struct {
	logger          *zap.SugaredLogger
	ankiConnectAddr string
}

func NewClient(logger *zap.SugaredLogger, ankiConnectAddr string) *Client {
	return &Client{
		logger:          logger,
		ankiConnectAddr: ankiConnectAddr,
	}
}

func (c *Client) AddCard(ctx context.Context, note Note) error {
	nb := noteBody{
		DeckName:  note.DeckName,
		ModelName: string(note.CardType),
		Fields: map[string]string{
			"Front": note.Front,
			"Back":  note.Back,
		},
	}

	if note.Audio != nil {
		nb.Audio = []audioBody{{
			Data:     base64.StdEncoding.EncodeToString(note.Audio.Data),
			Filename: note.Audio.Filename,
			Fields:   note.Audio.Fields,
		}}
	}

	req := ankiRequest{
		Action:  "addNote",
		Version: 6,
		Params:  addNoteParams{Note: nb},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://"+c.ankiConnectAddr, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer func(resp *http.Response) {
		if err := resp.Body.Close(); err != nil {
			c.logger.Errorw("failed to close response body", "err", err)
		}
	}(resp)

	var ankiResp ankiResponse
	if err := json.NewDecoder(resp.Body).Decode(&ankiResp); err != nil {
		return err
	}

	if ankiResp.Error != nil {
		return fmt.Errorf("anki: %s", *ankiResp.Error)
	}

	return nil
}
