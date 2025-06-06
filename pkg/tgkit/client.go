package tgkit

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	tgAPIHost           = "https://api.telegram.org/"
	tgMethodGetUpdates  = "getUpdates"
	tgMethodGetMe       = "getMe"
	tgMethodSendMessage = "sendMessage"
)

// ErrInternalServerErrorResponse is an error indicating Telegram API responded with a code >= 500.
var ErrInternalServerErrorResponse = errors.New("got internal server error response from Telegram API")

// NewDefaultClient creates a new default telegram Client instance.
func NewDefaultClient() *Client {
	return NewClientWithOptions(
		WithHTTPClient(getDefaultHTTPClient()),
		WithRetrier(getDefaultRequestRetrier()),
	)
}

// NewClientWithOptions creates new Client instance with a construction options.
func NewClientWithOptions(options ...ClientOption) *Client {
	c := &Client{}

	for _, opt := range options {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = getDefaultHTTPClient()
	}

	if c.retrier == nil {
		c.retrier = NewNoopRetrier()
	}

	return c
}

// Client is a tool to communicate with a Telegram API via the HTTPS.
type Client struct {
	httpClient HTTPClient
	retrier    RequestRetrier
}

// Get sends a GET request to the Telegram API.
// The successful response is decoded into the target.
func (c *Client) Get(bot Bot, method string, target any) error {
	url, urlDebug := getURLs(bot, method)

	resp, err := c.retrier.Do(func() (*http.Response, error) {
		return c.httpClient.Get(url)
	})
	if err != nil {
		return fmt.Errorf("send GET %s: %w", urlDebug, err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

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
func (c *Client) Post(bot Bot, method string, reqData any, target any) error {
	url, urlDebug := getURLs(bot, method)

	reqBody, err := encodeRequestBody(reqData)
	if err != nil {
		return err
	}

	resp, err := c.retrier.Do(func() (*http.Response, error) {
		return c.httpClient.Post(url, "application/json", reqBody)
	})
	if err != nil {
		return fmt.Errorf("send POST %s: %w", urlDebug, err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

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
func (c *Client) GetMe(bot Bot) (TgUser, error) {
	var resp TgUserResponse

	if err := c.Get(bot, tgMethodGetMe, &resp); err != nil {
		return TgUser{}, err
	}

	return resp.Result, nil
}

// GetUpdates returns last bot's updates.
// See https://core.telegram.org/bots/api#getupdates
// TODO: pass parameters.
func (c *Client) GetUpdates(bot Bot) ([]TgUpdate, error) {
	var resp TgUpdatesResponse

	if err := c.Get(bot, tgMethodGetUpdates, &resp); err != nil {
		return nil, err
	}

	return resp.Result, nil
}

// SendMessage sends a message from the bot via the Telegram API.
// See: https://core.telegram.org/bots/api#sendmessage
func (c *Client) SendMessage(bot Bot, msg TgMessageRequest) (TgMessage, error) {
	var resp TgMessageResponse

	if err := c.Post(bot, tgMethodSendMessage, msg, &resp); err != nil {
		return TgMessage{}, err
	}

	return resp.Result, nil
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
		defer func() {
			_ = resp.Body.Close()
		}()

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

func getURLs(bot Bot, method string) (url string, debug string) {
	return tgAPIHost + bot.GetIdentity().String() + "/" + method,
		tgAPIHost + bot.String() + "/" + method
}

func getDefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

func getDefaultRequestRetrier() RequestRetrier {
	return NewProgressiveRetrier(3, 500*time.Millisecond, 2)
}
