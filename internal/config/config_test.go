package config_test

import (
	"testing"
	"time"

	"github.com/kukymbr/tgnotifier/internal/config"
	"github.com/kukymbr/tgnotifier/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertBot(t *testing.T, conf *config.Config, name types.BotName, expected string) {
	bot, err := conf.Bots().GetBot(name)

	require.NoError(t, err)
	require.NotNil(t, bot)

	assert.Equal(t, expected, bot.GetIdentity().String())
}

func assertChat(t *testing.T, conf *config.Config, name types.ChatName, expected string) {
	chat, err := conf.Chats().GetChatID(name)

	require.NoError(t, err)
	require.NotEmpty(t, chat)

	assert.Equal(t, expected, chat.String())
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		Name       string
		Env        map[string]string
		ConfigFile string
		Assert     func(t *testing.T, conf *config.Config, err error)
	}{
		{
			Name: "With no env no file",
			Assert: func(t *testing.T, conf *config.Config, err error) {
				assert.Error(t, err)
				assert.Nil(t, conf)
			},
		},
		{
			Name: "With no chats registered",
			Env: map[string]string{
				"TGNOTIFIER_DEFAULT_BOT": "bot1:test1",
			},
			Assert: func(t *testing.T, conf *config.Config, err error) {
				assert.Error(t, err)
				assert.Nil(t, conf)
			},
		},
		{
			Name: "With env only",
			Env: map[string]string{
				"TGNOTIFIER_DEFAULT_BOT":  "bot2:test2",
				"TGNOTIFIER_DEFAULT_CHAT": "@testChat2",
			},
			Assert: func(t *testing.T, conf *config.Config, err error) {
				require.NoError(t, err)
				require.NotNil(t, conf)

				assert.Len(t, conf.Bots(), 1)
				assert.Len(t, conf.Chats(), 1)

				assert.Equal(t, types.DefaultBotName, conf.GetDefaultBotName())
				assert.Equal(t, types.DefaultChatName, conf.GetDefaultChatName())

				assertBot(t, conf, types.DefaultBotName, "bot2:test2")
				assertChat(t, conf, types.DefaultChatName, "@testChat2")

				assert.NotNil(t, conf.GetSilenceSchedule())

				replaces := conf.Replaces()

				assert.NotNil(t, replaces)
				assert.Len(t, replaces, 0)

				assert.Equal(t, conf.GRPC().GetAddress(), "127.0.0.1:80")
				assert.Equal(t, conf.HTTP().GetAddress(), "127.0.0.1:8080")

				assert.NotNil(t, conf.GetRequestRetrier())
				assert.Equal(t, 30*time.Second, conf.Client().GetTimeout())
			},
		},
		{
			Name:       "With config file only",
			ConfigFile: "./testdata/.tgnotifier.1.yml",
			Assert: func(t *testing.T, conf *config.Config, err error) {
				require.NoError(t, err)
				require.NotNil(t, conf)

				assert.Len(t, conf.Bots(), 2)
				assert.Len(t, conf.Chats(), 2)

				assert.Equal(t, types.BotName("first_bot"), conf.GetDefaultBotName())
				assert.Equal(t, types.ChatName("main_chat"), conf.GetDefaultChatName())

				assertBot(t, conf, "first_bot", "bot12345:FIRST_BOT_TOKEN")
				assertBot(t, conf, "second_bot", "bot54321:SECOND_BOT_TOKEN")

				assertChat(t, conf, "main_chat", "-12345")
				assertChat(t, conf, "secondary_chat", "@my_test_channel")

				interval := conf.GetSilenceSchedule()
				now := time.Now()

				assert.True(
					t,
					interval.Has(time.Date(
						now.Year(), now.Month(), now.Day(),
						23, 30, 0, 0,
						now.Location(),
					)),
				)

				assert.Len(t, conf.Replaces(), 1)

				assert.Equal(t, conf.GRPC().GetAddress(), "127.0.0.1:4080")
				assert.Equal(t, conf.HTTP().GetAddress(), "127.0.0.1:4081")

				assert.NotNil(t, conf.GetRequestRetrier())
				assert.Equal(t, 31*time.Second, conf.Client().GetTimeout())
			},
		},
		{
			Name:       "With config file and env",
			ConfigFile: "./testdata/.tgnotifier.1.yml",
			Env: map[string]string{
				"TGNOTIFIER_DEFAULT_BOT":  "bot3:test3",
				"TGNOTIFIER_DEFAULT_CHAT": "@testChat3",
			},
			Assert: func(t *testing.T, conf *config.Config, err error) {
				require.NoError(t, err)
				require.NotNil(t, conf)

				assert.Len(t, conf.Bots(), 3)
				assert.Len(t, conf.Chats(), 3)

				assert.Equal(t, types.DefaultBotName, conf.GetDefaultBotName())
				assert.Equal(t, types.DefaultChatName, conf.GetDefaultChatName())

				assertBot(t, conf, types.DefaultBotName, "bot3:test3")
				assertBot(t, conf, "first_bot", "bot12345:FIRST_BOT_TOKEN")
				assertBot(t, conf, "second_bot", "bot54321:SECOND_BOT_TOKEN")

				assertChat(t, conf, "main_chat", "-12345")
				assertChat(t, conf, "secondary_chat", "@my_test_channel")
				assertChat(t, conf, types.DefaultChatName, "@testChat3")
			},
		},
		{
			Name:       "With unknown file",
			ConfigFile: "./testdata/unknown.yml",
			Assert: func(t *testing.T, conf *config.Config, err error) {
				assert.Error(t, err)
				assert.Nil(t, conf)
			},
		},
		{
			Name:       "With invalid file (unmarshal)",
			ConfigFile: "./testdata/.tgnotifier.invalid.unmarshal.yml",
			Assert: func(t *testing.T, conf *config.Config, err error) {
				assert.Error(t, err)
				assert.Nil(t, conf)
			},
		},
		{
			Name:       "With invalid file (bot)",
			ConfigFile: "./testdata/.tgnotifier.invalid.bot.yml",
			Assert: func(t *testing.T, conf *config.Config, err error) {
				assert.Error(t, err)
				assert.Nil(t, conf)
			},
		},
		{
			Name:       "With invalid file (chat)",
			ConfigFile: "./testdata/.tgnotifier.invalid.chat.yml",
			Assert: func(t *testing.T, conf *config.Config, err error) {
				assert.Error(t, err)
				assert.Nil(t, conf)
			},
		},
		{
			Name:       "With invalid file (empty)",
			ConfigFile: "./testdata/.tgnotifier.invalid.empty.yml",
			Assert: func(t *testing.T, conf *config.Config, err error) {
				assert.Error(t, err)
				assert.Nil(t, conf)
			},
		},
		{
			Name:       "With invalid file (no chats)",
			ConfigFile: "./testdata/.tgnotifier.invalid.nochats.yml",
			Assert: func(t *testing.T, conf *config.Config, err error) {
				assert.Error(t, err)
				assert.Nil(t, conf)
			},
		},
		{
			Name:       "With invalid file (broken schedule)",
			ConfigFile: "./testdata/.tgnotifier.invalid.schedule.yml",
			Assert: func(t *testing.T, conf *config.Config, err error) {
				assert.Error(t, err)
				assert.Nil(t, conf)
			},
		},
		{
			Name:       "With invalid file (unknown retrier)",
			ConfigFile: "./testdata/.tgnotifier.invalid.retrier.yml",
			Assert: func(t *testing.T, conf *config.Config, err error) {
				assert.Error(t, err)
				assert.Nil(t, conf)
			},
		},
		{
			Name: "With invalid env",
			Env:  map[string]string{"TGNOTIFIER_DEFAULT_BOT": "invalid identity"},
			Assert: func(t *testing.T, conf *config.Config, err error) {
				assert.Error(t, err)
				assert.Nil(t, conf)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			for key, val := range test.Env {
				t.Setenv(key, val)
			}

			var (
				conf *config.Config
				err  error
			)

			conf, err = config.NewConfig(config.FromFile(test.ConfigFile))

			test.Assert(t, conf, err)
		})
	}
}

func TestConfig_BotChatGetters(t *testing.T) {
	t.Setenv("TGNOTIFIER_DEFAULT_BOT", "bot123:test")
	t.Setenv("TGNOTIFIER_DEFAULT_CHAT", "@testChat")

	conf, err := config.NewConfig()

	require.NoError(t, err)
	require.NotNil(t, conf)

	t.Run("get unknown bot", func(t *testing.T) {
		bot, err := conf.Bots().GetBot("unknown")

		assert.Nil(t, bot)
		assert.Error(t, err)
	})

	t.Run("get unknown chat", func(t *testing.T) {
		chat, err := conf.Chats().GetChatID("unknown")

		assert.Empty(t, chat)
		assert.Error(t, err)
	})
}
