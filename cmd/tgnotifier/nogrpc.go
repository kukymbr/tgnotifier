//go:build no_grpc

package main

import "github.com/spf13/cobra"

func getGRPCCommandDefinition() *cobra.Command {
	return nil
}
