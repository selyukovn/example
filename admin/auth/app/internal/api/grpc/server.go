package grpc

import (
	"context"
	"example/admin/auth/internal/api/grpc/pb"
	"google.golang.org/grpc"
	"net"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Server struct {
	s *grpc.Server
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewServer(router Router, interceptors ...grpc.UnaryServerInterceptor) Server {
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors...))
	pb.RegisterAuthServiceServer(s, router)

	return Server{s: s}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (s Server) Start(ctx context.Context) error {
	// todo : использовать контекст как базовый
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	return s.s.Serve(lis)
}

func (s Server) Stop(ctx context.Context) error {
	// todo : возможно, есть смысл ограничить по времени
	s.s.GracefulStop()
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
