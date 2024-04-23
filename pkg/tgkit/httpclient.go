package tgkit

import (
	"io"
	"net/http"
)

// HTTPClient is an HTTP client interface; the http.Client fits.
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}
