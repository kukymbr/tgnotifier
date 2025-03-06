package main

import (
	"context"
	"github.com/kukymbr/tgnotifier/internal/server/grpc"
	"github.com/kukymbr/tgnotifier/internal/server/http"
	"os"
	"os/signal"

	"github.com/kukymbr/tgnotifier/internal/config"
	"github.com/kukymbr/tgnotifier/internal/tgnotifier"
	"github.com/kukymbr/tgnotifier/internal/util"
	"github.com/spf13/cobra"
)

var (
	sendOptions    = tgnotifier.NewOptions()
	genericOptions = tgnotifier.GenericOptions{}
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cmd := getRootCommandDefinition(ctx)

	if err := cmd.ExecuteContext(ctx); err != nil {
		util.PrintlnError(err)

		os.Exit(1)
	}

	os.Exit(0)
}

func getRootCommandDefinition(ctx context.Context) *cobra.Command {
	cmd := getSendCommandDefinition(ctx, "tgnotifier")

	cmd.AddCommand(
		getSendCommandDefinition(ctx, "send"),
		getGRPCCommandDefinition(),
		getHTTPCommandDefinition(ctx),
	)

	return cmd
}

func getSendCommandDefinition(ctx context.Context, use string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     use,
		Short:   "Tool to send telegram notifications",
		Long:    `A tool send notifications via the Telegram HTTPS API.`,
		Version: tgnotifier.Version,

		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return tgnotifier.RunSendMessage(ctx, sendOptions)
		},
	}

	initSendFlags(cmd)

	return cmd
}

func getGRPCCommandDefinition() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "grpc",
		Short:   "tgnotifier gRPC server",
		Long:    `Starts an gRPC server to send notifications via the Telegram HTTPS API.`,
		Version: tgnotifier.Version,

		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: func(cmd *cobra.Command, args []string) error {
			genericOptions.Normalize()

			ctn, err := tgnotifier.NewDefaultDependencyContainer(genericOptions.ConfigPath)
			if err != nil {
				return err
			}

			server := grpc.New(ctn.Config, ctn.Sender)

			return server.Run()
		},
	}

	initGenericFlags(cmd, &genericOptions)

	return cmd
}

func getHTTPCommandDefinition(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "http",
		Short:   "tgnotifier HTTP server",
		Long:    `Starts an HTTP server to send notifications via the Telegram HTTPS API.`,
		Version: tgnotifier.Version,

		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: func(cmd *cobra.Command, args []string) error {
			genericOptions.Normalize()

			ctn, err := tgnotifier.NewDefaultDependencyContainer(genericOptions.ConfigPath)
			if err != nil {
				return err
			}

			return http.RunServer(ctx, ctn.Config, ctn.Sender)
		},
	}

	initGenericFlags(cmd, &genericOptions)

	return cmd
}

func initSendFlags(cmd *cobra.Command) {
	initGenericFlags(cmd, &sendOptions.GenericOptions)

	cmd.Flags().Var(
		&sendOptions.SendOptions.BotName,
		"bot",
		"Bot name to send message from (defined in config); "+
			"if not set, the default_bot directive or the bot "+
			"from the "+config.EnvDefaultBot+" env var will be used",
	)

	cmd.Flags().Var(
		&sendOptions.SendOptions.ChatName,
		"chat",
		"Chat name to send message to (defined in config); "+
			"if not set, the default_chat directive or the chat ID "+
			"from the "+config.EnvDefaultChat+" env var will be used",
	)

	cmd.Flags().StringVar(
		&sendOptions.SendOptions.Message.Text,
		"text",
		"",
		"Message text",
	)
	_ = cmd.MarkFlagRequired("text")

	cmd.Flags().Var(
		&sendOptions.SendOptions.Message.ParseMode,
		"parse-mode",
		"Parse mode (MarkdownV2|HTML)",
	)

	cmd.Flags().BoolVar(
		&sendOptions.SendOptions.Message.DisableNotification,
		"disable-notification",
		false,
		"Disable message sound notification",
	)

	cmd.Flags().BoolVar(
		&sendOptions.SendOptions.Message.ProtectContent,
		"protect-content",
		false,
		"Protect message content from copying and forwarding",
	)
}

func initGenericFlags(cmd *cobra.Command, opt *tgnotifier.GenericOptions) {
	cmd.Flags().StringVar(
		&opt.ConfigPath,
		"config",
		"",
		"Path to a config file",
	)

	cmd.Flags().BoolVar(
		&opt.IsDebug,
		"debug",
		false,
		"Enable the debug mode",
	)
}
