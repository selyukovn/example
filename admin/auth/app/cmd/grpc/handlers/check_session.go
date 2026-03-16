package handlers

import (
	"context"
	"example/admin/auth/cmd/grpc/container"
	"example/admin/auth/cmd/grpc/helpers"
	"example/admin/auth/cmd/grpc/pb"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/session"
	"example/admin/auth/internal/opera/use_cases/check_session"
	"github.com/selyukovn/go-std"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewCheckSession(ctr *container.Container) func(ctx context.Context, req *pb.CheckSessionRequest) (*pb.CheckSessionResponse, error) {
	return func(ctx context.Context, req *pb.CheckSessionRequest) (*pb.CheckSessionResponse, error) {
		cl, err := helpers.ParseClient(req.FromIp, req.FromUserAgent)
		if err != nil {
			ctr.Logger.CtxDebugFf(ctx, err.Error())
			return nil, helpers.ErrorInvalidArgument("кривой client")
		}

		sessId, err := session.IdFromString(req.SessionId)
		if err != nil {
			ctr.Logger.CtxDebugFf(ctx, err.Error())
			return nil, helpers.ErrorFailedPrecondition(&pb.ErrorValidationDetail{
				Field:   "sessionId",
				Message: err.Error(),
			})
		}

		// --

		res, err := ctr.UseCases.CheckSession.Execute(check_session.NewArgs(ctx, cl, sessId))
		switch vErr := err.(type) {
		case nil:
		case std.ErrorNotFound:
			return nil, helpers.ErrorNotFound()
		case account.ErrorDeactivated, account.ErrorIpWhitelist:
			return nil, helpers.ErrorFailedPrecondition(&pb.ErrorAccountAccessDeniedDetail{})
		case session.ErrorClosed:
			return nil, helpers.ErrorFailedPrecondition(&pb.ErrorSessionClosedDetail{
				IsExpired: vErr.IsExpired(),
				ClosedAt:  timestamppb.New(vErr.ClosedAt()),
			})
		case std.ErrorRuntime:
			ctr.Logger.CtxErrorFf(ctx, err.Error())
			return nil, helpers.ErrorInternal()
		default:
			panic(err)
		}

		// --

		return &pb.CheckSessionResponse{
			AccountId:       res.AccId().String(),
			SessionExpireAt: timestamppb.New(res.SessExpAt()),
		}, nil
	}
}
