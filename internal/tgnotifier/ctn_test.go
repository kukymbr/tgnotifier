package tgnotifier_test

import (
	"github.com/kukymbr/tgnotifier/internal/tgnotifier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewDefaultDependencyContainer(t *testing.T) {
	t.Setenv("TGNOTIFIER_DEFAULT_BOT", "bot1:test")
	t.Setenv("TGNOTIFIER_DEFAULT_CHAT", "-12345")

	var (
		ctn tgnotifier.DependencyContainer
		err error
	)

	assert.NotPanics(t, func() {
		ctn, err = tgnotifier.NewDefaultDependencyContainer("")
	})

	require.NoError(t, err)

	assert.NotNil(t, ctn.Config)
	assert.NotNil(t, ctn.Client)
	assert.NotNil(t, ctn.Sender)
	assert.NotNil(t, ctn.Processor)
}
