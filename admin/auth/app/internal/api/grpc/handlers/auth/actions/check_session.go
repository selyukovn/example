package actions

import (
	"context"
	"example/admin/auth/internal/api/grpc/handlers/auth/pb"
	"example/admin/auth/internal/api/grpc/kernel"
	"example/admin/auth/internal/api/grpc/kernel_extra"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/session"
	"example/admin/auth/internal/opera/use_cases/check_session"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ---------------------------------------------------------------------------------------------------------------------

type CheckSession = func(ctx context.Context, req *pb.CheckSessionRequest) (*pb.CheckSessionResponse, error)

// ---------------------------------------------------------------------------------------------------------------------

func NewCheckSession(ucCheckSession check_session.Command) CheckSession {
	return func(ctx context.Context, req *pb.CheckSessionRequest) (*pb.CheckSessionResponse, error) {
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

		res, err := ucCheckSession.Execute(ctx, cl, sessId)
		switch vErr := err.(type) {
		case nil:
		case std.ErrorNotFound:
			return nil, kernel.ErrorNotFound()
		case account.ErrorDeactivated, account.ErrorIpWhitelist:
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorAccountAccessDeniedDetail{})
		case session.ErrorClosed:
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorSessionClosedDetail{
				IsExpired: vErr.IsExpired(),
				ClosedAt:  timestamppb.New(vErr.ClosedAt()),
			})
		case std.ErrorRuntime:
			logger.ErrorFf(ctx, err.Error())
			return nil, kernel.ErrorInternal()
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

// ---------------------------------------------------------------------------------------------------------------------
