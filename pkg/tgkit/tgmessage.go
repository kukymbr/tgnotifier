package tgkit

// TgMessageRequest is a message request body.
// See https://core.telegram.org/bots/api#sendmessage
// TODO: add entities, link_preview_options, reply_parameters and reply_markup fields.
type TgMessageRequest struct {
	ChatID ChatID `json:"chat_id"`
	Text   string `json:"text"`

	BusinessConnectionId string `json:"business_connection_id,omitempty"`
	MessageThreadId      int    `json:"message_thread_id,omitempty"`
	ParseMode            string `json:"parse_mode,omitempty"`

	DisableNotification bool `json:"disable_notification,omitempty"`
	ProtectContent      bool `json:"protect_content,omitempty"`
}

// TgMessageResponse is a Telegram response with a single message result.
type TgMessageResponse struct {
	Ok     bool      `json:"ok"`
	Result TgMessage `json:"result"`
}

// TgMessage is a Telegram message model.
// See https://core.telegram.org/bots/api#message
// TODO: add a lot of fields :)
type TgMessage struct {
	MessageID       int    `json:"message_id"`
	MessageThreadId int    `json:"message_thread_id,omitempty"`
	Date            uint64 `json:"date"`

	From TgUser `json:"from"`
	Chat TgChat `json:"chat"`
}
