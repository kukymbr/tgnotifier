package tgkit

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	tgAPIHost           = "https://api.telegram.org/"
	tgMethodGetMe       = "getMe"
	tgMethodSendMessage = "sendMessage"
)

// NewDefaultClient creates a new default telegram Client instance.
func NewDefaultClient() *Client {
	return NewClient(&http.Client{
		Timeout: 30 * time.Second,
	})
}

// NewClient creates a new telegram Client instance.
func NewClient(httpClient HTTPClient) *Client {
	return &Client{httpClient: httpClient}
}

// Client is a tool to communicate with a Telegram API via the HTTPS.
type Client struct {
	httpClient HTTPClient
}

// Get sends a GET request to the Telegram API.
// The successful response is decoded into the target.
func (c *Client) Get(bot *Bot, method string, target any) error {
	url, urlDebug := getURLs(bot, method)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("send GET %s: %w", urlDebug, err)
	}

	errResp, err := parseResponse(resp, target)
	if err != nil {
		return fmt.Errorf("parse response from GET %s: %w", urlDebug, err)
	}

	if errResp != nil {
		return fmt.Errorf("non-OK response from GET %s: %w", urlDebug, errResp)
	}

	return nil
}

// Post sends a POST request to the Telegram API.
// The successful response is decoded into the target.
func (c *Client) Post(bot *Bot, method string, reqData any, target any) error {
	url, urlDebug := getURLs(bot, method)

	reqBody, err := encodeRequestBody(reqData)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Post(url, "application/json", reqBody)
	if err != nil {
		return fmt.Errorf("send POST %s: %w", urlDebug, err)
	}

	errResp, err := parseResponse(resp, target)
	if err != nil {
		return fmt.Errorf("parse response from POST %s: %w", urlDebug, err)
	}

	if errResp != nil {
		return fmt.Errorf("non-OK response from POST %s: %w", urlDebug, errResp)
	}

	return nil
}

// GetMe returns information about the Bot in the TgUser format.
// See: https://core.telegram.org/bots/api#getme
func (c *Client) GetMe(bot *Bot) (*TgUser, error) {
	var resp *TgUserResponse

	if err := c.Get(bot, tgMethodGetMe, &resp); err != nil {
		return nil, err
	}

	return &resp.Result, nil
}

// SendMessage sends a message from the bot via the Telegram API.
// See: https://core.telegram.org/bots/api#sendmessage
func (c *Client) SendMessage(bot *Bot, msg TgMessageRequest) (*TgMessage, error) {
	var resp *TgMessageResponse

	if err := c.Post(bot, tgMethodSendMessage, msg, &resp); err != nil {
		return nil, err
	}

	return &resp.Result, nil
}

func encodeRequestBody(data any) (io.Reader, error) {
	encoded, err := jsoniter.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("encode request body: %w", err)
	}

	return bytes.NewReader(encoded), nil
}

func parseResponse(resp *http.Response, target any) (*TgErrorResponse, error) {
	var (
		body []byte
		err  error
	)

	if resp.Body != nil {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response body: %w", err)
		}
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		if err := jsoniter.Unmarshal(body, target); err != nil {
			return nil, fmt.Errorf("decode response body (%s): %w", string(body), err)
		}

		return nil, nil
	}

	var errResp *TgErrorResponse

	if err := jsoniter.Unmarshal(body, &errResp); err != nil {
		return nil, fmt.Errorf("failed to decode error response body (%s): %w", string(body), err)
	}

	return errResp, nil
}

func getURLs(bot *Bot, method string) (url string, debug string) {
	return tgAPIHost + bot.GetIdentity().String() + "/" + method,
		tgAPIHost + bot.String() + "/" + method
}
