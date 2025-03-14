package sender

import (
	"context"
	"time"

	"github.com/kukymbr/tgnotifier/internal/config"
	"github.com/kukymbr/tgnotifier/internal/msgproc"
	"github.com/kukymbr/tgnotifier/internal/types"
	"github.com/kukymbr/tgnotifier/pkg/tgkit"
)

// New returns a new Sender instance.
func New(conf *config.Config, client *tgkit.Client, msgProc msgproc.MessageProcessor) *Sender {
	return &Sender{
		conf:    conf,
		client:  client,
		msgProc: msgProc,
	}
}

// Sender is a tool to send a message via the tgkit.Client.
type Sender struct {
	conf    *config.Config
	client  *tgkit.Client
	msgProc msgproc.MessageProcessor
}

// Send sends a message from the bot to the chat.
func (s *Sender) Send(ctx context.Context, opt types.SendOptions) (*tgkit.TgMessage, error) {
	if err := opt.Validate(); err != nil {
		return nil, err
	}

	bot, err := s.conf.Bots().FindByName(opt.BotName)
	if err != nil {
		return nil, err
	}

	chatID, err := s.conf.Chats().FindByName(opt.ChatName)
	if err != nil {
		return nil, err
	}

	opt.Message.Text = s.msgProc.Process(opt.Message.Text)

	disableNotification := opt.Message.DisableNotification
	if s.conf.GetSilenceSchedule().Has(time.Now()) {
		disableNotification = true
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return s.client.SendMessage(bot, tgkit.TgMessageRequest{
		ChatID:              chatID,
		Text:                opt.Message.Text,
		ParseMode:           opt.Message.ParseMode.String(),
		DisableNotification: disableNotification,
		ProtectContent:      opt.Message.ProtectContent,
	})
}
