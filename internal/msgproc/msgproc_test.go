package msgproc_test

import (
	"testing"

	"github.com/kukymbr/tgnotifier/internal/config"
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
			Input:     "Hello, @testUser1!",
			ParseMode: types.ParseModeMarkdown2,
			Expected:  "Hello, [@testUser1](tg://user?id=1)!",
		},
		{
			Input:     "Hello, @testUser2!",
			ParseMode: types.ParseModeHTML,
			Expected:  `Hello, <a href="tg://user?id=2">@testUser2</a>!`,
		},
		{
			Input:     "Hello, @testUser3!",
			ParseMode: types.ParseModeDefault,
			Expected:  "Hello, [@testUser3](tg://user?id=3)!",
		},
		{
			Input:     "Hello, @testUser4!",
			ParseMode: types.ParseModeMarkdown2,
			Expected:  "Hello, @testUser4!",
		},
		{
			Input:     "Hello, @testUser5!",
			ParseMode: types.ParseModeHTML,
			Expected:  "Hello, @testUser5!",
		},
		{
			Input:     "Hello, @testUser6!",
			ParseMode: types.ParseModeDefault,
			Expected:  "Hello, @testUser6!",
		},
		{
			Input:     "Hello, @testUser7!",
			ParseMode: "unknown",
			Expected:  "Hello, @testUser7!",
		},
	}

	users := config.UsersIndex{
		"testUser1": 1,
		"testUser2": 2,
		"testUser3": 3,
	}

	proc := msgproc.NewMessageProcessor(users)

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			msg := proc.Process(test.Input, test.ParseMode)

			assert.Equal(t, test.Expected, msg)
		})
	}
}

func TestNewMessageProcessor_Process_WhenNilUsers_ExpectNoPanic(t *testing.T) {
	proc := msgproc.NewMessageProcessor(nil)
	msg := ""
	inp := "Hello, @user!"

	assert.NotPanics(t, func() {
		msg = proc.Process(inp, types.ParseModeMarkdown2)
	})

	assert.Equal(t, inp, msg)
}
