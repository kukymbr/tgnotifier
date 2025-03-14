package tgkit

import (
	"fmt"
	"regexp"
	"strconv"
)

var rxBotIdentity = regexp.MustCompile(`(?i)^(?:bot)*([0-9]+):([a-zA-Z0-9\-_]+)$`)

// BotID is a Bot ID number.
type BotID uint64

func (id BotID) String() string {
	return fmt.Sprintf("%d", id)
}

// BotToken is a Bot token string.
type BotToken string

func (t BotToken) String() string {
	return string(t)
}

// BotIdentity is a Bot identity string.
type BotIdentity string

func (t BotIdentity) String() string {
	return string(t)
}

// NewBot receives a string of the bot identity ("<botID>:<botToken>")
// and returns the new Bot instance.
// Returns an error if given identity is not in valid format.
//
// Examples of the identity strings:
// - `bot12345:botToken1`
// - `54321:botToken2`
func NewBot(identity string) (Bot, error) {
	matches := rxBotIdentity.FindAllStringSubmatch(identity, -1)

	if len(matches) == 0 {
		return Bot{}, fmt.Errorf("invalid bot identity string format")
	}

	id, _ := strconv.ParseUint(matches[0][1], 10, 64)

	return Bot{
		id:    BotID(id),
		token: BotToken(matches[0][2]),
	}, nil
}

// MustNewBot creates new Bot and panics on failure.
func MustNewBot(identity string) Bot {
	bot, err := NewBot(identity)
	if err != nil {
		panic(err)
	}

	return bot
}

// Bot is a telegram bot model.
type Bot struct {
	id    BotID
	token BotToken
}

// GetID returns a Bot ID number.
func (b Bot) GetID() BotID {
	return b.id
}

// GetToken returns a Bot token string.
func (b Bot) GetToken() BotToken {
	return b.token
}

// GetIdentity returns a Bot identity string to pass it to the API.
func (b Bot) GetIdentity() BotIdentity {
	if b.GetID() == BotID(0) {
		return ""
	}

	return BotIdentity("bot" + b.GetID().String() + ":" + b.GetToken().String())
}

// String returns a Bot identity with a masked token for the debug purposes.
func (b Bot) String() string {
	if b.GetID() == BotID(0) {
		return ""
	}

	token := b.GetToken().String()
	if token != "" {
		token = "*****"
	}

	return "bot" + b.GetID().String() + ":" + token
}
