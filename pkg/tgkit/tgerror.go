package tgkit

import "fmt"

// TgErrorResponse is a Telegram response describing the error.
type TgErrorResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}

func (r TgErrorResponse) Error() string {
	return fmt.Sprintf("telegram error response (ok=%t, error_code=%d): %s", r.Ok, r.ErrorCode, r.Description)
}
