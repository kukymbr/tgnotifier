//go:build gui

package main

import (
	"github.com/kukymbr/tgnotifier/internal/gui"
	"github.com/kukymbr/tgnotifier/internal/tgnotifier"
	flag "github.com/spf13/pflag"
	"log"
	"os"
)

var genericOptions = tgnotifier.GenericOptions{}

func main() {
	initGenericFlags(os.Args)

	genericOptions.Normalize()

	ctn, err := tgnotifier.NewDefaultDependencyContainer(genericOptions.ConfigPath)
	if err != nil {
		gui.RunFailed(err)

		return
	}

	gui.Run(ctn.Config, ctn.Sender, ctn.Client)
}

func initGenericFlags(args []string) {
	flags := flag.NewFlagSet("", flag.ExitOnError)

	flags.SetOutput(os.Stderr)

	flags.StringVar(
		&genericOptions.ConfigPath,
		"config",
		"",
		"Path to a config file",
	)

	flags.BoolVar(
		&genericOptions.IsDebug,
		"debug",
		false,
		"Enable the debug mode",
	)

	if err := flags.Parse(args); err != nil {
		log.Fatalf("Error parsing flags: %s", err.Error())
	}
}
