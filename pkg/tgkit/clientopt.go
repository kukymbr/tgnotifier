package tgkit

import "time"

// ClientOption
type ClientOption func(*Client)

func WithHTTPClient(httpClient HTTPClient) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithRetry(attempts int, delay time.Duration) ClientOption {
	return func(c *Client) {
		c.retryAttempts = attempts
		c.retryDelay = delay
	}
}
