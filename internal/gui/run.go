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

func Run(conf *config.Config, sender Sender) {
	ui.Init()

	ctrl := newController(conf, sender)

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
