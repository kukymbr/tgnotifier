//go:build no_http

package main

import (
	"context"
	"github.com/spf13/cobra"
)

func getHTTPCommandDefinition(_ context.Context) *cobra.Command {
	return nil
}
