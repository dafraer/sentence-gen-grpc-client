package rpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func newTestClient(t *testing.T) *Client {
	t.Helper()
	logger, err := zap.NewDevelopment()
	assert.NoError(t, err)

	client, err := NewClient("localhost:50051", logger.Sugar())
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, client.Close())
	})

	return client
}

func TestGenerateSentence(t *testing.T) {
	client := newTestClient(t)

	resp, err := client.GenerateSentence(context.Background(), &GenerateSentenceRequest{
		Word:                "cat",
		WordLanguage:        "en-US",
		TranslationLanguage: "ru-RU",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.OriginalSentence)
	assert.NotEmpty(t, resp.TranslatedSentence)
}

func TestGenerateSentenceWithAudio(t *testing.T) {
	client := newTestClient(t)

	resp, err := client.GenerateSentence(context.Background(), &GenerateSentenceRequest{
		Word:                "cat",
		WordLanguage:        "en-US",
		TranslationLanguage: "ru-RU",
		IncludeAudio:        true,
		VoiceGender:         Female,
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.OriginalSentence)
	assert.NotEmpty(t, resp.TranslatedSentence)
	assert.NotEmpty(t, resp.Audio)

}

func TestGenerateSentenceWithHint(t *testing.T) {
	client := newTestClient(t)

	resp, err := client.GenerateSentence(context.Background(), &GenerateSentenceRequest{
		Word:                "run",
		WordLanguage:        "en-US",
		TranslationLanguage: "ru-RU",
		TranslationHint:     "use past tense",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.OriginalSentence)
	assert.NotEmpty(t, resp.TranslatedSentence)
}
