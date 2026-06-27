package actions

import (
	"context"
	"example/admin/auth/internal/api/grpc/handlers/auth/pb"
	"example/admin/auth/internal/api/grpc/kernel"
	"example/admin/auth/internal/api/grpc/kernel_extra"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/action_request"
	"example/admin/auth/internal/domain/cfm"
	"example/admin/auth/internal/opera/use_cases/sign_in_request_retry"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ---------------------------------------------------------------------------------------------------------------------

type SignInRequestRetry = func(ctx context.Context, req *pb.SignInRequestRetryRequest) (*pb.SignInRequestRetryResponse, error)

// ---------------------------------------------------------------------------------------------------------------------

func NewSignInRequestRetry(ucSignInRequestRetry sign_in_request_retry.Command) SignInRequestRetry {
	return func(ctx context.Context, req *pb.SignInRequestRetryRequest) (*pb.SignInRequestRetryResponse, error) {
		cl, err := kernel_extra.ParseClient(req.FromIp, req.FromUserAgent)
		if err != nil {
			logger.DebugFf(ctx, err.Error())
			return nil, kernel.ErrorInvalidArgument("кривой client")
		}

		signInId, err := action_request.IdFromString(req.SignInId)
		if err != nil {
			logger.DebugFf(ctx, err.Error())
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorValidationDetail{
				Field:   "SignInId",
				Message: err.Error(),
			})
		}

		// --

		res, err := ucSignInRequestRetry.Execute(ctx, cl, signInId)
		switch vErr := err.(type) {
		case nil:
		case std.ErrorNotFound:
			return nil, kernel.ErrorNotFound()
		case account.ErrorDeactivated, account.ErrorIpWhitelist:
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorAccountAccessDeniedDetail{})
		case cfm.ErrorFinished:
			logger.WarnFf(ctx, "%v обратился к завершенному SignIn %q: %#v", cl, signInId, vErr)
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorSignInFinishedDetail{
				IsPassed:  vErr.IsAsPassed(),
				IsFailed:  vErr.IsAsFailed(),
				IsExpired: vErr.IsAsExpired(),
			})
		case cfm.ErrorNoAttemptsLeft:
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorNoAttemptsLeftDetail{})
		case cfm.ErrorRequestsFrequency:
			// фронт обновляет данные, а не реагирует на статус, поэтому отвечаем как при успехе
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorRequestsFrequencyDetail{
				CanReqAfter:        timestamppb.New(vErr.CanReqAfter()),
				CanReqAttemptsLeft: int32(vErr.CanReqAttemptsLeft()),
			})
		case std.ErrorUnprocessable:
			// todo : по логике это дубликат IsAsPassed случая cfm.ErrorFinished, но...
			logger.WarnFf(ctx, "%v обратился к завершенному SignIn %q: %#v", cl, signInId, vErr)
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorUnprocessableDetail{})
		case std.ErrorRuntime:
			logger.ErrorFf(ctx, err.Error())
			return nil, kernel.ErrorInternal()
		default:
			panic(err)
		}

		// --

		// can again
		if res.CanReqAgain() {
			return &pb.SignInRequestRetryResponse{
				SignInId:    res.SignInId().String(),
				RetriesLeft: int32(res.RetriesLeft()),
				CanRetryAt:  timestamppb.New(res.CanRetryAt()),
			}, nil
		}

		// last
		return &pb.SignInRequestRetryResponse{
			SignInId:    res.SignInId().String(),
			RetriesLeft: 0,
			CanRetryAt:  nil,
		}, nil
	}
}

// ---------------------------------------------------------------------------------------------------------------------
