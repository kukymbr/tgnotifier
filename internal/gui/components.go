//go:build gui

package gui

import (
	"github.com/roblillack/spot"
	"github.com/roblillack/spot/ui"
	"strings"
)

var (
	currentPosition int
	windowHeight    int
)

const (
	windowWidth = columnWidth*2 + margin

	columnWidth  = 320
	margin       = 20
	elementWidth = columnWidth - margin*2

	labelHeight    = 15
	dropdownHeight = 20
	textareaHeight = 150

	buttonHeight = 20
	buttonWidth  = (columnWidth - margin*2) / 2
)

func createComponents(ctx *spot.RenderContext, ctrl *controller) (components []spot.Component, windowH int) {
	ctrl.responseContent, ctrl.setResponseContent = spot.UseState(ctx, "")

	return []spot.Component{
		label("Bot:"),
		dropdown(
			ctrl.getDefaultBotIndex(),
			func(index int) {
				ctrl.selectedBotIndex = index
			},
			ctrl.getBots()...,
		),
		button("Check bot...", func() {
			ctrl.getMe()
		}),
		button("Get updates...", func() {
			ctrl.getUpdates()
		}),

		label("Send to:"),
		dropdown(
			ctrl.getDefaultChatIndex(),
			func(index int) {
				ctrl.selectedChatIndex = index
			},
			ctrl.getChats()...,
		),

		label("Message"),
		textarea(func(s string) {
			ctrl.messageContent = s
		}),

		button("Send it now!", func() {
			ctrl.sendMessage()
		}),

		&ui.TextView{
			X:      columnWidth + margin/2,
			Y:      margin,
			Width:  elementWidth,
			Height: windowHeight - margin*2,
			Text:   ctrl.responseContent,
		},
	}, windowHeight
}

func registerComponentSize(height int) (y int) {
	currentPosition += margin

	y = currentPosition

	currentPosition += height

	windowHeight = currentPosition + margin

	return y
}

func label(text string) *ui.Label {
	return &ui.Label{
		Value:  text,
		X:      margin,
		Y:      registerComponentSize(labelHeight),
		Width:  elementWidth,
		Height: labelHeight,
		Align:  ui.LabelAlignmentLeft,
	}
}

func textarea(onChange func(string), value ...string) *ui.TextEditor {
	return &ui.TextEditor{
		X:        margin,
		Y:        registerComponentSize(textareaHeight),
		Width:    elementWidth,
		Height:   textareaHeight,
		Text:     strings.Join(value, ""),
		OnChange: onChange,
	}
}

func dropdown(index int, onChange func(int), items ...string) *ui.Dropdown {
	return &ui.Dropdown{
		X:                    margin,
		Y:                    registerComponentSize(dropdownHeight),
		Width:                elementWidth,
		Height:               dropdownHeight,
		Items:                items,
		SelectedIndex:        index,
		OnSelectionDidChange: onChange,
	}
}

func button(title string, onClick func()) *ui.Button {
	return &ui.Button{
		X:       margin,
		Y:       registerComponentSize(buttonHeight),
		Width:   buttonWidth,
		Height:  buttonHeight,
		Title:   title,
		OnClick: onClick,
	}
}
