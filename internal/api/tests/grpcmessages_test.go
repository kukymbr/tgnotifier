//go:build grpc_tests

package api_test

import (
	"context"
	"github.com/kukymbr/tgnotifier/internal/api/grpc"
)

func (s *grpcTestSuite) TestSendMessage() {
	resp, err := s.messagesClient.Send(context.Background(), &grpc.SendMessageRequest{
		BotName:  s.botName,
		ChatName: s.chatName,
		Message: &grpc.MessageRequest{
			Text:                "😎 test_message",
			DisableNotification: true,
		},
	})

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.True(resp.GetOk())
	s.True(resp.GetResult().GetFrom().GetIsBot())
}
