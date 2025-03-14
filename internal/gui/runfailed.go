//go:build gui

package gui

import (
	"fmt"
	"github.com/roblillack/spot"
	"github.com/roblillack/spot/ui"
	"log"
)

const (
	failedWindowWidth  = 400
	failedWindowHeight = 400
)

func RunFailed(err error) {
	ui.Init()

	errStr := fmt.Sprintf("Failed to start UI: %s", err.Error())

	log.Printf(errStr)

	spot.MountFn(func(ctx *spot.RenderContext) spot.Component {
		return &ui.Window{
			Title:  "tgnotifier UI",
			Width:  failedWindowWidth,
			Height: failedWindowHeight,
			Children: []spot.Component{
				&ui.TextView{
					X:      margin,
					Y:      margin,
					Width:  failedWindowWidth - margin*2,
					Height: failedWindowHeight - margin*2,
					Text:   errStr,
				},
			},
		}
	})

	ui.Run()
}
