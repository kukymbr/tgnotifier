//go:build grpc_tests

package api_test

import (
	"context"
	"github.com/kukymbr/tgnotifier/internal/api/grpc"
)

func (s *grpcTestSuite) TestSendMessage() {
	resp, err := s.messagesClient.Send(context.Background(), &grpc.SendMessageRequest{
		BotName:  &s.botName,
		ChatName: &s.chatName,
		Message: &grpc.MessageRequest{
			Text:                getPtr("😎 test_message"),
			DisableNotification: getPtr(true),
		},
	})

	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.True(resp.GetOk())
	s.True(resp.GetResult().GetFrom().GetIsBot())
}

func getPtr[T any](val T) *T {
	return &val
}
