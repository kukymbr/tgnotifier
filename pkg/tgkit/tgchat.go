package tgkit

// TgChat is a model representing the Telegram API chat object.
// See: https://core.telegram.org/bots/api#chat
//
// id	Integer	Unique identifier for this chat. This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this identifier.
// type	String	Type of the chat, can be either “private”, “group”, “supergroup” or “channel”
// title	String	Optional. Title, for supergroups, channels and group chats
// username	String	Optional. Username, for private chats, supergroups and channels if available
// first_name	String	Optional. First name of the other party in a private chat
// last_name	String	Optional. Last name of the other party in a private chat
// is_forum	True	Optional. True, if the supergroup chat is a forum (has topics enabled)
type TgChat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	IsForum   bool   `json:"is_forum,omitempty"`
}
