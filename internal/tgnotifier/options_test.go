package tgnotifier_test

import (
	"github.com/kukymbr/tgnotifier/internal/types"
	"testing"

	"github.com/kukymbr/tgnotifier/internal/tgnotifier"
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
				SendOptions: types.SendOptions{
					BotName:  "",
					ChatName: "",
					Message: types.MessageOptions{
						Text: "test1",
					},
				},
			},
			Normalized: tgnotifier.Options{
				GenericOptions: tgnotifier.GenericOptions{
					ConfigPath: tgnotifier.DefaultConfigPath,
				},
				SendOptions: types.SendOptions{
					Message: types.MessageOptions{
						Text: "test1",
					},
				},
			},
			Valid: true,
		},
		{
			Name: "valid 2",
			Options: tgnotifier.Options{
				SendOptions: types.SendOptions{
					BotName:  "testBot",
					ChatName: "testChat",
					Message:  types.MessageOptions{Text: "Test"},
				},
			},
			Normalized: tgnotifier.Options{
				GenericOptions: tgnotifier.GenericOptions{
					ConfigPath: tgnotifier.DefaultConfigPath,
				},
				SendOptions: types.SendOptions{
					BotName:  "testBot",
					ChatName: "testChat",
					Message:  types.MessageOptions{Text: "Test"},
				},
			},
			Valid: true,
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
