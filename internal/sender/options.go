package sender

// MessageOptions is a sending message options.
type MessageOptions struct {
	Text      string
	ParseMode string

	DisableNotification bool
	ProtectContent      bool
}
