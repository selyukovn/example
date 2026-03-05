package grpc

import (
	"example/admin/cfm/cmd/common/container"
	"example/admin/cfm/cmd/grpc/interceptors"
	"example/admin/cfm/cmd/grpc/pb"
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
func NewServer(ctr *container.Container, apiKey string) *Server {
	assert.NotNilDeepMust(ctr)
	assert.Str().NotEmpty().Must(apiKey)

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.NewBoundary(ctr),
		interceptors.NewAccessKey(apiKey),
		interceptors.NewPrometheusMetrics(),
	))
	pb.RegisterCfmServiceServer(s, newRouter(ctr))

	return &Server{s: s}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	return s.s.Serve(lis)
}

func (s *Server) Stop() error {
	s.s.GracefulStop()
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
