package tgkit

// ClientOption is a client constructor configuration option.
type ClientOption func(*Client)

// WithHTTPClient defines an HTTPClient instance to use with a Client.
func WithHTTPClient(httpClient HTTPClient) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithRetrier defines a way or requests retries using the RequestRetrier instance.
func WithRetrier(retrier RequestRetrier) ClientOption {
	return func(c *Client) {
		c.retrier = retrier
	}
}
