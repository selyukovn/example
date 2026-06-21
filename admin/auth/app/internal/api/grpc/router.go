package grpc

import (
	"context"
	"example/admin/auth/internal/api/grpc/handlers"
	"example/admin/auth/internal/api/grpc/pb"
	"example/admin/auth/internal/opera/use_cases/check_session"
	"example/admin/auth/internal/opera/use_cases/sign_in_confirm"
	"example/admin/auth/internal/opera/use_cases/sign_in_request"
	"example/admin/auth/internal/opera/use_cases/sign_in_request_retry"
	"example/admin/auth/internal/opera/use_cases/sign_out"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ pb.AuthServiceServer = Router{}

type Router struct {
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

func NewRouter(
	ucSignInRequest sign_in_request.Command,
	ucSignInRequestRetry sign_in_request_retry.Command,
	ucSignInConfirm sign_in_confirm.Command,
	ucSignOut sign_out.Command,
	ucCheckSession check_session.Command,
) Router {
	return Router{
		signInRequest:      handlers.NewSignInRequest(ucSignInRequest),
		signInRequestRetry: handlers.NewSignInRequestRetry(ucSignInRequestRetry),
		signInConfirm:      handlers.NewSignInConfirm(ucSignInConfirm),
		signOut:            handlers.NewSignOut(ucSignOut),
		checkSession:       handlers.NewCheckSession(ucCheckSession),
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (r Router) SignInRequest(ctx context.Context, req *pb.SignInRequestRequest) (*pb.SignInRequestResponse, error) {
	return r.signInRequest(ctx, req)
}

func (r Router) SignInRequestRetry(ctx context.Context, req *pb.SignInRequestRetryRequest) (*pb.SignInRequestRetryResponse, error) {
	return r.signInRequestRetry(ctx, req)
}

func (r Router) SignInConfirm(ctx context.Context, req *pb.SignInConfirmRequest) (*pb.SignInConfirmResponse, error) {
	return r.signInConfirm(ctx, req)
}

func (r Router) SignOut(ctx context.Context, req *pb.SignOutRequest) (*pb.SignOutResponse, error) {
	return r.signOut(ctx, req)
}

func (r Router) CheckSession(ctx context.Context, req *pb.CheckSessionRequest) (*pb.CheckSessionResponse, error) {
	return r.checkSession(ctx, req)
}

// ---------------------------------------------------------------------------------------------------------------------
