package handlers

import (
	"context"
	"example/admin/auth/cmd/grpc/container"
	"example/admin/auth/cmd/grpc/helpers"
	"example/admin/auth/cmd/grpc/pb"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/session"
	"example/admin/auth/internal/opera/use_cases/sign_out"
	"github.com/selyukovn/go-std"
)

func NewSignOut(ctr *container.Container) func(ctx context.Context, req *pb.SignOutRequest) (*pb.SignOutResponse, error) {
	return func(ctx context.Context, req *pb.SignOutRequest) (*pb.SignOutResponse, error) {
		cl, err := helpers.ParseClient(req.FromIp, req.FromUserAgent)
		if err != nil {
			ctr.Logger.CtxDebugFf(ctx, err.Error())
			return nil, helpers.ErrorInvalidArgument("кривой client")
		}

		sessId, err := session.IdFromString(req.SessionId)
		if err != nil {
			ctr.Logger.CtxDebugFf(ctx, err.Error())
			return nil, helpers.ErrorFailedPrecondition(&pb.ErrorValidationDetail{
				Field:   "SessionId",
				Message: err.Error(),
			})
		}

		// --

		err = ctr.UseCases.SignOut.Execute(sign_out.NewArgs(ctx, cl, sessId))
		switch err.(type) {
		case nil:
		case std.ErrorNotFound:
			return nil, helpers.ErrorNotFound()
		case account.ErrorDeactivated, account.ErrorIpWhitelist:
			return nil, helpers.ErrorFailedPrecondition(&pb.ErrorAccountAccessDeniedDetail{})
		case std.ErrorAlreadyDone:
			return nil, helpers.ErrorFailedPrecondition(&pb.ErrorAlreadyDoneDetail{})
		case std.ErrorRuntime:
			ctr.Logger.CtxErrorFf(ctx, err.Error())
			return nil, helpers.ErrorInternal()
		default:
			panic(err)
		}

		// --

		return &pb.SignOutResponse{}, nil
	}
}
