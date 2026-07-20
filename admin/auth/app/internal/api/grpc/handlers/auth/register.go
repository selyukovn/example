package auth

import (
	"example/admin/auth/internal/api/grpc/handlers/auth/actions"
	"example/admin/auth/internal/api/grpc/handlers/auth/pb"
	"example/admin/auth/internal/opera/use_cases/check_session"
	"example/admin/auth/internal/opera/use_cases/sign_in_confirm"
	"example/admin/auth/internal/opera/use_cases/sign_in_request"
	"example/admin/auth/internal/opera/use_cases/sign_in_request_retry"
	"example/admin/auth/internal/opera/use_cases/sign_out"
	assert "github.com/selyukovn/go-wm-assert"
	"google.golang.org/grpc"
)

func Register(
	grpcServer *grpc.Server,
	ucSignInRequest sign_in_request.Command,
	ucSignInRequestRetry sign_in_request_retry.Command,
	ucSignInConfirm sign_in_confirm.Command,
	ucSignOut sign_out.Command,
	ucCheckSession check_session.Command,
	sMws []func(pb.AuthServiceServer) pb.AuthServiceServer,
) {
	assert.NotNilDeepMust(grpcServer)

	var service pb.AuthServiceServer = newServiceDefault(
		actions.NewSignInRequest(ucSignInRequest),
		actions.NewSignInRequestRetry(ucSignInRequestRetry),
		actions.NewSignInConfirm(ucSignInConfirm),
		actions.NewSignOut(ucSignOut),
		actions.NewCheckSession(ucCheckSession),
	)

	for i := len(sMws) - 1; i >= 0; i-- {
		service = sMws[i](service)
	}

	pb.RegisterAuthServiceServer(grpcServer, service)
}
