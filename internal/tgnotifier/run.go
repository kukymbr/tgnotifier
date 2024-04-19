package tgnotifier

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/kukymbr/tgnotifier/internal/config"
	"github.com/kukymbr/tgnotifier/internal/sender"
	"github.com/kukymbr/tgnotifier/internal/util"
	"github.com/kukymbr/tgnotifier/pkg/tgkit"
)

// Run executes the message sending.
func Run(ctx context.Context, opt *Options) error {
	ctn, err := getCtn(opt)
	if err != nil {
		return err
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

func getCtn(opt *Options) (*DependencyContainer, error) {
	if opt == nil {
		panic("nil options given")
	}

	opt.Normalize()

	if err := opt.Validate(); err != nil {
		return nil, fmt.Errorf("arguments are invalid: %w", err)
	}

	conf, err := config.NewConfigFromFile(opt.ConfigPath)
	if err != nil {
		return nil, err
	}

	client := tgkit.NewDefaultClient()

	return &DependencyContainer{
		Config: conf,
		Client: client,
		Sender: sender.NewSender(conf, client),
	}, nil
}
