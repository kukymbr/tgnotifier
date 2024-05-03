package tgnotifier

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kukymbr/tgnotifier/internal/sender"
	"github.com/kukymbr/tgnotifier/internal/types"
)

// DefaultConfigPath is a path to a config file by default.
const DefaultConfigPath = ".tgnotifier.yml"

// NewOptions creates new Options with default values.
func NewOptions() *Options {
	opt := &Options{}

	return opt
}

// Options is a tgnotifier CLI options.
type Options struct {
	ConfigPath string

	BotName  types.BotName
	ChatName types.ChatName

	Message sender.MessageOptions

	IsDebug bool
}

// Normalize normalizes the values and sets the default values if missing.
func (opt *Options) Normalize() {
	if opt.ConfigPath == "" {
		if _, err := os.Stat(DefaultConfigPath); err == nil {
			opt.ConfigPath = DefaultConfigPath
		}
	}

	if opt.ConfigPath != "" {
		opt.ConfigPath = filepath.Clean(opt.ConfigPath)
	}

	opt.Message.Text = strings.TrimSpace(opt.Message.Text)
}

// Validate validates if given values are valid.
func (opt *Options) Validate() error {
	if opt.Message.Text == "" {
		return fmt.Errorf("no message text given")
	}

	return nil
}
