//go:build grpc_tests

package api_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kukymbr/tgnotifier/internal/api/grpc"
	"github.com/stretchr/testify/suite"
	grpcpkg "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"
)

const (
	envTestGRPCHost = "TEST_GRPC_HOST"
	envTestGRPCPort = "TEST_GRPC_PORT"

	testGRPCHostDefault = "localhost"
	testGRPCPortDefault = "50051"

	envTestBotName  = "TEST_BOT_NAME"
	envTestChatName = "TEST_CHAT_NAME"
)

func TestGRPCSuite(t *testing.T) {
	suite.Run(t, new(grpcTestSuite))
}

type grpcTestSuite struct {
	suite.Suite

	botName        string
	chatName       string
	messagesClient grpc.MessagesClient
}

func (s *grpcTestSuite) SetupSuite() {
	s.botName = os.Getenv(envTestBotName)
	s.chatName = os.Getenv(envTestChatName)

	s.Require().NotEmpty(s.botName)
	s.Require().NotEmpty(s.chatName)

	s.T().Logf("botName=%s; chatName=%s", s.botName, s.chatName)

	s.runServer()

	grpcHost := os.Getenv(envTestGRPCHost)
	if grpcHost == "" {
		grpcHost = testGRPCHostDefault
	}

	grpcPort := os.Getenv(envTestGRPCPort)
	if grpcPort == "" {
		grpcPort = testGRPCPortDefault
	}

	grpcAddress := net.JoinHostPort(grpcHost, grpcPort)

	s.T().Logf("gRPC host to test: %s", grpcAddress)

	client, err := grpcpkg.NewClient(grpcAddress, grpcpkg.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.FailNow("failed to create gRPC client: " + err.Error())
	}

	s.messagesClient = grpc.NewMessagesClient(client)
}

func (s *grpcTestSuite) runServer() {
	ctx, cancel := context.WithCancel(context.Background())

	s.T().Cleanup(cancel)
	s.T().Logf("Running test gRPC server...")

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd := exec.Command(
		"go",
		"build",
		"-o", "./testdata/grpc_test_tgnotifier",
		"./../../../cmd/tgnotifier",
	)

	s.T().Cleanup(func() {
		_ = os.Remove("./testdata/grpc_test_tgnotifier")
	})

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	s.Require().NoError(cmd.Run(), fmt.Sprintf("stdout=%s; stderr=%s", stdout.String(), stderr.String()))

	cmd = exec.CommandContext(
		ctx,
		"./testdata/grpc_test_tgnotifier",
		"grpc",
		"--config=./../../../internal/api/tests/testdata/configs/.tgnotifier.grpc_tests.yml",
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	s.Require().NoError(cmd.Start())

	go func() {
		s.Require().NoError(cmd.Wait())
	}()

	timer := time.NewTimer(500 * time.Millisecond)

	select {
	case <-ctx.Done():
		return
	case <-timer.C:
		s.T().Logf("server stdout: %s", stdout.String())
		s.T().Logf("server stderr: %s", stderr.String())

		break
	}

	s.Require().Empty(stderr.String())
}
