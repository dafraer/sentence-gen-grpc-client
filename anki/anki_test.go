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

// newTestClient is a helper function that creates new client for testing
func newTestClient(t *testing.T) *Client {
	t.Helper()
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	//Create new client
	return NewClient(logger.Sugar(), "localhost:8765")
}

// setupTestDeck is a helper function that creates new test deck
func setupTestDeck(t *testing.T, client *Client) {
	t.Helper()

	//Build the request
	create := ankiRequest{
		Action:  "createDeck",
		Version: ankiConnectVersion,
		Params:  map[string]string{"deck": testDeckName},
	}

	//Make the request
	assert.NoError(t, doAnkiRequest(client, create))

	//Make a cleanup function
	t.Cleanup(func() {
		del := ankiRequest{
			Action:  "deleteDecks",
			Version: ankiConnectVersion,
			Params:  map[string]interface{}{"decks": []string{testDeckName}, "cardsToo": true},
		}
		assert.NoError(t, doAnkiRequest(client, del))
	})
}

// doAnkiRequest is a helper function that does the deck setup request
func doAnkiRequest(client *Client, req ankiRequest) error {
	//Marshal the request
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	//Make the request
	resp, err := http.Post("http://"+client.ankiConnectAddr, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//Parse the response
	var ankiResp ankiResponse
	if err := json.NewDecoder(resp.Body).Decode(&ankiResp); err != nil {
		return err
	}

	//Check for anki error
	if ankiResp.Error != nil {
		return assert.AnError
	}

	return nil
}

// TestAddBasicCard tests adding basic card
func TestAddBasicCard(t *testing.T) {
	//Create new client and test deck
	client := newTestClient(t)
	setupTestDeck(t, client)

	//Test the request
	err := client.AddCard(context.Background(), Note{
		DeckName: testDeckName,
		CardType: Basic,
		Front:    "Hello",
		Back:     "Привет",
	})

	assert.NoError(t, err)
}

// TestAddBasicAndReverseCard tests adding two-sided card
func TestAddBasicAndReverseCard(t *testing.T) {
	//Create new client and set up test deck
	client := newTestClient(t)
	setupTestDeck(t, client)

	//Add the card
	err := client.AddCard(context.Background(), Note{
		DeckName: testDeckName,
		CardType: BasicAndReverse,
		Front:    "Cat",
		Back:     "Кошка",
	})

	assert.NoError(t, err)
}
