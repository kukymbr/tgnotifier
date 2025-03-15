//go:build gui

package gui

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/kukymbr/tgnotifier/internal/config"
	"github.com/kukymbr/tgnotifier/internal/types"
)

func newController(conf *config.Config, sender Sender, tgClient Client) *controller {
	ctrl := &controller{
		conf:     conf,
		sender:   sender,
		tgClient: tgClient,
	}

	ctrl.selectedBotIndex = ctrl.getDefaultBotIndex()
	ctrl.selectedChatIndex = ctrl.getDefaultChatIndex()

	return ctrl
}

type controller struct {
	conf     *config.Config
	sender   Sender
	tgClient Client

	selectedBotIndex  int
	selectedChatIndex int
	messageContent    string

	responseContent    string
	setResponseContent func(string)
}

func (c *controller) sendMessage() {
	c.printWait()

	opt := types.SendOptions{}

	_ = opt.BotName.Set(c.getBots()[c.selectedBotIndex])
	_ = opt.ChatName.Set(c.getChats()[c.selectedChatIndex])

	opt.Message = types.MessageOptions{
		Text: c.messageContent,
	}

	if err := opt.Validate(); err != nil {
		c.printError(err)

		return
	}

	resp, err := c.sender.Send(context.Background(), opt)
	if err != nil {
		c.printError(err)

		return
	}

	c.printResponse(resp)
}

func (c *controller) getMe() {
	c.printWait()

	bot, err := c.conf.Bots().FindByNameIndex(c.selectedBotIndex)
	if err != nil {
		c.printError(err)

		return
	}

	resp, err := c.tgClient.GetMe(bot)
	if err != nil {
		c.printError(err)

		return
	}

	c.printResponse(resp)
}

func (c *controller) getUpdates() {
	c.printWait()

	bot, err := c.conf.Bots().FindByNameIndex(c.selectedBotIndex)
	if err != nil {
		c.printError(err)

		return
	}

	resp, err := c.tgClient.GetUpdates(bot)
	if err != nil {
		c.printError(err)

		return
	}

	c.printResponse(resp)
}

func (c *controller) printResponse(resp any) {
	data, err := jsoniter.MarshalIndent(resp, "", "  ")
	if err != nil {
		c.printError(err)

		return
	}

	c.setResponseContent(string(data))
}

func (c *controller) printWait() {
	c.setResponseContent("Please wait...")
}

func (c *controller) printError(err error) {
	c.setResponseContent(fmt.Sprintf("ERROR: %s", err.Error()))
}

func (c *controller) getBots() []string {
	names := make([]string, 0, c.conf.Bots().Len())

	for _, name := range c.conf.Bots().GetNames() {
		names = append(names, name.String())
	}

	return names
}

func (c *controller) getChats() []string {
	names := make([]string, 0, c.conf.Chats().Len())

	for _, name := range c.conf.Chats().GetNames() {
		names = append(names, name.String())
	}

	return names
}

func (c *controller) getDefaultBotIndex() int {
	dft := c.conf.Bots().GetDefaultName()
	if dft == "" {
		return 0
	}

	for i, name := range c.getBots() {
		if name == dft.String() {
			return i
		}
	}

	return 0
}

func (c *controller) getDefaultChatIndex() int {
	dft := c.conf.Chats().GetDefaultName()
	if dft == "" {
		return 0
	}

	for i, name := range c.getChats() {
		if name == dft.String() {
			return i
		}
	}

	return 0
}
