package config

import (
	"fmt"
	"io"
	"os"

	"github.com/kukymbr/tgnotifier/internal/types"
	"github.com/kukymbr/tgnotifier/pkg/tgkit"
	"gopkg.in/yaml.v3"
)

const (
	EnvDefaultBot  = "TGNOTIFIER_DEFAULT_BOT"
	EnvDefaultChat = "TGNOTIFIER_DEFAULT_CHAT"
)

// NewConfig reads config from the file if existing file given,
// and from the env if values are presented.
func NewConfig(path string, getEnv func(string) string) (*Config, error) {
	var (
		conf = &Config{}
		err  error
	)

	if path != "" {
		conf, err = NewConfigFromFile(path)
		if err != nil {
			return nil, err
		}
	}

	if err := readDefaultsFromEnv(conf, getEnv); err != nil {
		return nil, err
	}

	if len(conf.bots) == 0 {
		return nil, fmt.Errorf("no bots registered in config file or env")
	}

	if len(conf.chats) == 0 {
		return nil, fmt.Errorf("no chats registered in config file or env")
	}

	return conf, nil
}

// NewConfigFromFile reads Config from the file.
func NewConfigFromFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("read config file %s: %w", path, err)
	}

	defer func() {
		_ = f.Close()
	}()

	return NewConfigFromReader(f)
}

// NewConfigFromReader reads Config from io.Reader.
func NewConfigFromReader(inp io.Reader) (*Config, error) {
	data, err := io.ReadAll(inp)
	if err != nil {
		return nil, fmt.Errorf("read config data from reader: %w", err)
	}

	var raw *config

	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("decode config data: %w", err)
	}

	if raw == nil || len(raw.Bots) == 0 {
		return nil, fmt.Errorf("invalid config: no bots given")
	}

	if len(raw.Chats) == 0 {
		return nil, fmt.Errorf("invalid config: no chats given")
	}

	conf := &Config{}
	conf.init()

	for botName, identity := range raw.Bots {
		bot, err := tgkit.NewBot(identity)
		if err != nil {
			return nil, fmt.Errorf("read bot from config: %w", err)
		}

		conf.bots[botName] = bot
	}

	for chatName, chatIDStr := range raw.Chats {
		chatID := tgkit.ChatID(chatIDStr)
		if chatID == "" {
			return nil, fmt.Errorf("empty chat ID for %s in config", chatName)
		}

		conf.chats[chatName] = chatID
	}

	conf.defaultBotName = raw.DefaultBot
	conf.defaultChatName = raw.DefaultChat

	return conf, nil
}

func readDefaultsFromEnv(conf *Config, getEnv func(string) string) error {
	conf.init()

	if identity := getEnv(EnvDefaultBot); identity != "" {
		bot, err := tgkit.NewBot(identity)
		if err != nil {
			return fmt.Errorf("read bot from %s: %w", EnvDefaultBot, err)
		}

		conf.bots[types.DefaultBotName] = bot
		conf.defaultBotName = types.DefaultBotName
	}

	if chatIDStr := getEnv(EnvDefaultChat); chatIDStr != "" {
		conf.chats[types.DefaultChatName] = tgkit.ChatID(chatIDStr)
		conf.defaultChatName = types.DefaultChatName
	}

	return nil
}

// Config is a tgnotifier configuration.
type Config struct {
	bots  BotsIndex
	chats ChatsIndex

	defaultBotName  types.BotName
	defaultChatName types.ChatName
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

type config struct {
	Bots  map[types.BotName]string  `json:"bots" yaml:"bots"`
	Chats map[types.ChatName]string `json:"chats" yaml:"chats"`

	DefaultBot  types.BotName  `json:"default_bot" yaml:"default_bot"`
	DefaultChat types.ChatName `json:"default_chat" yaml:"default_chat"`
}
