package config

import (
	"fmt"
	"github.com/kukymbr/tgnotifier/internal/types"
	"github.com/kukymbr/tgnotifier/pkg/tgkit"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

const (
	defaultGRPCPort = 80
)

// NewConfig reads config from the file if existing file given,
// and from the env if values are presented.
func NewConfig(readerFactory ...SourceReaderFactory) (*Config, error) {
	var (
		conf   = &Config{}
		reader io.ReadCloser
		err    error
	)

	if len(readerFactory) > 0 {
		reader, err = readerFactory[0]()
		if err != nil {
			return nil, fmt.Errorf("failed to create config reader: %w", err)
		}
	}

	if reader != nil {
		defer func() {
			_ = reader.Close()
		}()

		conf, err = newConfigFromReader(reader)
		if err != nil {
			return nil, err
		}
	}

	if err := setupConfigValues(conf); err != nil {
		return nil, err
	}

	return conf, nil
}

// newConfigFromReader reads Config from io.Reader.
func newConfigFromReader(inp io.Reader) (*Config, error) {
	data, err := io.ReadAll(inp)
	if err != nil {
		return nil, fmt.Errorf("read config data from reader: %w", err)
	}

	var raw *configDTO

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

	conf.silenceSchedule, err = parseTimeSchedule(raw.SilenceSchedule)
	if err != nil {
		return nil, fmt.Errorf("parse silence schedule from config: %w", err)
	}

	conf.replaces = raw.Replaces
	conf.grpc = GRPCServerConfig{
		port: raw.GRPC.Port,
	}

	return conf, nil
}

func setupConfigValues(conf *Config) error {
	if err := readDefaultsFromEnv(conf); err != nil {
		return err
	}

	if len(conf.bots) == 0 {
		return fmt.Errorf("no bots registered in config file or env")
	}

	if len(conf.chats) == 0 {
		return fmt.Errorf("no chats registered in config file or env")
	}

	if conf.grpc.port == 0 {
		conf.grpc.port = defaultGRPCPort
	}

	return nil
}

func readDefaultsFromEnv(conf *Config) error {
	conf.init()

	if identity := os.Getenv(EnvDefaultBot); identity != "" {
		bot, err := tgkit.NewBot(identity)
		if err != nil {
			return fmt.Errorf("read bot from %s: %w", EnvDefaultBot, err)
		}

		conf.bots[types.DefaultBotName] = bot
		conf.defaultBotName = types.DefaultBotName
	}

	if chatIDStr := os.Getenv(EnvDefaultChat); chatIDStr != "" {
		conf.chats[types.DefaultChatName] = tgkit.ChatID(chatIDStr)
		conf.defaultChatName = types.DefaultChatName
	}

	return nil
}

func parseTimeSchedule(raw []silenceScheduleItem) (*types.TimeSchedule, error) {
	schedule := &types.TimeSchedule{}

	if len(raw) == 0 {
		return schedule, nil
	}

	for _, item := range raw {
		from, err := types.ParseKitchenTime(item.From)
		if err != nil {
			return nil, fmt.Errorf("parse 'from' value '%s': %w", item.From, err)
		}

		to, err := types.ParseKitchenTime(item.To)
		if err != nil {
			return nil, fmt.Errorf("parse 'to' value '%s': %w", item.From, err)
		}

		schedule.AddInterval(types.TimeInterval{
			From: from,
			To:   to,
		})
	}

	return schedule, nil
}
