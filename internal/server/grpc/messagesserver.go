//go:build !no_grpc

package grpc

import (
	"context"
	"fmt"
	"github.com/kukymbr/tgnotifier/internal/api/grpc"
	"github.com/kukymbr/tgnotifier/internal/config"
	"github.com/kukymbr/tgnotifier/internal/types"
	"github.com/kukymbr/tgnotifier/pkg/tgkit"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Sender interface {
	Send(ctx context.Context, options types.SendOptions) (tgkit.TgMessage, error)
}

func registerMessagesServer(grpcSrv *grpcpkg.Server, conf *config.Config, sender Sender) {
	grpc.RegisterMessagesServer(grpcSrv, &messagesServer{
		sender: sender,
		conf:   conf,
	})
}

type messagesServer struct {
	grpc.UnimplementedMessagesServer

	sender Sender
	conf   *config.Config
}

func (s *messagesServer) Send(ctx context.Context, req *grpc.SendMessageRequest) (*grpc.SendMessageResponse, error) {
	opt, err := s.getSendOptions(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tgResp, err := s.sender.Send(ctx, opt)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpc.SendMessageResponse{
		Ok: true,
		Result: &grpc.TgMessage{
			MessageId:       int64(tgResp.MessageID),
			Date:            tgResp.Date,
			MessageThreadId: int64(tgResp.MessageThreadId),
			From: &grpc.TgUser{
				Id:                      tgResp.From.ID,
				FirstName:               tgResp.From.FirstName,
				LastName:                tgResp.From.LastName,
				Username:                tgResp.From.Username,
				LanguageCode:            tgResp.From.LanguageCode,
				IsBot:                   tgResp.From.IsBot,
				IsPremium:               tgResp.From.IsPremium,
				AddedToAttachmentMenu:   tgResp.From.AddedToAttachmentMenu,
				CanJoinGroups:           tgResp.From.CanJoinGroups,
				CanReadAllGroupMessages: tgResp.From.CanReadAllGroupMessages,
				SupportsInlineQueries:   tgResp.From.SupportsInlineQueries,
				CanConnectToBusiness:    tgResp.From.CanConnectToBusiness,
			},
		},
	}, nil
}

func (s *messagesServer) getSendOptions(req *grpc.SendMessageRequest) (types.SendOptions, error) {
	opt := types.SendOptions{}

	if err := opt.BotName.Set(req.GetBotName()); err != nil {
		return types.SendOptions{}, err
	}

	if err := opt.ChatName.Set(req.GetChatName()); err != nil {
		return types.SendOptions{}, err
	}

	msg := req.GetMessage()
	if msg == nil {
		return types.SendOptions{}, fmt.Errorf("message is required")
	}

	if err := opt.Message.ParseMode.Set(msg.GetParseMode()); err != nil {
		return types.SendOptions{}, err
	}

	opt.Message = types.MessageOptions{
		Text:                msg.GetText(),
		DisableNotification: msg.GetDisableNotification(),
		ProtectContent:      msg.GetProtectContent(),
	}

	if err := opt.Validate(); err != nil {
		return types.SendOptions{}, err
	}

	return opt, nil
}
