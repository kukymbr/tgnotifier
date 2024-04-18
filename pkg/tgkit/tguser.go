package tgkit

// TgUserResponse is a response from the Telegram with a single TgUser object inside.
type TgUserResponse struct {
	Ok     bool
	Result TgUser
}

// TgUser is a model of the Telegram User.
// See https://core.telegram.org/bots/api#user
type TgUser struct {
	ID int64 `json:"id"`

	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`

	IsBot     bool `json:"is_bot"`
	IsPremium bool `json:"is_premium,omitempty"`

	AddedToAttachmentMenu   bool `json:"added_to_attachment_menu,omitempty"`
	CanJoinGroups           bool `json:"can_join_groups,omitempty"`
	CanReadAllGroupMessages bool `json:"can_read_all_group_messages,omitempty"`
	SupportsInlineQueries   bool `json:"supports_inline_queries,omitempty"`
	CanConnectToBusiness    bool `json:"can_connect_to_business,omitempty"`
}
