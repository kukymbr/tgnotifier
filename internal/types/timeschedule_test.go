package types_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kukymbr/tgnotifier/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestTimeSchedule(t *testing.T) {
	schedule := types.TimeSchedule{}

	schedule.AddInterval(types.TimeInterval{
		From: types.MustParseKitchenTime("09:00"),
		To:   types.MustParseKitchenTime("10:10"),
	})

	schedule.AddInterval(types.TimeInterval{
		From: types.MustParseKitchenTime("12:00"),
		To:   types.MustParseKitchenTime("13:10"),
	})

	schedule.AddInterval(types.TimeInterval{
		From: types.MustParseKitchenTime("22:00"),
		To:   types.MustParseKitchenTime("07:00"),
	})

	now := time.Now()

	tests := []struct {
		GetTime func() time.Time
		Has     bool
	}{
		{
			GetTime: func() time.Time {
				return time.Date(
					now.Year(), now.Month(), now.Day(), 9, 0, 0, 0,
					now.Location(),
				)
			},
			Has: true,
		},
		{
			GetTime: func() time.Time {
				return time.Date(
					now.Year(), now.Month(), now.Day(), 10, 10, 1, 0,
					now.Location(),
				)
			},
			Has: false,
		},
		{
			GetTime: func() time.Time {
				return time.Date(
					now.Year(), now.Month(), now.Day(), 12, 30, 0, 0,
					now.Location(),
				)
			},
			Has: true,
		},
		{
			GetTime: func() time.Time {
				return time.Date(
					now.Year(), now.Month(), now.Day(), 13, 30, 0, 0,
					now.Location(),
				)
			},
			Has: false,
		},
		{
			GetTime: func() time.Time {
				return time.Date(
					now.Year(), now.Month(), now.Day(), 15, 0, 0, 0,
					now.Location(),
				)
			},
			Has: false,
		},
		{
			GetTime: func() time.Time {
				return time.Date(
					now.Year(), now.Month(), now.Day(), 23, 0, 0, 0,
					now.Location(),
				)
			},
			Has: true,
		},
		{
			GetTime: func() time.Time {
				return time.Date(
					now.Year(), now.Month(), now.Day(), 5, 0, 0, 0,
					now.Location(),
				)
			},
			Has: true,
		},
		{
			GetTime: func() time.Time {
				return time.Date(
					now.Year(), now.Month(), now.Day(), 7, 1, 0, 0,
					now.Location(),
				)
			},
			Has: false,
		},
		{
			GetTime: func() time.Time {
				return time.Date(
					now.Year(), now.Month(), now.Day(), 21, 59, 0, 0,
					now.Location(),
				)
			},
			Has: false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			val := test.GetTime()

			assert.Equal(t, test.Has, schedule.Has(val))
		})
	}
}
