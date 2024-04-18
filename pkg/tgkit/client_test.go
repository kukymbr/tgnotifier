package tgkit_test

import (
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
