//go:build !no_grpc

package grpc

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/kukymbr/tgnotifier/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

// New creates new gRPC server instance.
func New(conf *config.Config, sender Sender) *Server {
	grpcSrv := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(
			recovery.WithRecoveryHandler(func(p interface{}) (err error) {
				return status.Errorf(codes.Internal, fmt.Sprintf("unhandled panic: %v", p))
			}),
		),
	))

	return &Server{
		conf:       conf,
		sender:     sender,
		grpcServer: grpcSrv,
	}
}

// Server is an gRPC server.
type Server struct {
	conf       *config.Config
	sender     Sender
	grpcServer *grpc.Server
}

// Run starts the gRPC server.
func (s *Server) Run() error {
	registerMessagesServer(s.grpcServer, s.conf, s.sender)

	port := s.conf.GRPC().GetPort()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to create TCP listener: %w", err)
	}

	fmt.Printf("Starting gRPC server on port %d...\n", port)

	if err := s.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to start gRPC server: %w", err)
	}

	return nil
}

// Close stops the gRPC server.
func (s *Server) Close() error {
	fmt.Println("Stopping the gRPC server...")

	s.grpcServer.GracefulStop()

	return nil
}
