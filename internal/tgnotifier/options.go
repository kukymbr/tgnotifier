package tgnotifier

import (
	"github.com/kukymbr/tgnotifier/internal/types"
	"os"
	"path/filepath"
)

// DefaultConfigPath is a path to a config file by default.
const DefaultConfigPath = ".tgnotifier.yml"

// NewOptions creates new Options with default values.
func NewOptions() Options {
	return Options{}
}

// GenericOptions are the options shared between commands.
type GenericOptions struct {
	ConfigPath string
	IsDebug    bool
}

// Normalize normalizes the values and sets the default values if missing.
func (opt *GenericOptions) Normalize() {
	if opt.ConfigPath == "" {
		if _, err := os.Stat(DefaultConfigPath); err == nil {
			opt.ConfigPath = DefaultConfigPath
		}
	}

	if opt.ConfigPath != "" {
		opt.ConfigPath = filepath.Clean(opt.ConfigPath)
	}
}

// Options is a tgnotifier CLI options.
type Options struct {
	GenericOptions GenericOptions
	SendOptions    types.SendOptions
}

// Normalize normalizes the values and sets the default values if missing.
func (opt *Options) Normalize() {
	opt.GenericOptions.Normalize()
}

// Validate validates if given values are valid.
func (opt *Options) Validate() error {
	return nil
}
