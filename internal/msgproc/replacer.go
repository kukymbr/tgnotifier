package msgproc

import (
	"strings"

	"github.com/kukymbr/tgnotifier/internal/types"
)

// NewReplacer returns a MessageProcessor replacing substrings in the message.
func NewReplacer(replaces types.KeyVal) MessageProcessor {
	return &replacer{
		replaces: replaces,
	}
}

type replacer struct {
	replaces types.KeyVal
}

func (mp *replacer) Process(msg string) string {
	if len(mp.replaces) == 0 {
		return msg
	}

	for oldS, newS := range mp.replaces {
		msg = strings.ReplaceAll(msg, oldS, newS)
	}

	return msg
}
