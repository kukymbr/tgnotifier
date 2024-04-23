package tgkit_test

import (
	"net/http"
	"testing"

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
			assert.Nil(t, resp)
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
		Assert      func(t *testing.T, bot *tgkit.Bot, client *tgkit.Client)
	}{
		{
			BotIdentity: "1:test1",
			Assert: func(t *testing.T, bot *tgkit.Bot, client *tgkit.Client) {
				user, err := client.GetMe(bot)

				require.NoError(t, err)
				require.NotNil(t, user)

				assert.True(t, user.IsBot)
				assert.Equal(t, "Test Bot 1", user.FirstName)
			},
		},
		{
			BotIdentity: "2:test2",
			Assert: func(t *testing.T, bot *tgkit.Bot, client *tgkit.Client) {
				user, err := client.GetMe(bot)

				require.NoError(t, err)
				require.NotNil(t, user)

				assert.True(t, user.IsBot)
				assert.Equal(t, "Test Bot 2", user.FirstName)
			},
		},
		{
			BotIdentity: "3:test3",
			Assert: func(t *testing.T, bot *tgkit.Bot, client *tgkit.Client) {
				user, err := client.GetMe(bot)

				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.BotIdentity, func(t *testing.T) {
			client := tgkit.NewClient(httpClient)
			bot, err := tgkit.NewBot(test.BotIdentity)

			require.NoError(t, err)

			test.Assert(t, bot, client)
		})
	}
}
