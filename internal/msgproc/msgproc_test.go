package msgproc_test

import (
	"testing"

	"github.com/kukymbr/tgnotifier/internal/msgproc"
	"github.com/kukymbr/tgnotifier/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewMessageProcessor_Process(t *testing.T) {
	tests := []struct {
		Input     string
		ParseMode types.ParseMode
		Expected  string
	}{
		{
			Input:     " Hello, @testUser1! ",
			ParseMode: types.ParseModeMarkdown2,
			Expected:  "Hello, @testUser1!",
		},
		{
			Input:     " Hello, @testUser2! ",
			ParseMode: types.ParseModeHTML,
			Expected:  `Hello, @testUser2!`,
		},
	}

	proc := msgproc.NewMessageProcessor()

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			msg := proc.Process(test.Input)

			assert.Equal(t, test.Expected, msg)
		})
	}
}
