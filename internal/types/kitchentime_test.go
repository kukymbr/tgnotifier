package types_test

import (
	"testing"

	"github.com/kukymbr/tgnotifier/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestParseKitchenTime(t *testing.T) {
	tests := []struct {
		Input       string
		Expected    string
		ErrExpected bool
	}{
		{
			Input:    "11:00",
			Expected: "11:00:00",
		},
		{
			Input:    "12:32",
			Expected: "12:32:00",
		},
		{
			Input:    "12:32:33",
			Expected: "12:32:33",
		},
		{
			Input:       "2024-01-01 00:00:00",
			ErrExpected: true,
		},
		{
			Input:       "0:00:00",
			ErrExpected: true,
		},
		{
			Input:       "invalid",
			ErrExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Input, func(t *testing.T) {
			k, err := types.ParseKitchenTime(test.Input)

			if !test.ErrExpected {
				assert.NoError(t, err)
				assert.Equal(t, test.Expected, k.String())
			} else {
				assert.Error(t, err)
			}
		})
	}
}
