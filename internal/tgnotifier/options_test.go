package tgnotifier_test

import (
	"testing"

	"github.com/kukymbr/tgnotifier/internal/sender"
	"github.com/kukymbr/tgnotifier/internal/tgnotifier"
	"github.com/kukymbr/tgnotifier/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	tests := []struct {
		Name       string
		Options    tgnotifier.Options
		Normalized tgnotifier.Options
		Valid      bool
	}{
		{
			Name: "valid 1",
			Options: tgnotifier.Options{
				BotName:  "",
				ChatName: "",
				Message: sender.MessageOptions{
					Text: " test1 ",
				},
			},
			Normalized: tgnotifier.Options{
				ConfigPath: tgnotifier.DefaultConfigPath,
				BotName:    types.DefaultBotName,
				ChatName:   types.DefaultChatName,
				Message: sender.MessageOptions{
					Text: "test1",
				},
			},
			Valid: true,
		},
		{
			Name: "valid 2",
			Options: tgnotifier.Options{
				BotName:  "testBot",
				ChatName: "testChat",
				Message:  sender.MessageOptions{Text: "Test"},
			},
			Normalized: tgnotifier.Options{
				ConfigPath: tgnotifier.DefaultConfigPath,
				BotName:    "testBot",
				ChatName:   "testChat",
				Message:    sender.MessageOptions{Text: "Test"},
			},
			Valid: true,
		},
		{
			Name:    "invalid 1",
			Options: *tgnotifier.NewOptions(),
			Valid:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Options.Normalize()

			err := test.Options.Validate()

			if test.Valid {
				assert.EqualExportedValues(t, test.Normalized, test.Options)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
