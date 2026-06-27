package cfm

import (
	"context"
	"example/admin/cfm/internal/api/grpc/handlers/cfm/actions"
	"example/admin/cfm/internal/api/grpc/handlers/cfm/pb"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ pb.CfmServiceServer = sDefault{}

type sDefault struct {
	pb.UnimplementedCfmServiceServer
	createForEmail actions.CreateForEmail
	request        actions.Request
	confirm        actions.Confirm
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func newServiceDefault(
	createForEmail actions.CreateForEmail,
	request actions.Request,
	confirm actions.Confirm,
) sDefault {
	assert.NotNilDeepMust(createForEmail)
	assert.NotNilDeepMust(request)
	assert.NotNilDeepMust(confirm)

	return sDefault{
		createForEmail: createForEmail,
		request:        request,
		confirm:        confirm,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (s sDefault) CreateForEmail(ctx context.Context, req *pb.CreateForEmailRequest) (*pb.CreateForEmailResponse, error) {
	return s.createForEmail(ctx, req)
}

func (s sDefault) Request(ctx context.Context, req *pb.RequestRequest) (*pb.RequestResponse, error) {
	return s.request(ctx, req)
}

func (s sDefault) Confirm(ctx context.Context, req *pb.ConfirmRequest) (*pb.ConfirmResponse, error) {
	return s.confirm(ctx, req)
}

// ---------------------------------------------------------------------------------------------------------------------
