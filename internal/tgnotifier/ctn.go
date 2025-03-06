package tgnotifier

import (
	"github.com/kukymbr/tgnotifier/internal/config"
	"github.com/kukymbr/tgnotifier/internal/msgproc"
	"github.com/kukymbr/tgnotifier/internal/sender"
	"github.com/kukymbr/tgnotifier/pkg/tgkit"
	"net/http"
)

// NewDefaultDependencyContainer builds new default DependencyContainer.
func NewDefaultDependencyContainer(configFilePath string) (DependencyContainer, error) {
	conf, err := config.NewConfig(config.FromDefaultFileDiscovery(configFilePath))

	if err != nil {
		return DependencyContainer{}, err
	}

	client := tgkit.NewClientWithOptions(
		tgkit.WithHTTPClient(&http.Client{
			Timeout: conf.Client().GetTimeout(),
		}),
		tgkit.WithRetrier(conf.GetRequestRetrier()),
	)

	proc := msgproc.NewProcessingChain(
		msgproc.NewTextNormalizer(),
		msgproc.NewReplacer(conf.Replaces()),
	)

	return DependencyContainer{
		Config:    conf,
		Client:    client,
		Sender:    sender.New(conf, client, proc),
		Processor: proc,
	}, nil
}

// DependencyContainer is a struct containing components of the tgnotifier.
type DependencyContainer struct {
	Config    *config.Config
	Client    *tgkit.Client
	Sender    *sender.Sender
	Processor msgproc.MessageProcessor
}
