package tgnotifier

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/kukymbr/tgnotifier/internal/util"
)

// RunSendMessage executes the message sending.
func RunSendMessage(ctx context.Context, opt Options) error {
	opt.Normalize()

	if err := opt.Validate(); err != nil {
		return fmt.Errorf("arguments are invalid: %w", err)
	}

	ctn, err := NewDefaultDependencyContainer(opt.GenericOptions.ConfigPath)
	if err != nil {
		return err
	}

	msg, err := ctn.Sender.Send(ctx, opt.SendOptions)
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
