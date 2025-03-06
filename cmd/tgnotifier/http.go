//go:build !no_http

package main

import (
	"context"
	"github.com/kukymbr/tgnotifier/internal/server/http"
	"github.com/kukymbr/tgnotifier/internal/tgnotifier"
	"github.com/spf13/cobra"
)

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
