//go:build gui

package main

import (
	"github.com/kukymbr/tgnotifier/internal/buildinfo"
	"github.com/kukymbr/tgnotifier/internal/gui"
	"github.com/kukymbr/tgnotifier/internal/tgnotifier"
	flag "github.com/spf13/pflag"
	"log"
	"os"
)

var (
	genericOptions = tgnotifier.GenericOptions{}
	isVersion      = false
)

func main() {
	initGenericFlags(os.Args)

	if isVersion {
		buildinfo.PrintVersion()

		os.Exit(0)
	}

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

	flags.BoolVar(&isVersion, "version", false, "Show tgnotifier version")

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
