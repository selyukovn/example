package grpc

import (
	"context"
	"example/admin/cfm/internal/api/grpc/handlers"
	"example/admin/cfm/internal/api/grpc/pb"
	"example/admin/cfm/internal/opera/use_cases/confirm"
	"example/admin/cfm/internal/opera/use_cases/create_for_email"
	"example/admin/cfm/internal/opera/use_cases/request"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ pb.CfmServiceServer = Router{}

type Router struct {
	pb.UnimplementedCfmServiceServer
	createForEmail func(ctx context.Context, req *pb.CreateForEmailRequest) (*pb.CreateForEmailResponse, error)
	request        func(ctx context.Context, req *pb.RequestRequest) (*pb.RequestResponse, error)
	confirm        func(ctx context.Context, req *pb.ConfirmRequest) (*pb.ConfirmResponse, error)
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewRouter(
	ucCreateForEmail create_for_email.Command,
	ucRequest request.Command,
	ucConfirm confirm.Command,
) Router {
	return Router{
		createForEmail: handlers.NewCreateForEmail(ucCreateForEmail),
		request:        handlers.NewRequest(ucRequest),
		confirm:        handlers.NewConfirm(ucConfirm),
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (r Router) CreateForEmail(ctx context.Context, req *pb.CreateForEmailRequest) (*pb.CreateForEmailResponse, error) {
	return r.createForEmail(ctx, req)
}

func (r Router) Request(ctx context.Context, req *pb.RequestRequest) (*pb.RequestResponse, error) {
	return r.request(ctx, req)
}

func (r Router) Confirm(ctx context.Context, req *pb.ConfirmRequest) (*pb.ConfirmResponse, error) {
	return r.confirm(ctx, req)
}

// ---------------------------------------------------------------------------------------------------------------------
