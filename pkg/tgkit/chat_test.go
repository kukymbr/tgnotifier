package tgkit_test

import (
	"testing"

	"github.com/kukymbr/tgnotifier/pkg/tgkit"
	"github.com/stretchr/testify/assert"
)

func TestChatID_String(t *testing.T) {
	tests := []struct {
		Input    string
		Expected string
	}{
		{Input: "100", Expected: "100"},
		{Input: "-123", Expected: "-123"},
		{Input: "test1", Expected: "@test1"},
		{Input: "@test2", Expected: "@test2"},
	}

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			chatID := tgkit.ChatID(test.Input)

			assert.Equal(t, chatID.String(), test.Expected)
		})
	}
}
