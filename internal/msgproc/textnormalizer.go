package msgproc

import "strings"

// NewTextNormalizer returns a MessageProcessor normalizing the text.
func NewTextNormalizer() MessageProcessor {
	return &textNormalizer{}
}

type textNormalizer struct{}

func (mp *textNormalizer) Process(msg string) string {
	msg = strings.TrimSpace(msg)

	return msg
}
