package tgkit

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// RequestRetrier is a tool to retry failed requests.
type RequestRetrier interface {
	Do(reqFn sendRequestFn) (*http.Response, error)
}

// NewLinearRetrier returns a RequestRetrier attempting to send request with a fixed delay.
func NewLinearRetrier(attempts uint, delay time.Duration) RequestRetrier {
	return &linearRetrier{
		attempts: attempts,
		delay:    delay,
	}
}

// NewProgressiveRetrier returns a RequestRetrier attempting to send request with an increasing delay.
func NewProgressiveRetrier(attempts uint, initialDelay time.Duration, multiplier float64) RequestRetrier {
	return &progressiveRetrier{
		attempts:     attempts,
		initialDelay: initialDelay,
		multiplier:   multiplier,
	}
}

// NewNoopRetrier returns a RequestRetrier with a single request attempt.
func NewNoopRetrier() RequestRetrier {
	return &noopRetrier{}
}

type linearRetrier struct {
	attempts uint
	delay    time.Duration
}

func (r *linearRetrier) Do(reqFn sendRequestFn) (*http.Response, error) {
	return doReqFnWithRetries(reqFn, r.attempts, func(_ uint) time.Duration {
		return r.delay
	})
}

type progressiveRetrier struct {
	attempts     uint
	initialDelay time.Duration
	multiplier   float64
}

func (r *progressiveRetrier) Do(reqFn sendRequestFn) (*http.Response, error) {
	return doReqFnWithRetries(reqFn, r.attempts, func(attemptN uint) time.Duration {
		if attemptN == 0 {
			return r.initialDelay
		}

		multiplier := r.multiplier
		if multiplier <= 0.001 {
			multiplier = 1
		}

		return r.initialDelay * time.Duration(float64(attemptN)*multiplier)
	})
}

type noopRetrier struct{}

func (r *noopRetrier) Do(reqFn sendRequestFn) (*http.Response, error) {
	return doReqFn(reqFn)
}

func doReqFn(reqFn sendRequestFn) (*http.Response, error) {
	resp, err := reqFn()

	if err == nil && resp.StatusCode < http.StatusInternalServerError {
		return resp, nil
	}

	if err == nil {
		body, _ := io.ReadAll(resp.Body)
		defer func() {
			_ = resp.Body.Close()
		}()

		return nil, fmt.Errorf(
			"%w: code %d: %s",
			ErrInternalServerErrorResponse,
			resp.StatusCode,
			string(body),
		)
	}

	return nil, err
}

func doReqFnWithRetries(reqFn sendRequestFn, attempts uint, delayFn retrierDelayGetterFn) (*http.Response, error) {
	if attempts == 0 {
		attempts = 1
	}

	var lastErr error

	for i := uint(0); i < attempts; i++ {
		resp, err := doReqFn(reqFn)
		lastErr = err

		if err == nil {
			return resp, nil
		}

		delay := delayFn(i)
		if delay <= 0 {
			delay = 500 * time.Millisecond
		}

		<-time.After(delay)
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", attempts, lastErr)
}

type sendRequestFn func() (*http.Response, error)

type retrierDelayGetterFn func(attemptN uint) time.Duration
