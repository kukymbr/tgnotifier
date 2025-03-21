package config

import (
	"net"
	"strconv"
	"time"

	"github.com/kukymbr/tgnotifier/internal/types"
	"github.com/kukymbr/tgnotifier/pkg/tgkit"
)

const (
	EnvPrefix      = "TGNOTIFIER_"
	EnvConfigPath  = EnvPrefix + "CONFIG_PATH"
	EnvDefaultBot  = EnvPrefix + "DEFAULT_BOT"
	EnvDefaultChat = EnvPrefix + "DEFAULT_CHAT"
)

// Config is a tgnotifier configuration.
type Config struct {
	bots  *types.Named[types.BotName, tgkit.Bot]
	chats *types.Named[types.ChatName, tgkit.ChatID]

	silenceSchedule *types.TimeSchedule

	replaces map[string]string

	grpc ServerConfig
	http ServerConfig

	client  ClientConfig
	retrier tgkit.RequestRetrier
}

// Bots - returns registered bots index.
func (c *Config) Bots() *types.Named[types.BotName, tgkit.Bot] {
	return c.bots
}

// Chats returns registered chats index.
func (c *Config) Chats() *types.Named[types.ChatName, tgkit.ChatID] {
	return c.chats
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
