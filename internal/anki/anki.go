package anki

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

const (
	ankiConnectVersion = 6
	actionDeckNames    = "deckNames"
	actionAddNote      = "addNote"
)

var (
	ErrAnkiError = errors.New("anki connect error")
)

type Client struct {
	logger          *zap.SugaredLogger
	ankiConnectAddr string
}

// NewClient creates a new Anki Connect client
func NewClient(logger *zap.SugaredLogger, ankiConnectAddr string) *Client {
	return &Client{
		logger:          logger,
		ankiConnectAddr: ankiConnectAddr,
	}
}

// GetDeckNames gets deck names from the Anki Connect
func (c *Client) GetDeckNames(ctx context.Context) ([]string, error) {
	//Build the request
	req := ankiRequest{
		Action:  actionDeckNames,
		Version: ankiConnectVersion,
	}

	//Make the request
	resp, err := c.makeRequest(ctx, req)
	if err != nil {
		return nil, wrapError(err)
	}

	//Close the response body
	defer func(resp *http.Response) {
		if err := resp.Body.Close(); err != nil {
			c.logger.Errorw("failed to close response body", "err", err)
		}
	}(resp)

	//Parse the response
	var ankiResp deckNamesResponse
	if err := json.NewDecoder(resp.Body).Decode(&ankiResp); err != nil {
		c.logger.Errorw("failed to decode response", "err", err)
		return nil, wrapError(err)
	}

	//Check anki connect error
	if ankiResp.Error != nil {
		c.logger.Errorw("AnkiConnect returned error for GetDeckNames", "ankiErr", *ankiResp.Error)
		return nil, wrapError(fmt.Errorf("anki: %s", *ankiResp.Error))
	}

	return ankiResp.Result, nil
}

// AddCard adds note to the specified deck
func (c *Client) AddCard(ctx context.Context, note Note) error {
	c.logger.Debugw("adding card to Anki", "deck", note.DeckName, "front", note.Front, "hasAudio", note.Audio != nil)
	//Build note body
	nb := noteBody{
		DeckName:  note.DeckName,
		ModelName: string(note.CardType),
		Fields: map[string]string{
			"Front": note.Front,
			"Back":  note.Back,
		},
	}

	//Add the audio
	if note.Audio != nil {
		nb.Audio = []audioBody{{
			Data:     base64.StdEncoding.EncodeToString(note.Audio.Data),
			Filename: note.Audio.Filename,
			Fields:   note.Audio.Fields,
		}}
	}

	//Build the request
	req := ankiRequest{
		Action:  actionAddNote,
		Version: ankiConnectVersion,
		Params:  addNoteParams{Note: nb},
	}

	//Make the request
	resp, err := c.makeRequest(ctx, req)
	if err != nil {
		return wrapError(err)
	}

	//Close the response body
	defer func(resp *http.Response) {
		if err := resp.Body.Close(); err != nil {
			c.logger.Errorw("failed to close response body", "err", err)
		}
	}(resp)

	//Parse the response
	var ankiResp ankiResponse
	if err := json.NewDecoder(resp.Body).Decode(&ankiResp); err != nil {
		c.logger.Errorw("failed to decode response", "err", err)
		return wrapError(err)
	}

	//Check for anki error
	if ankiResp.Error != nil {
		c.logger.Errorw("AnkiConnect returned error for AddCard", "deck", note.DeckName, "ankiErr", *ankiResp.Error)
		return wrapError(fmt.Errorf("anki: %s", *ankiResp.Error))
	}

	c.logger.Infow("card added successfully", "deck", note.DeckName, "front", note.Front)
	return nil
}

func (c *Client) makeRequest(ctx context.Context, req ankiRequest) (*http.Response, error) {
	//Marshal json body
	body, err := json.Marshal(req)
	if err != nil {
		c.logger.Errorw("failed to marshal request", "err", err)
		return nil, err
	}

	//Create http request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://"+c.ankiConnectAddr, bytes.NewReader(body))
	if err != nil {
		c.logger.Errorw("failed to build AddCard HTTP request", "err", err)
		return nil, err
	}

	//Set the header
	httpReq.Header.Set("Content-Type", "application/json")

	//Make the request
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		c.logger.Errorw("HTTP request failed", "err", err)
		return nil, err
	}

	return resp, nil
}

// wrapError joins error with ErrAnkiError for later identification
func wrapError(err error) error {
	return errors.Join(err, ErrAnkiError)
}
