package sender

import (
	"context"

	"github.com/kukymbr/tgnotifier/internal/config"
	"github.com/kukymbr/tgnotifier/internal/types"
	"github.com/kukymbr/tgnotifier/pkg/tgkit"
)

// NewSender returns a new Sender instance.
func NewSender(conf *config.Config, client *tgkit.Client) *Sender {
	return &Sender{
		conf:   conf,
		client: client,
	}
}

// Sender is a tool to send a message via the tgkit.Client.
type Sender struct {
	conf   *config.Config
	client *tgkit.Client
}

// Send sends a message from the bot to the chat.
func (s *Sender) Send(
	ctx context.Context,
	botName types.BotName,
	chatName types.ChatName,
	msg MessageOptions,
) (*tgkit.TgMessage, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	bot, err := s.conf.Bots().GetBot(botName)
	if err != nil {
		return nil, err
	}

	chatID, err := s.conf.Chats().GetChatID(chatName)
	if err != nil {
		return nil, err
	}

	return s.client.SendMessage(bot, tgkit.TgMessageRequest{
		ChatID:              chatID,
		Text:                msg.Text,
		ParseMode:           msg.ParseMode,
		DisableNotification: msg.DisableNotification,
		ProtectContent:      msg.ProtectContent,
	})
}
