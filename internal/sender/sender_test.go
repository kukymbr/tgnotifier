package sender_test

import (
	"context"
	"github.com/kukymbr/tgnotifier/internal/config"
	"github.com/kukymbr/tgnotifier/internal/msgproc"
	"github.com/kukymbr/tgnotifier/internal/sender"
	"github.com/kukymbr/tgnotifier/internal/types"
	"github.com/kukymbr/tgnotifier/pkg/tgkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func getTestConfig() string {
	return `
bots:
  test_bot: "12345:test_token"
chats:
  test_chat: "-12345"
default_bot: "test_bot"
default_chat: "test_chat"
`
}

func TestSender(t *testing.T) {
	tests := []struct {
		Name        string
		SetupClient func(t *testing.T, httpMock *tgkit.HTTPClientJSONMock)
		SendOptions types.SendOptions
		Assert      func(t *testing.T, resp *tgkit.TgMessage, err error)
	}{
		{
			Name: "when ok (default bot)",
			SetupClient: func(t *testing.T, httpMock *tgkit.HTTPClientJSONMock) {
				httpMock.RegisterResponse(
					http.MethodPost,
					"https://api.telegram.org/bot12345:test_token/sendMessage",
					http.StatusOK,
					&tgkit.TgMessageResponse{
						Ok: true,
						Result: tgkit.TgMessage{
							MessageID: 1,
						},
					},
				)
			},
			SendOptions: types.SendOptions{
				Message: types.MessageOptions{
					Text: "test1",
				},
			},
			Assert: func(t *testing.T, resp *tgkit.TgMessage, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)

				assert.Equal(t, 1, resp.MessageID)
			},
		},
		{
			Name: "when ok (specified bot)",
			SetupClient: func(t *testing.T, httpMock *tgkit.HTTPClientJSONMock) {
				httpMock.RegisterResponse(
					http.MethodPost,
					"https://api.telegram.org/bot12345:test_token/sendMessage",
					http.StatusOK,
					&tgkit.TgMessageResponse{
						Ok: true,
						Result: tgkit.TgMessage{
							MessageID: 1,
						},
					},
				)
			},
			SendOptions: types.SendOptions{
				BotName: "test_bot",
				Message: types.MessageOptions{
					Text: "test1",
				},
			},
			Assert: func(t *testing.T, resp *tgkit.TgMessage, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)

				assert.Equal(t, 1, resp.MessageID)
			},
		},
		{
			Name: "when non-ok response",
			SetupClient: func(t *testing.T, httpMock *tgkit.HTTPClientJSONMock) {
				httpMock.RegisterResponse(
					http.MethodPost,
					"https://api.telegram.org/bot12345:test_token/sendMessage",
					http.StatusBadRequest,
					&tgkit.TgErrorResponse{
						Ok:          false,
						ErrorCode:   http.StatusBadRequest,
						Description: "something went wrong",
					},
				)
			},
			SendOptions: types.SendOptions{
				Message: types.MessageOptions{
					Text: "test1",
				},
			},
			Assert: func(t *testing.T, resp *tgkit.TgMessage, err error) {
				assert.Error(t, err)
				assert.Nil(t, resp)
			},
		},
		{
			Name: "when bot invalid",
			SendOptions: types.SendOptions{
				BotName: "unknown_bot",
				Message: types.MessageOptions{
					Text: "test",
				},
			},
			Assert: func(t *testing.T, resp *tgkit.TgMessage, err error) {
				assert.Error(t, err)
			},
		},
		{
			Name:        "when message invalid",
			SendOptions: types.SendOptions{},
			Assert: func(t *testing.T, resp *tgkit.TgMessage, err error) {
				assert.Error(t, err)
			},
		},
		{
			Name: "when chat invalid",
			SendOptions: types.SendOptions{
				ChatName: "unknown_chat",
				Message: types.MessageOptions{
					Text: "test",
				},
			},
			Assert: func(t *testing.T, resp *tgkit.TgMessage, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			conf, err := config.New(config.FromString(getTestConfig()))

			require.NoError(t, err)
			require.NotNil(t, conf)

			httpMock := &tgkit.HTTPClientJSONMock{}

			if test.SetupClient != nil {
				test.SetupClient(t, httpMock)
			}

			client := tgkit.NewClientWithOptions(tgkit.WithHTTPClient(httpMock))
			service := sender.New(conf, client, msgproc.NewProcessingChain())

			require.NotNil(t, service)

			res, err := service.Send(context.Background(), test.SendOptions)

			test.Assert(t, res, err)
		})
	}
}
