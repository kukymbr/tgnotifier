package config

import (
	"fmt"
	"github.com/kukymbr/tgnotifier/internal/types"
	"github.com/kukymbr/tgnotifier/pkg/tgkit"
	"net"
	"strconv"
	"time"
)

const (
	EnvDefaultBot  = "TGNOTIFIER_DEFAULT_BOT"
	EnvDefaultChat = "TGNOTIFIER_DEFAULT_CHAT"
)

// Config is a tgnotifier configuration.
type Config struct {
	bots  BotsIndex
	chats ChatsIndex

	defaultBotName  types.BotName
	defaultChatName types.ChatName

	silenceSchedule *types.TimeSchedule

	replaces map[string]string

	grpc ServerConfig
	http ServerConfig

	client  ClientConfig
	retrier tgkit.RequestRetrier
}

// Bots - returns registered bots index.
func (c *Config) Bots() BotsIndex {
	return c.bots
}

// Chats returns registered chats index.
func (c *Config) Chats() ChatsIndex {
	return c.chats
}

// GetDefaultBotName returns a default bot name if no bot name defined in arguments.
func (c *Config) GetDefaultBotName() types.BotName {
	return c.defaultBotName
}

// GetDefaultChatName returns a default chat name if no chat name defined in arguments.
func (c *Config) GetDefaultChatName() types.ChatName {
	return c.defaultChatName
}

// GetSilenceSchedule returns a schedule when all the messages should be sent without a sound.
func (c *Config) GetSilenceSchedule() *types.TimeSchedule {
	if c.silenceSchedule == nil {
		c.silenceSchedule = &types.TimeSchedule{}
	}

	return c.silenceSchedule
}

// Replaces returns substrings to be replaced in the messages.
func (c *Config) Replaces() map[string]string {
	if c.replaces == nil {
		c.replaces = make(map[string]string)
	}

	return c.replaces
}

// Client returns Telegram API client configuration.
func (c *Config) Client() ClientConfig {
	return c.client
}

// GRPC returns configuration of the gRPC server.
func (c *Config) GRPC() ServerConfig {
	return c.grpc
}

// HTTP returns configuration of the HTTP server.
func (c *Config) HTTP() ServerConfig {
	return c.http
}

// GetRequestRetrier returns tgkit.RequestRetrier instance.
func (c *Config) GetRequestRetrier() tgkit.RequestRetrier {
	return c.retrier
}

func (c *Config) init() {
	if c.bots == nil {
		c.bots = make(BotsIndex)
	}

	if c.chats == nil {
		c.chats = make(ChatsIndex)
	}
}

// BotsIndex is an index of the registered bots.
type BotsIndex map[types.BotName]*tgkit.Bot

// GetBot finds registered bot by the name.
func (b BotsIndex) GetBot(name types.BotName) (*tgkit.Bot, error) {
	bot, ok := b[name]
	if !ok {
		return nil, fmt.Errorf("bot %s is not registered", name)
	}

	return bot, nil
}

// ChatsIndex is an index of the registered chats.
type ChatsIndex map[types.ChatName]tgkit.ChatID

// GetChatID finds registered chat ID by the name.
func (b ChatsIndex) GetChatID(name types.ChatName) (tgkit.ChatID, error) {
	chatID, ok := b[name]
	if !ok {
		return "", fmt.Errorf("chat %s is not registered", name)
	}

	return chatID, nil
}

type ServerConfig struct {
	host string
	port int
}

func (c ServerConfig) GetAddress() string {
	return net.JoinHostPort(c.GetHost(), strconv.Itoa(c.GetPort()))
}

func (c ServerConfig) GetHost() string {
	return c.host
}

func (c ServerConfig) GetPort() int {
	return c.port
}

type ClientConfig struct {
	timeout time.Duration
}

func (c ClientConfig) GetTimeout() time.Duration {
	return c.timeout
}
