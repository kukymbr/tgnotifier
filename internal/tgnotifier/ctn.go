package tgnotifier

import (
	"github.com/kukymbr/tgnotifier/internal/config"
	"github.com/kukymbr/tgnotifier/internal/sender"
	"github.com/kukymbr/tgnotifier/pkg/tgkit"
)

type DependencyContainer struct {
	Config *config.Config
	Client *tgkit.Client
	Sender *sender.Sender
}
