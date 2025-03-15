package types_test

import (
	"github.com/kukymbr/tgnotifier/internal/types"
	"github.com/kukymbr/tgnotifier/pkg/tgkit"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNamed(t *testing.T) {
	constructor1 := func(t *testing.T) *types.Named[types.BotName, tgkit.Bot] {
		list := types.NewNamed[types.BotName, tgkit.Bot]("bot2", func(bot tgkit.Bot) types.BotName {
			return types.BotToName(bot)
		})

		list.Add("test1", tgkit.MustNewBot("bot1:test1"))
		list.Add("test1", tgkit.MustNewBot("bot1:test1"))
		list.Add("", tgkit.MustNewBot("bot2:test2"))

		return list
	}

	constructor2 := func(t *testing.T) *types.Named[types.BotName, tgkit.Bot] {
		list := types.NewNamed[types.BotName, tgkit.Bot]("")

		list.Add("test3", tgkit.MustNewBot("bot3:test3"))

		return list
	}

	tests := []struct {
		Name        string
		Constructor func(*testing.T) *types.Named[types.BotName, tgkit.Bot]
		Assert      func(t *testing.T, list *types.Named[types.BotName, tgkit.Bot])
	}{
		{
			Name:        "assert counts",
			Constructor: constructor1,
			Assert: func(t *testing.T, list *types.Named[types.BotName, tgkit.Bot]) {
				assert.Equal(t, 2, list.Len())
				assert.Len(t, list.GetNames(), 2)
			},
		},
		{
			Name:        "when existing name requested, expect found",
			Constructor: constructor1,
			Assert: func(t *testing.T, list *types.Named[types.BotName, tgkit.Bot]) {
				bot, err := list.FindByName("test1")

				assert.NoError(t, err)
				assert.Equal(t, "bot1:test1", bot.GetIdentity().String())
			},
		},
		{
			Name:        "when default bot requested, expect found",
			Constructor: constructor1,
			Assert: func(t *testing.T, list *types.Named[types.BotName, tgkit.Bot]) {
				bot, err := list.FindByName("")

				assert.NoError(t, err)
				assert.Equal(t, "bot2:test2", bot.GetIdentity().String())
			},
		},
		{
			Name:        "when existing index requested, expect found",
			Constructor: constructor1,
			Assert: func(t *testing.T, list *types.Named[types.BotName, tgkit.Bot]) {
				bot, err := list.FindByNameIndex(0)

				assert.NoError(t, err)
				assert.Equal(t, "bot2:test2", bot.GetIdentity().String())
			},
		},
		{
			Name:        "when non-existing name requested, expect error",
			Constructor: constructor1,
			Assert: func(t *testing.T, list *types.Named[types.BotName, tgkit.Bot]) {
				bot, err := list.FindByName("test2")

				assert.Error(t, err)
				assert.Empty(t, bot.GetIdentity().String())
			},
		},
		{
			Name:        "when non-existing index requested, expect error",
			Constructor: constructor1,
			Assert: func(t *testing.T, list *types.Named[types.BotName, tgkit.Bot]) {
				bot, err := list.FindByNameIndex(10)

				assert.Error(t, err)
				assert.Empty(t, bot.GetIdentity().String())
			},
		},
		{
			Name:        "when no default bot, expect error",
			Constructor: constructor2,
			Assert: func(t *testing.T, list *types.Named[types.BotName, tgkit.Bot]) {
				bot, err := list.FindByName("")

				assert.Error(t, err)
				assert.Empty(t, bot.GetIdentity().String())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			list := test.Constructor(t)

			test.Assert(t, list)
		})
	}
}
