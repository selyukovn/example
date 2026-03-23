package grpc

import (
	"context"
	"example/admin/auth/cmd/grpc/container"
	"example/admin/auth/cmd/grpc/handlers"
	"example/admin/auth/cmd/grpc/pb"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ pb.AuthServiceServer = router{}

type router struct {
	pb.UnimplementedAuthServiceServer
	signInRequest      func(context.Context, *pb.SignInRequestRequest) (*pb.SignInRequestResponse, error)
	signInRequestRetry func(context.Context, *pb.SignInRequestRetryRequest) (*pb.SignInRequestRetryResponse, error)
	signInConfirm      func(context.Context, *pb.SignInConfirmRequest) (*pb.SignInConfirmResponse, error)
	signOut            func(context.Context, *pb.SignOutRequest) (*pb.SignOutResponse, error)
	checkSession       func(context.Context, *pb.CheckSessionRequest) (*pb.CheckSessionResponse, error)
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func newRouter(ctr *container.Container) router {
	return router{
		signInRequest:      handlers.NewSignInRequest(ctr),
		signInRequestRetry: handlers.NewSignInRequestRetry(ctr),
		signInConfirm:      handlers.NewSignInConfirm(ctr),
		signOut:            handlers.NewSignOut(ctr),
		checkSession:       handlers.NewCheckSession(ctr),
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (r router) SignInRequest(ctx context.Context, req *pb.SignInRequestRequest) (*pb.SignInRequestResponse, error) {
	return r.signInRequest(ctx, req)
}

func (r router) SignInRequestRetry(ctx context.Context, req *pb.SignInRequestRetryRequest) (*pb.SignInRequestRetryResponse, error) {
	return r.signInRequestRetry(ctx, req)
}

func (r router) SignInConfirm(ctx context.Context, req *pb.SignInConfirmRequest) (*pb.SignInConfirmResponse, error) {
	return r.signInConfirm(ctx, req)
}

func (r router) SignOut(ctx context.Context, req *pb.SignOutRequest) (*pb.SignOutResponse, error) {
	return r.signOut(ctx, req)
}

func (r router) CheckSession(ctx context.Context, req *pb.CheckSessionRequest) (*pb.CheckSessionResponse, error) {
	return r.checkSession(ctx, req)
}

// ---------------------------------------------------------------------------------------------------------------------
