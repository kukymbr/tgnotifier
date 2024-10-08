package main

import (
	"context"
	"os"

	"github.com/kukymbr/tgnotifier/internal/config"
	"github.com/kukymbr/tgnotifier/internal/tgnotifier"
	"github.com/kukymbr/tgnotifier/internal/util"
	"github.com/spf13/cobra"
)

var opt = tgnotifier.NewOptions()

func main() {
	ctx := context.Background()
	cmd := getRootCommandDefinition(ctx)

	if err := cmd.ExecuteContext(ctx); err != nil {
		util.PrintlnError(err)

		os.Exit(1)
	}

	os.Exit(0)
}

func getRootCommandDefinition(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tgnotifier",
		Short:   "Tool to send telegram notifications",
		Long:    `A tool send notifications via the Telegram HTTPS API.`,
		Version: version,

		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return tgnotifier.Run(ctx, opt)
		},
	}

	initFlags(cmd)

	return cmd
}

func initFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&opt.ConfigPath,
		"config",
		"",
		"Path to a config file",
	)

	cmd.Flags().Var(
		&opt.BotName,
		"bot",
		"Bot name to send message from (defined in config); "+
			"if not set, the default_bot directive or the bot "+
			"from the "+config.EnvDefaultBot+" env var will be used",
	)

	cmd.Flags().Var(
		&opt.ChatName,
		"chat",
		"Chat name to send message to (defined in config); "+
			"if not set, the default_chat directive or the chat ID "+
			"from the "+config.EnvDefaultChat+" env var will be used",
	)

	cmd.Flags().StringVar(
		&opt.Message.Text,
		"text",
		"",
		"Message text",
	)
	_ = cmd.MarkFlagRequired("text")

	cmd.Flags().Var(
		&opt.Message.ParseMode,
		"parse-mode",
		"Parse mode (MarkdownV2|HTML)",
	)

	cmd.Flags().BoolVar(
		&opt.Message.DisableNotification,
		"disable-notification",
		false,
		"Disable message sound notification",
	)

	cmd.Flags().BoolVar(
		&opt.Message.ProtectContent,
		"protect-content",
		false,
		"Protect message content from copying and forwarding",
	)

	cmd.Flags().BoolVar(
		&opt.IsDebug,
		"debug",
		false,
		"Enable the debug mode",
	)
}
