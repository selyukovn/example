package grpc

import (
	"context"
	"example/admin/auth/cmd/grpc/container"
	"example/admin/auth/cmd/grpc/interceptors"
	"example/admin/auth/cmd/grpc/pb"
	assert "github.com/selyukovn/go-wm-assert"
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

// NewServer
//
// Паникует при нулевых аргументах.
func NewServer(ctr *container.Container, apiKey string) Server {
	assert.NotNilDeepMust(ctr)
	assert.Str().NotEmpty().Must(apiKey)

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.NewBoundary(),
		interceptors.NewAccessKey(apiKey),
		interceptors.NewPrometheusMetrics(),
	))
	pb.RegisterAuthServiceServer(s, newRouter(ctr))

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
