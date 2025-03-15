package types

import (
	"fmt"
	"strings"
)

// SendOptions are an options of the sending message.
type SendOptions struct {
	BotName  BotName
	ChatName ChatName

	Message MessageOptions
}

func (opt SendOptions) GetNormalized(defaultBotName BotName, defaultChatName ChatName) SendOptions {
	normalized := SendOptions{
		BotName:  opt.BotName,
		ChatName: opt.ChatName,
		Message:  opt.Message,
	}

	if normalized.BotName == "" {
		normalized.BotName = defaultBotName
	}

	if normalized.ChatName == "" {
		normalized.ChatName = defaultChatName
	}

	normalized.Message.Text = strings.TrimSpace(normalized.Message.Text)

	return normalized
}

func (opt SendOptions) Validate() error {
	if err := opt.Message.Validate(); err != nil {
		return fmt.Errorf("message content is invalid: %w", err)
	}

	return nil
}

// MessageOptions is a sending message content options.
type MessageOptions struct {
	Text      string
	ParseMode ParseMode

	DisableNotification bool
	ProtectContent      bool
}

func (m MessageOptions) Validate() error {
	if strings.TrimSpace(m.Text) == "" {
		return fmt.Errorf("no message text given")
	}

	return nil
}
