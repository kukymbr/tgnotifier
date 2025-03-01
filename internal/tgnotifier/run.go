package tgnotifier

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/kukymbr/tgnotifier/internal/util"
)

// Run executes the message sending.
func Run(ctx context.Context, opt Options) error {
	opt.Normalize()

	if err := opt.Validate(); err != nil {
		return fmt.Errorf("arguments are invalid: %w", err)
	}

	ctn, err := NewDefaultDependencyContainer(opt.ConfigPath)
	if err != nil {
		return err
	}

	if opt.BotName == "" {
		opt.BotName = ctn.Config.GetDefaultBotName()
	}

	if opt.ChatName == "" {
		opt.ChatName = ctn.Config.GetDefaultChatName()
	}

	msg, err := ctn.Sender.Send(ctx, opt.BotName, opt.ChatName, opt.Message)
	if err != nil {
		return err
	}

	data, err := jsoniter.Marshal(msg)
	if err != nil {
		util.PrintlnError(fmt.Errorf("failed to encode response: %w", err))
	}

	fmt.Print(string(data))

	return nil
}
