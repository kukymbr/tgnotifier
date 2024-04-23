package tgkit

import (
	"strconv"
	"strings"
)

// ChatID is a telegram chat ID/username.
type ChatID string

func (id ChatID) String() string {
	str := string(id)

	if str == "" {
		return ""
	}

	if _, err := strconv.ParseInt(string(id), 10, 64); err == nil {
		return str
	}

	if !strings.HasPrefix(str, "@") {
		str = "@" + str
	}

	return str
}
