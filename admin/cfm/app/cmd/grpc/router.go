package grpc

import (
	"context"
	"example/admin/cfm/cmd/common/container"
	"example/admin/cfm/cmd/grpc/handlers"
	"example/admin/cfm/cmd/grpc/pb"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ pb.CfmServiceServer = router{}

type router struct {
	pb.UnimplementedCfmServiceServer
	createForEmail func(ctx context.Context, req *pb.CreateForEmailRequest) (*pb.CreateForEmailResponse, error)
	request        func(ctx context.Context, req *pb.RequestRequest) (*pb.RequestResponse, error)
	confirm        func(ctx context.Context, req *pb.ConfirmRequest) (*pb.ConfirmResponse, error)
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func newRouter(ctr *container.Container) router {
	return router{
		createForEmail: handlers.NewCreateForEmail(ctr),
		request:        handlers.NewRequest(ctr),
		confirm:        handlers.NewConfirm(ctr),
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (r router) CreateForEmail(ctx context.Context, req *pb.CreateForEmailRequest) (*pb.CreateForEmailResponse, error) {
	return r.createForEmail(ctx, req)
}

func (r router) Request(ctx context.Context, req *pb.RequestRequest) (*pb.RequestResponse, error) {
	return r.request(ctx, req)
}

func (r router) Confirm(ctx context.Context, req *pb.ConfirmRequest) (*pb.ConfirmResponse, error) {
	return r.confirm(ctx, req)
}

// ---------------------------------------------------------------------------------------------------------------------
