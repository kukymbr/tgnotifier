package types

import (
	"fmt"
	"time"
)

// MustParseKitchenTime parses input string into the KitchenTime value and panics on the error.
func MustParseKitchenTime(input string) KitchenTime {
	t, err := ParseKitchenTime(input)
	if err != nil {
		panic(err)
	}

	return t
}

// ParseKitchenTime parses input string into the KitchenTime value.
func ParseKitchenTime(input string) (KitchenTime, error) {
	var format string

	switch len(input) {
	case 4:
		input = "0" + input
		format = "15:04"
	case 5:
		format = "15:04"
	case 8:
		format = "15:04:05"
	}

	t, err := time.Parse(format, input)
	if err != nil {
		return KitchenTime{}, fmt.Errorf("failed to parse kitchen time: %w", err)
	}

	return KitchenTime{
		hour:   t.Hour(),
		minute: t.Minute(),
		second: t.Second(),
	}, nil
}

// KitchenTime is a hh:mm:ss time.
type KitchenTime struct {
	hour   int
	minute int
	second int
}

// Time converts KitchenTime into the time.Time value.
func (k KitchenTime) Time() time.Time {
	now := time.Now()

	return time.Date(now.Year(), now.Month(), now.Day(), k.hour, k.minute, k.second, 0, now.Location())
}

func (k KitchenTime) String() string {
	return k.Time().Format("15:04:05")
}
