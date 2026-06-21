package handlers

import (
	"context"
	"example/admin/auth/internal/api/grpc/kernel"
	"example/admin/auth/internal/api/grpc/kernel_extra"
	"example/admin/auth/internal/api/grpc/pb"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/session"
	"example/admin/auth/internal/opera/use_cases/sign_out"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
)

func NewSignOut(ucSignOut sign_out.Command) func(ctx context.Context, req *pb.SignOutRequest) (*pb.SignOutResponse, error) {
	return func(ctx context.Context, req *pb.SignOutRequest) (*pb.SignOutResponse, error) {
		cl, err := kernel_extra.ParseClient(req.FromIp, req.FromUserAgent)
		if err != nil {
			logger.DebugFf(ctx, err.Error())
			return nil, kernel.ErrorInvalidArgument("кривой client")
		}

		sessId, err := session.IdFromString(req.SessionId)
		if err != nil {
			logger.DebugFf(ctx, err.Error())
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorValidationDetail{
				Field:   "sessionId",
				Message: err.Error(),
			})
		}

		// --

		err = ucSignOut.Execute(sign_out.NewArgs(ctx, cl, sessId))
		switch err.(type) {
		case nil:
		case std.ErrorNotFound:
			return nil, kernel.ErrorNotFound()
		case account.ErrorDeactivated, account.ErrorIpWhitelist:
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorAccountAccessDeniedDetail{})
		case std.ErrorAlreadyDone:
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorAlreadyDoneDetail{})
		case std.ErrorRuntime:
			logger.ErrorFf(ctx, err.Error())
			return nil, kernel.ErrorInternal()
		default:
			panic(err)
		}

		// --

		return &pb.SignOutResponse{}, nil
	}
}
