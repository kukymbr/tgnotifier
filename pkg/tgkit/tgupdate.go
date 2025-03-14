package tgkit

type TgUpdatesResponse struct {
	Ok     bool       `json:"ok"`
	Result []TgUpdate `json:"result"`
}

// TgUpdate is struct representing the Telegram's Update object.
// See https://core.telegram.org/bots/api#update
// TODO: add a lot of fields :)
type TgUpdate struct {
	UpdateID int       `json:"update_id"`
	Message  TgMessage `json:"message"`
}
