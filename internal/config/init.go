package config

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kukymbr/tgnotifier/internal/types"
	"github.com/kukymbr/tgnotifier/pkg/tgkit"
	"gopkg.in/yaml.v3"
)

const (
	defaultGRPCHost = "127.0.0.1"
	defaultGRPCPort = 80

	defaultHTTPHost = "127.0.0.1"
	defaultHTTPPort = 8080

	defaultClientTimeout = 30 * time.Second

	retrierNoop        = "noop"
	retrierLinear      = "linear"
	retrierProgressive = "progressive"
)

// New reads config from the file if existing file given,
// and from the env if values are presented.
func New(readerFactory ...SourceReaderFactory) (*Config, error) {
	conf := &Config{
		bots: types.NewNamed[types.BotName, tgkit.Bot]("", func(bot tgkit.Bot) types.BotName {
			return types.BotToName(bot)
		}),
		chats: types.NewNamed[types.ChatName, tgkit.ChatID]("", func(chatID tgkit.ChatID) types.ChatName {
			return types.ChatIDToName(chatID)
		}),
	}

	var (
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

		if err := readConfigFromReader(conf, reader); err != nil {
			return nil, err
		}
	}

	if err := setupConfigValues(conf); err != nil {
		return nil, err
	}

	return conf, nil
}

func readConfigFromReader(conf *Config, inp io.Reader) error {
	data, err := io.ReadAll(inp)
	if err != nil {
		return fmt.Errorf("read config data from reader: %w", err)
	}

	var dto configDTO

	if err := yaml.Unmarshal(data, &dto); err != nil {
		return fmt.Errorf("decode config data: %w", err)
	}

	return setValuesFromDTO(conf, dto)
}

func setValuesFromDTO(conf *Config, dto configDTO) error {
	var err error

	for botName, identity := range dto.Bots {
		bot, err := tgkit.NewBot(identity)
		if err != nil {
			return fmt.Errorf("read bot from config: %w", err)
		}

		conf.Bots().Add(botName, bot)
	}

	for chatName, chatIDStr := range dto.Chats {
		chatID := tgkit.ChatID(chatIDStr)
		if chatID == "" {
			return fmt.Errorf("empty chat ID for %s in config", chatName)
		}

		conf.chats.Add(chatName, chatID)
	}

	conf.Bots().SetDefaultName(dto.DefaultBot)
	conf.Chats().SetDefaultName(dto.DefaultChat)

	conf.silenceSchedule, err = parseTimeSchedule(dto.SilenceSchedule)
	if err != nil {
		return fmt.Errorf("parse silence schedule from config: %w", err)
	}

	conf.retrier, err = newRetrier(dto.Retrier)
	if err != nil {
		return fmt.Errorf("failed to initialize request retrier from config: %w", err)
	}

	conf.client = ClientConfig{
		timeout: dto.Client.Timeout,
	}

	conf.replaces = dto.Replaces
	conf.grpc = ServerConfig{
		host: dto.GRPC.Host,
		port: dto.GRPC.Port,
	}
	conf.http = ServerConfig{
		host: dto.HTTP.Host,
		port: dto.HTTP.Port,
	}

	return nil
}

func setupConfigValues(conf *Config) error {
	if err := readDefaultsFromEnv(conf); err != nil {
		return err
	}

	if conf.Bots().Len() == 0 {
		return fmt.Errorf("no bots registered in config file or env")
	}

	setServerDefaults(&conf.grpc, defaultGRPCHost, defaultGRPCPort)
	setServerDefaults(&conf.http, defaultHTTPHost, defaultHTTPPort)

	if conf.retrier == nil {
		conf.retrier = tgkit.NewProgressiveRetrier(3, 500*time.Millisecond, 2)
	}

	if conf.client.timeout <= 0 {
		conf.client.timeout = defaultClientTimeout
	}

	return nil
}

func readDefaultsFromEnv(conf *Config) error {
	if identity := os.Getenv(EnvDefaultBot); identity != "" {
		bot, err := tgkit.NewBot(identity)
		if err != nil {
			return fmt.Errorf("read bot from %s: %w", EnvDefaultBot, err)
		}

		name := types.BotToName(bot)

		conf.Bots().Add(name, bot)
		conf.Bots().SetDefaultName(name)
	}

	if chatIDStr := os.Getenv(EnvDefaultChat); chatIDStr != "" {
		chatID := tgkit.ChatID(chatIDStr)
		name := types.ChatIDToName(chatID)

		conf.Chats().Add(name, chatID)
		conf.Chats().SetDefaultName(name)
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

func setServerDefaults(conf *ServerConfig, defaultHost string, defaultPort int) {
	if conf.host == "" {
		conf.host = defaultHost
	}

	if conf.port == 0 {
		conf.port = defaultPort
	}
}

func newRetrier(dto retrierConfigDTO) (tgkit.RequestRetrier, error) {
	switch dto.Type {
	case retrierNoop, "":
		return tgkit.NewNoopRetrier(), nil
	case retrierLinear:
		return tgkit.NewLinearRetrier(dto.Attempts, dto.Delay), nil
	case retrierProgressive:
		return tgkit.NewProgressiveRetrier(dto.Attempts, dto.Delay, dto.Multiplier), nil
	}

	return nil, fmt.Errorf("unknown retrier type: %s", dto.Type)
}
