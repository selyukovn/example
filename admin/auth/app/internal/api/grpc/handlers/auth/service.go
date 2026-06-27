package auth

import (
	"context"
	"example/admin/auth/internal/api/grpc/handlers/auth/actions"
	"example/admin/auth/internal/api/grpc/handlers/auth/pb"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ pb.AuthServiceServer = sDefault{}

type sDefault struct {
	pb.UnimplementedAuthServiceServer
	signInRequest      actions.SignInRequest
	signInRequestRetry actions.SignInRequestRetry
	signInConfirm      actions.SignInConfirm
	signOut            actions.SignOut
	checkSession       actions.CheckSession
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func newServiceDefault(
	signInRequest actions.SignInRequest,
	signInRequestRetry actions.SignInRequestRetry,
	signInConfirm actions.SignInConfirm,
	signOut actions.SignOut,
	checkSession actions.CheckSession,
) sDefault {
	assert.NotNilDeepMust(signInRequest)
	assert.NotNilDeepMust(signInRequestRetry)
	assert.NotNilDeepMust(signInConfirm)
	assert.NotNilDeepMust(signOut)
	assert.NotNilDeepMust(checkSession)

	return sDefault{
		signInRequest:      signInRequest,
		signInRequestRetry: signInRequestRetry,
		signInConfirm:      signInConfirm,
		signOut:            signOut,
		checkSession:       checkSession,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (r sDefault) SignInRequest(ctx context.Context, req *pb.SignInRequestRequest) (*pb.SignInRequestResponse, error) {
	return r.signInRequest(ctx, req)
}

func (r sDefault) SignInRequestRetry(ctx context.Context, req *pb.SignInRequestRetryRequest) (*pb.SignInRequestRetryResponse, error) {
	return r.signInRequestRetry(ctx, req)
}

func (r sDefault) SignInConfirm(ctx context.Context, req *pb.SignInConfirmRequest) (*pb.SignInConfirmResponse, error) {
	return r.signInConfirm(ctx, req)
}

func (r sDefault) SignOut(ctx context.Context, req *pb.SignOutRequest) (*pb.SignOutResponse, error) {
	return r.signOut(ctx, req)
}

func (r sDefault) CheckSession(ctx context.Context, req *pb.CheckSessionRequest) (*pb.CheckSessionResponse, error) {
	return r.checkSession(ctx, req)
}

// ---------------------------------------------------------------------------------------------------------------------
