package actions

import (
	"context"
	"example/admin/auth/internal/api/grpc/handlers/auth/pb"
	"example/admin/auth/internal/api/grpc/kernel"
	"example/admin/auth/internal/api/grpc/kernel_extra"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/action_request"
	"example/admin/auth/internal/domain/cfm"
	"example/admin/auth/internal/opera/use_cases/sign_in_confirm"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ---------------------------------------------------------------------------------------------------------------------

type SignInConfirm = func(ctx context.Context, req *pb.SignInConfirmRequest) (*pb.SignInConfirmResponse, error)

// ---------------------------------------------------------------------------------------------------------------------

func NewSignInConfirm(ucSignInConfirm sign_in_confirm.Command) SignInConfirm {
	return func(ctx context.Context, req *pb.SignInConfirmRequest) (*pb.SignInConfirmResponse, error) {
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
		code, err := cfm.CodeFromString(req.Code)
		if err != nil {
			logger.DebugFf(ctx, err.Error())
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorValidationDetail{
				Field:   "Code",
				Message: err.Error(),
			})
		}

		// --

		res, err := ucSignInConfirm.Execute(ctx, cl, signInId, code)
		switch vErr := err.(type) {
		case nil:
		case std.ErrorNotFound:
			return nil, kernel.ErrorNotFound()
		case account.ErrorDeactivated, account.ErrorIpWhitelist:
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorAccountAccessDeniedDetail{})
		case cfm.ErrorFinished:
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorSignInFinishedDetail{
				IsPassed:  vErr.IsAsPassed(),
				IsFailed:  vErr.IsAsFailed(),
				IsExpired: vErr.IsAsExpired(),
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

		return &pb.SignInConfirmResponse{
			IsPassed:        res.IsPassed(),
			AttemptsLeft:    int32(res.AttemptsLeft()),
			SessionId:       res.SessId().String(),
			SessionExpireAt: timestamppb.New(res.SessExpAt()),
		}, nil
	}
}

// ---------------------------------------------------------------------------------------------------------------------
