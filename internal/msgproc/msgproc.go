package msgproc

import (
	"fmt"
	"strings"

	"github.com/kukymbr/tgnotifier/internal/config"
	"github.com/kukymbr/tgnotifier/internal/types"
)

// NewMessageProcessor returns new MessageProcessor instance.
func NewMessageProcessor(users config.UsersIndex) *MessageProcessor {
	return &MessageProcessor{
		users: users,
	}
}

// MessageProcessor is a tool to process messages before sending.
type MessageProcessor struct {
	users config.UsersIndex
}

// Process a message.
func (mp *MessageProcessor) Process(msg string, parseMode types.ParseMode) string {
	msg = mp.tagUsers(msg, parseMode)

	return msg
}

func (mp *MessageProcessor) tagUsers(msg string, parseMode types.ParseMode) string {
	if len(mp.users) == 0 {
		return msg
	}

	for name, id := range mp.users {
		find := "@" + name.String()
		repl := find

		switch parseMode {
		case types.ParseModeHTML:
			repl = fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, id, find)
		case types.ParseModeMarkdown2, types.ParseModeDefault:
			repl = fmt.Sprintf(`[%s](tg://user?id=%d)`, find, id)
		default:
			continue
		}

		msg = strings.ReplaceAll(msg, find, repl)
	}

	return msg
}
