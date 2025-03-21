package tgkit

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

// HTTPClientJSONMock is a test client returning the predefined response.
type HTTPClientJSONMock struct {
	responses map[string]*http.Response
}

// RegisterResponse registers response for the url.
func (m *HTTPClientJSONMock) RegisterResponse(method string, url string, responseCode int, responseBody any) {
	var (
		body []byte
		err  error
	)

	switch r := responseBody.(type) {
	case []byte:
		body = r
	case string:
		body = []byte(r)
	case nil:
	default:
		body, err = jsoniter.Marshal(responseBody)
		if err != nil {
			panic("cannot marshal mock response body: " + err.Error())
		}
	}

	resp := &http.Response{
		Status:        http.StatusText(responseCode),
		StatusCode:    responseCode,
		Body:          io.NopCloser(bytes.NewReader(body)),
		ContentLength: int64(len(body)),
	}

	if m.responses == nil {
		m.responses = make(map[string]*http.Response)
	}

	m.responses[m.key(method, url)] = resp
}

func (m *HTTPClientJSONMock) Get(url string) (resp *http.Response, err error) {
	return m.getResponse(http.MethodGet, url), nil
}

func (m *HTTPClientJSONMock) Post(url, _ string, _ io.Reader) (resp *http.Response, err error) {
	return m.getResponse(http.MethodPost, url), nil
}

func (m *HTTPClientJSONMock) getResponse(method string, url string) *http.Response {
	noRespErr := fmt.Errorf("no test response is registered for '%s %s'", method, url)

	if m.responses == nil {
		panic(noRespErr)
	}

	resp, ok := m.responses[m.key(method, url)]
	if !ok {
		panic(noRespErr)
	}

	return resp
}

func (m *HTTPClientJSONMock) key(method string, url string) string {
	return method + " " + url
}
