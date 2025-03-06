//go:build !no_grpc

package main

import (
	"github.com/kukymbr/tgnotifier/internal/server/grpc"
	"github.com/kukymbr/tgnotifier/internal/tgnotifier"
	"github.com/spf13/cobra"
)

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
