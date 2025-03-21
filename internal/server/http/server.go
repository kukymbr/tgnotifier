//go:build !no_http

package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	httpstubs "github.com/kukymbr/tgnotifier/internal/api/http"
	"github.com/kukymbr/tgnotifier/internal/config"
	"github.com/kukymbr/tgnotifier/internal/sender"
	"github.com/kukymbr/tgnotifier/internal/types"
	"github.com/kukymbr/tgnotifier/internal/util"
	"github.com/kukymbr/tgnotifier/pkg/tgkit"
)

const readHeaderTimeout = 30 * time.Second

// RunServer runs an HTTP server.
func RunServer(ctx context.Context, conf *config.Config, sender *sender.Sender) error {
	handler := httpstubs.Handler(&messagesService{
		conf:   conf,
		sender: sender,
	})

	server := &http.Server{
		Addr:              conf.HTTP().GetAddress(),
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		fmt.Printf("Staring HTTP server on %s\n", server.Addr)

		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			util.PrintlnError(fmt.Errorf("HTTP server failed: %w", err))
		}
	}()

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		<-ctx.Done()

		if err := server.Shutdown(ctx); err != nil {
			util.PrintlnError(fmt.Errorf("HTTP server shutdown failed: %w", err))
		}
	}()

	wg.Wait()

	return nil
}

type messagesService struct {
	conf   *config.Config
	sender *sender.Sender
}

func (s *messagesService) Send(w http.ResponseWriter, r *http.Request) {
	opt, err := s.parseRequest(r)
	if err != nil {
		s.responseWithError(w, http.StatusBadRequest, err)

		return
	}

	tgResp, err := s.sender.Send(r.Context(), opt)
	if err != nil {
		s.responseWithError(w, http.StatusInternalServerError, err)

		return
	}

	s.responseOk(w, tgResp)
}

func (s *messagesService) parseRequest(r *http.Request) (types.SendOptions, error) {
	if r.Body == nil {
		return types.SendOptions{}, fmt.Errorf("empty request body")
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			util.PrintlnError(fmt.Errorf("failed to close request body: %w", err))
		}
	}()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return types.SendOptions{}, fmt.Errorf("failed to read request body: %w", err)
	}

	var req httpstubs.SendMessageRequest

	if err := jsoniter.Unmarshal(body, &req); err != nil {
		return types.SendOptions{}, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	opt := types.SendOptions{
		Message: types.MessageOptions{
			Text:                req.Message.Text,
			DisableNotification: getValue(req.Message.DisableNotification),
			ProtectContent:      getValue(req.Message.ProtectContent),
		},
	}

	if err := opt.BotName.Set(getValue(req.BotName)); err != nil {
		return types.SendOptions{}, err
	}

	if err := opt.ChatName.Set(getValue(req.ChatName)); err != nil {
		return types.SendOptions{}, err
	}

	if err := opt.Message.ParseMode.Set(getValue(req.Message.ParseMode)); err != nil {
		return types.SendOptions{}, err
	}

	if err := opt.Validate(); err != nil {
		return types.SendOptions{}, err
	}

	return opt, nil
}

func (s *messagesService) responseOk(w http.ResponseWriter, resp tgkit.TgMessage) {
	w.Header().Add("Content-Type", "application/json")

	data, err := jsoniter.Marshal(resp)
	if err != nil {
		s.responseWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to marshal response: %w", err))

		return
	}

	_, _ = w.Write(data)
}

func (s *messagesService) responseWithError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, _ = w.Write([]byte(`{"error": "` + err.Error() + `"}`))
}

func getValue[T any](ptr *T) T {
	var empty T

	if ptr == nil {
		return empty
	}

	return *ptr
}
