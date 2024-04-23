package sender

import "github.com/kukymbr/tgnotifier/internal/types"

// MessageOptions is a sending message options.
type MessageOptions struct {
	Text      string
	ParseMode types.ParseMode

	DisableNotification bool
	ProtectContent      bool
}
