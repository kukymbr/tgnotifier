package tgkit_test

import (
	"testing"

	"github.com/kukymbr/tgnotifier/pkg/tgkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBot_WhenValid_ExpectNoError(t *testing.T) {
	tests := []struct {
		Identity         string
		ExpectedID       tgkit.BotID
		ExpectedToken    tgkit.BotToken
		ExpectedIdentity tgkit.BotIdentity
		ExpectedString   string
	}{
		{
			Identity:         "12345:testToken1",
			ExpectedID:       12345,
			ExpectedToken:    "testToken1",
			ExpectedIdentity: "bot12345:testToken1",
			ExpectedString:   "bot12345:*****",
		},
		{
			Identity:         "bot54321:testToken2",
			ExpectedID:       54321,
			ExpectedToken:    "testToken2",
			ExpectedIdentity: "bot54321:testToken2",
			ExpectedString:   "bot54321:*****",
		},
	}

	for _, test := range tests {
		t.Run(test.Identity, func(t *testing.T) {
			bot, err := tgkit.NewBot(test.Identity)

			require.NoError(t, err)

			assert.Equal(t, test.ExpectedID, bot.GetID())
			assert.Equal(t, test.ExpectedToken, bot.GetToken())
			assert.Equal(t, test.ExpectedIdentity, bot.GetIdentity())
			assert.Equal(t, test.ExpectedString, bot.String())

			assert.NotPanics(t, func() {
				tgkit.MustNewBot(test.Identity)
			})
		})
	}
}

func TestNewBot_WhenInvalid_ExpectError(t *testing.T) {
	tests := []struct {
		Identity string
	}{
		{"invalid_ident"},
		{"invalid_id:testToken1"},
		{"12345:invalid/token"},
	}

	for _, test := range tests {
		t.Run(test.Identity, func(t *testing.T) {
			bot, err := tgkit.NewBot(test.Identity)

			assert.Error(t, err)
			assert.Empty(t, bot.String())
		})
	}
}
