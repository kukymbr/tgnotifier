package tgkit_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/kukymbr/tgnotifier/pkg/tgkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetMe_WhenInvalid_ExpectError(t *testing.T) {
	client := tgkit.NewDefaultClient()

	tests := []struct {
		Identity string
	}{
		{"1000:not_a_real_token"},
		{"0:test"},
	}

	for _, test := range tests {
		t.Run(test.Identity, func(t *testing.T) {
			bot, err := tgkit.NewBot(test.Identity)
			require.NoError(t, err)

			resp, err := client.GetMe(bot)

			assert.Error(t, err)
			assert.Empty(t, resp)
		})
	}
}

func TestClient_GetMe(t *testing.T) {
	httpClient := &tgkit.HTTPClientJSONMock{}
	httpClient.RegisterResponse(
		http.MethodGet,
		"https://api.telegram.org/bot1:test1/getMe",
		http.StatusOK,
		&tgkit.TgUserResponse{
			Ok: true,
			Result: tgkit.TgUser{
				IsBot:     true,
				FirstName: "Test Bot 1",
			},
		},
	)
	httpClient.RegisterResponse(
		http.MethodGet,
		"https://api.telegram.org/bot2:test2/getMe",
		http.StatusOK,
		&tgkit.TgUserResponse{
			Ok: true,
			Result: tgkit.TgUser{
				IsBot:     true,
				FirstName: "Test Bot 2",
			},
		},
	)
	httpClient.RegisterResponse(
		http.MethodGet,
		"https://api.telegram.org/bot3:test3/getMe",
		http.StatusNotFound,
		&tgkit.TgErrorResponse{
			Ok:          false,
			ErrorCode:   http.StatusNotFound,
			Description: "Not found",
		},
	)

	tests := []struct {
		BotIdentity string
		Assert      func(t *testing.T, user tgkit.TgUser, err error)
	}{
		{
			BotIdentity: "1:test1",
			Assert: func(t *testing.T, user tgkit.TgUser, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, user)

				assert.True(t, user.IsBot)
				assert.Equal(t, "Test Bot 1", user.FirstName)
			},
		},
		{
			BotIdentity: "2:test2",
			Assert: func(t *testing.T, user tgkit.TgUser, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, user)

				assert.True(t, user.IsBot)
				assert.Equal(t, "Test Bot 2", user.FirstName)
			},
		},
		{
			BotIdentity: "3:test3",
			Assert: func(t *testing.T, user tgkit.TgUser, err error) {
				assert.Error(t, err)
				assert.Empty(t, user)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.BotIdentity, func(t *testing.T) {
			client := tgkit.NewClientWithOptions(tgkit.WithHTTPClient(httpClient))
			bot, err := tgkit.NewBot(test.BotIdentity)

			require.NoError(t, err)

			user, err := client.GetMe(bot)

			test.Assert(t, user, err)
		})
	}
}

func TestClient_SendMessage(t *testing.T) {
	httpClient := &tgkit.HTTPClientJSONMock{}
	httpClient.RegisterResponse(
		http.MethodPost,
		"https://api.telegram.org/bot1:test1/sendMessage",
		http.StatusOK,
		&tgkit.TgMessageResponse{
			Ok: true,
			Result: tgkit.TgMessage{
				MessageID: 1,
			},
		},
	)
	httpClient.RegisterResponse(
		http.MethodPost,
		"https://api.telegram.org/bot2:test2/sendMessage",
		http.StatusNotFound,
		&tgkit.TgErrorResponse{
			Ok:          false,
			ErrorCode:   http.StatusNotFound,
			Description: "Not found",
		},
	)
	httpClient.RegisterResponse(
		http.MethodPost,
		"https://api.telegram.org/bot3:test3/sendMessage",
		http.StatusInternalServerError,
		"Internal Server Error",
	)

	tests := []struct {
		BotIdentity string
		ChatID      string
		Text        string
		Assert      func(t *testing.T, msg tgkit.TgMessage, err error)
	}{
		{
			BotIdentity: "1:test1",
			ChatID:      "1",
			Text:        "Test 1",
			Assert: func(t *testing.T, msg tgkit.TgMessage, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, msg)
				assert.Equal(t, 1, msg.MessageID)
			},
		},
		{
			BotIdentity: "2:test2",
			ChatID:      "1",
			Text:        "Test 2",
			Assert: func(t *testing.T, msg tgkit.TgMessage, err error) {
				assert.Error(t, err)
				assert.Empty(t, msg)
			},
		},
		{
			BotIdentity: "bot3:test3",
			ChatID:      "1",
			Text:        "Test 3",
			Assert: func(t *testing.T, msg tgkit.TgMessage, err error) {
				assert.Error(t, err)
				assert.Empty(t, msg)
				assert.Contains(t, err.Error(), "failed after 2 attempts")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.BotIdentity, func(t *testing.T) {
			t.Parallel()

			client := tgkit.NewClientWithOptions(
				tgkit.WithHTTPClient(httpClient),
				tgkit.WithRetrier(tgkit.NewProgressiveRetrier(2, 3*time.Millisecond, 2)),
			)
			bot, err := tgkit.NewBot(test.BotIdentity)

			require.NoError(t, err)

			msg, err := client.SendMessage(bot, tgkit.TgMessageRequest{
				ChatID: tgkit.ChatID(test.ChatID),
				Text:   test.Text,
			})

			test.Assert(t, msg, err)
		})
	}
}

func TestClient_GetUpdates(t *testing.T) {
	httpClient := &tgkit.HTTPClientJSONMock{}
	httpClient.RegisterResponse(
		http.MethodGet,
		"https://api.telegram.org/bot1:test1/getUpdates",
		http.StatusOK,
		&tgkit.TgUpdatesResponse{
			Ok: true,
			Result: []tgkit.TgUpdate{
				{
					UpdateID: 1,
					Message: tgkit.TgMessage{
						MessageID: 1,
					},
				},
			},
		},
	)
	httpClient.RegisterResponse(
		http.MethodGet,
		"https://api.telegram.org/bot2:test2/getUpdates",
		http.StatusOK,
		`{"ok":true, "result": []}`,
	)

	tests := []struct {
		BotIdentity string
		Assert      func(t *testing.T, updates []tgkit.TgUpdate, err error)
	}{
		{
			BotIdentity: "1:test1",
			Assert: func(t *testing.T, updates []tgkit.TgUpdate, err error) {
				require.NoError(t, err)
				require.Len(t, updates, 1)

				assert.Equal(t, updates[0].UpdateID, 1)
			},
		},
		{
			BotIdentity: "2:test2",
			Assert: func(t *testing.T, updates []tgkit.TgUpdate, err error) {
				require.NoError(t, err)

				assert.NotNil(t, updates)
				assert.Len(t, updates, 0)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.BotIdentity, func(t *testing.T) {
			client := tgkit.NewClientWithOptions(tgkit.WithHTTPClient(httpClient))
			bot, err := tgkit.NewBot(test.BotIdentity)

			require.NoError(t, err)

			updates, err := client.GetUpdates(bot)

			test.Assert(t, updates, err)
		})
	}
}
