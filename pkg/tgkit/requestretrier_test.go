package tgkit_test

import (
	"bytes"
	"fmt"
	"github.com/kukymbr/tgnotifier/pkg/tgkit"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
	"time"
)

func testReqFnOK() (*http.Response, error) {
	body := io.NopCloser(bytes.NewReader([]byte("OK")))

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       body,
	}, nil
}

func testReqFn500() (*http.Response, error) {
	body := io.NopCloser(bytes.NewReader([]byte("Internal Server Error")))

	return &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       body,
	}, nil
}

func testReqFnError() (*http.Response, error) {
	return nil, fmt.Errorf("test_error")
}

func TestRequestRetrier(t *testing.T) {
	tests := []struct {
		Name       string
		GetRetrier func() tgkit.RequestRetrier
		ReqFn      func() (*http.Response, error)
		AssertResp func(t *testing.T, resp *http.Response, err error)
		AssertTime func(t *testing.T, took time.Duration)
	}{
		{
			Name: "noop retrier",
			GetRetrier: func() tgkit.RequestRetrier {
				return tgkit.NewNoopRetrier()
			},
			ReqFn: testReqFn500,
			AssertResp: func(t *testing.T, resp *http.Response, err error) {
				assert.Nil(t, resp)
				assert.ErrorIs(t, err, tgkit.ErrInternalServerErrorResponse)
			},
			AssertTime: func(t *testing.T, took time.Duration) {
				assert.Less(t, took, time.Millisecond)
			},
		},
		{
			Name: "linear retrier when success",
			GetRetrier: func() tgkit.RequestRetrier {
				return tgkit.NewLinearRetrier(3, time.Second)
			},
			ReqFn: testReqFnOK,
			AssertResp: func(t *testing.T, resp *http.Response, err error) {
				assert.NotNil(t, resp)
				assert.NoError(t, err)
			},
			AssertTime: func(t *testing.T, took time.Duration) {
				assert.Less(t, took, time.Millisecond)
			},
		},
		{
			Name: "linear retrier when 500",
			GetRetrier: func() tgkit.RequestRetrier {
				return tgkit.NewLinearRetrier(3, 10*time.Millisecond)
			},
			ReqFn: testReqFn500,
			AssertResp: func(t *testing.T, resp *http.Response, err error) {
				assert.Nil(t, resp)
				assert.ErrorIs(t, err, tgkit.ErrInternalServerErrorResponse)
			},
			AssertTime: func(t *testing.T, took time.Duration) {
				assert.GreaterOrEqual(t, took, 30*time.Millisecond)
			},
		},
		{
			Name: "linear retrier when error",
			GetRetrier: func() tgkit.RequestRetrier {
				return tgkit.NewLinearRetrier(3, 10*time.Millisecond)
			},
			ReqFn: testReqFnError,
			AssertResp: func(t *testing.T, resp *http.Response, err error) {
				assert.Nil(t, resp)
				assert.Error(t, err)
			},
			AssertTime: func(t *testing.T, took time.Duration) {
				assert.GreaterOrEqual(t, took, 30*time.Millisecond)
			},
		},
		{
			Name: "progressive retrier when success",
			GetRetrier: func() tgkit.RequestRetrier {
				return tgkit.NewProgressiveRetrier(3, time.Second, 2)
			},
			ReqFn: testReqFnOK,
			AssertResp: func(t *testing.T, resp *http.Response, err error) {
				assert.NotNil(t, resp)
				assert.NoError(t, err)
			},
			AssertTime: func(t *testing.T, took time.Duration) {
				assert.Less(t, took, time.Millisecond)
			},
		},
		{
			Name: "progressive retrier when 500",
			GetRetrier: func() tgkit.RequestRetrier {
				return tgkit.NewProgressiveRetrier(3, 5*time.Millisecond, 2)
			},
			ReqFn: testReqFn500,
			AssertResp: func(t *testing.T, resp *http.Response, err error) {
				assert.Nil(t, resp)
				assert.ErrorIs(t, err, tgkit.ErrInternalServerErrorResponse)
			},
			AssertTime: func(t *testing.T, took time.Duration) {
				assert.GreaterOrEqual(t, took, 35*time.Millisecond)
			},
		},
		{
			Name: "progressive retrier when error",
			GetRetrier: func() tgkit.RequestRetrier {
				return tgkit.NewProgressiveRetrier(3, 5*time.Millisecond, 2)
			},
			ReqFn: testReqFnError,
			AssertResp: func(t *testing.T, resp *http.Response, err error) {
				assert.Nil(t, resp)
				assert.Error(t, err)
			},
			AssertTime: func(t *testing.T, took time.Duration) {
				assert.GreaterOrEqual(t, took, 35*time.Millisecond)
			},
		},
		{
			Name: "when zero arguments",
			GetRetrier: func() tgkit.RequestRetrier {
				return tgkit.NewLinearRetrier(0, 0)
			},
			ReqFn: testReqFn500,
			AssertResp: func(t *testing.T, resp *http.Response, err error) {
				assert.Nil(t, resp)
				assert.ErrorIs(t, err, tgkit.ErrInternalServerErrorResponse)
			},
			AssertTime: func(t *testing.T, took time.Duration) {
				assert.GreaterOrEqual(t, took, 500*time.Millisecond)
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			retrier := test.GetRetrier()

			start := time.Now()
			resp, err := retrier.Do(test.ReqFn)
			took := time.Since(start)

			test.AssertResp(t, resp, err)
			test.AssertTime(t, took)
		})
	}
}
