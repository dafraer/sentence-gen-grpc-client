package anki

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const testDeckName = "TestDeck"

func newTestClient(t *testing.T) *Client {
	t.Helper()
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)
	return NewClient(logger.Sugar(), "localhost:8765")
}

func setupTestDeck(t *testing.T, client *Client) {
	t.Helper()

	create := ankiRequest{
		Action:  "createDeck",
		Version: 6,
		Params:  map[string]string{"deck": testDeckName},
	}
	assert.NoError(t, doAnkiRequest(client, create))

	t.Cleanup(func() {
		del := ankiRequest{
			Action:  "deleteDecks",
			Version: 6,
			Params:  map[string]interface{}{"decks": []string{testDeckName}, "cardsToo": true},
		}
		assert.NoError(t, doAnkiRequest(client, del))
	})
}

func doAnkiRequest(client *Client, req ankiRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://"+client.ankiConnectAddr, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var ankiResp ankiResponse
	if err := json.NewDecoder(resp.Body).Decode(&ankiResp); err != nil {
		return err
	}

	if ankiResp.Error != nil {
		return assert.AnError
	}

	return nil
}

func TestAddBasicCard(t *testing.T) {
	client := newTestClient(t)
	setupTestDeck(t, client)

	err := client.AddCard(context.Background(), Note{
		DeckName: testDeckName,
		CardType: Basic,
		Front:    "Hello",
		Back:     "Привет",
	})

	assert.NoError(t, err)
}

func TestAddBasicAndReverseCard(t *testing.T) {
	client := newTestClient(t)
	setupTestDeck(t, client)

	err := client.AddCard(context.Background(), Note{
		DeckName: testDeckName,
		CardType: BasicAndReverse,
		Front:    "Cat",
		Back:     "Кошка",
	})

	assert.NoError(t, err)
}
