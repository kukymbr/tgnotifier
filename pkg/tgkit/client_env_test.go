//go:build envtest

// Tests in this file expect environment values to be set:
// - TGKIT_ENVTEST_BOT_IDENTITY: identity of the real telegram bot to run the success tests;
// - TGKIT_ENVTEST_CHAT_ID: test chat ID.
package tgkit_test

import (
	"os"
	"testing"
	"time"

	"github.com/kukymbr/tgnotifier/pkg/tgkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func requireBot(t *testing.T) *tgkit.Bot {
	ident := os.Getenv("TGKIT_ENVTEST_BOT_IDENTITY")

	bot, err := tgkit.NewBot(ident)

	require.NoError(t, err)

	return bot
}

func requireChatID(t *testing.T) tgkit.ChatID {
	id := os.Getenv("TGKIT_ENVTEST_CHAT_ID")

	require.NotEmpty(t, id)

	return tgkit.ChatID(id)
}

func TestClient_GetMe(t *testing.T) {
	bot := requireBot(t)

	client := tgkit.NewDefaultClient()

	resp, err := client.GetMe(bot)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.True(t, resp.IsBot)
}

func TestClient_SendMessage(t *testing.T) {
	bot := requireBot(t)
	chatID := requireChatID(t)

	client := tgkit.NewDefaultClient()

	resp, err := client.SendMessage(bot, tgkit.TgMessageRequest{
		ChatID: chatID,
		Text:   "Test text " + time.Now().Format("2006-01-02 15:04:05"),
	})

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.NotEmpty(t, resp.MessageID)
}
