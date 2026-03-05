package handlers

import (
	"context"
	"example/admin/auth/cmd/grpc/container"
	"example/admin/auth/cmd/grpc/helpers"
	"example/admin/auth/cmd/grpc/pb"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/action_request"
	"example/admin/auth/internal/domain/cfm"
	"example/admin/auth/internal/opera/use_cases/sign_in_confirm"
	"github.com/selyukovn/go-std"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewSignInConfirm(ctr *container.Container) func(ctx context.Context, req *pb.SignInConfirmRequest) (*pb.SignInConfirmResponse, error) {
	return func(ctx context.Context, req *pb.SignInConfirmRequest) (*pb.SignInConfirmResponse, error) {
		cl, err := helpers.ParseClient(req.FromIp, req.FromUserAgent)
		if err != nil {
			ctr.Logger.CtxDebugFf(ctx, err.Error())
			return nil, helpers.ErrorInvalidArgument("кривой client")
		}

		signInId, err := action_request.IdFromString(req.SignInId)
		if err != nil {
			ctr.Logger.CtxDebugFf(ctx, err.Error())
			return nil, helpers.ErrorFailedPrecondition(&pb.ErrorValidationDetail{
				Field:   "SignInId",
				Message: err.Error(),
			})
		}
		code, err := cfm.CodeFromString(req.Code)
		if err != nil {
			ctr.Logger.CtxDebugFf(ctx, err.Error())
			return nil, helpers.ErrorFailedPrecondition(&pb.ErrorValidationDetail{
				Field:   "Code",
				Message: err.Error(),
			})
		}

		// --

		res, err := ctr.UseCases.SignInConfirm.Execute(sign_in_confirm.NewArgs(ctx, cl, signInId, code))
		switch vErr := err.(type) {
		case nil:
		case std.ErrorNotFound:
			return nil, helpers.ErrorNotFound()
		case account.ErrorDeactivated, account.ErrorIpWhitelist:
			return nil, helpers.ErrorFailedPrecondition(&pb.ErrorAccountAccessDeniedDetail{})
		case cfm.ErrorFinished:
			return nil, helpers.ErrorFailedPrecondition(&pb.ErrorSignInFinishedDetail{
				IsPassed:  vErr.IsAsPassed(),
				IsFailed:  vErr.IsAsFailed(),
				IsExpired: vErr.IsAsExpired(),
			})
		case std.ErrorUnprocessable:
			// todo : по логике это дубликат IsAsPassed случая cfm.ErrorFinished, но...
			ctr.Logger.CtxWarnFf(ctx, "%v обратился к завершенному SignIn %q: %#v", cl, signInId, vErr)
			return nil, helpers.ErrorFailedPrecondition(&pb.ErrorUnprocessableDetail{})
		case std.ErrorRuntime:
			ctr.Logger.CtxErrorFf(ctx, err.Error())
			return nil, helpers.ErrorInternal()
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
