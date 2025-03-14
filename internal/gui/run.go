//go:build gui

package gui

import (
	"context"
	"github.com/kukymbr/tgnotifier/internal/config"
	"github.com/kukymbr/tgnotifier/internal/types"
	"github.com/kukymbr/tgnotifier/pkg/tgkit"
	"github.com/roblillack/spot"
	"github.com/roblillack/spot/ui"
)

type Sender interface {
	Send(ctx context.Context, options types.SendOptions) (*tgkit.TgMessage, error)
}

type Client interface {
	GetMe(bot tgkit.Bot) (*tgkit.TgUser, error)
	GetUpdates(bot tgkit.Bot) ([]tgkit.TgUpdate, error)
}

func Run(conf *config.Config, sender Sender, tgClient Client) {
	ui.Init()

	ctrl := newController(conf, sender, tgClient)

	spot.MountFn(func(ctx *spot.RenderContext) spot.Component {
		components, windowH := createComponents(ctx, ctrl)

		return &ui.Window{
			Title:    "tgnotifier UI",
			Width:    windowWidth,
			Height:   windowH,
			Children: components,
		}
	})

	ui.Run()
}
