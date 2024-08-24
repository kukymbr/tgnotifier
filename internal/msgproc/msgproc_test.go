package msgproc_test

import (
	"fmt"
	"testing"

	"github.com/kukymbr/tgnotifier/internal/msgproc"
	"github.com/kukymbr/tgnotifier/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestMessageProcessor_Process(t *testing.T) {
	tests := []struct {
		Input     string
		Processor msgproc.MessageProcessor
		Expected  string
	}{
		{
			Input:     " Hello, @testUser1! ",
			Processor: msgproc.NewTextNormalizer(),
			Expected:  "Hello, @testUser1!",
		},
		{
			Input:     " Hello, @testUser2! ",
			Processor: msgproc.NewReplacer(types.KeyVal{"@testUser2": "@id04041"}),
			Expected:  ` Hello, @id04041! `,
		},
		{
			Input: " Hello, @testUser3! ",
			Processor: msgproc.NewProcessingChain(
				msgproc.NewTextNormalizer(),
				msgproc.NewReplacer(types.KeyVal{"@testUser3": "@id04042"}),
			),
			Expected: `Hello, @id04042!`,
		},
	}

	for i, test := range tests {
		test := test

		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			msg := test.Processor.Process(test.Input)

			assert.Equal(t, test.Expected, msg)
		})
	}
}
