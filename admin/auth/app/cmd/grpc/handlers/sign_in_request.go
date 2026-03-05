package handlers

import (
	"context"
	"example/admin/auth/cmd/grpc/container"
	"example/admin/auth/cmd/grpc/helpers"
	"example/admin/auth/cmd/grpc/pb"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/opera/use_cases/sign_in_request"
	"github.com/selyukovn/go-std"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewSignInRequest(ctr *container.Container) func(ctx context.Context, req *pb.SignInRequestRequest) (*pb.SignInRequestResponse, error) {
	return func(ctx context.Context, req *pb.SignInRequestRequest) (*pb.SignInRequestResponse, error) {
		cl, err := helpers.ParseClient(req.FromIp, req.FromUserAgent)
		if err != nil {
			ctr.Logger.CtxDebugFf(ctx, err.Error())
			return nil, helpers.ErrorInvalidArgument("кривой client")
		}

		email, err := std.EmailFromString(req.Email)
		if err != nil {
			ctr.Logger.CtxDebugFf(ctx, err.Error())
			return nil, helpers.ErrorFailedPrecondition(&pb.ErrorValidationDetail{
				Field:   "Email",
				Message: err.Error(),
			})
		}

		// --

		res, err := ctr.UseCases.SignInRequest.Execute(sign_in_request.NewArgs(ctx, cl, email))
		switch err.(type) {
		case nil:
		case std.ErrorNotFound:
			return nil, helpers.ErrorNotFound()
		case account.ErrorDeactivated, account.ErrorIpWhitelist:
			return nil, helpers.ErrorFailedPrecondition(&pb.ErrorAccountAccessDeniedDetail{})
		case std.ErrorRuntime:
			ctr.Logger.CtxErrorFf(ctx, err.Error())
			return nil, helpers.ErrorInternal()
		default:
			panic(err)
		}

		// --

		return &pb.SignInRequestResponse{
			SignInId:    res.SignInId().String(),
			RetriesLeft: int32(res.RetriesLeft()),
			// это всегда первый запрос -- тут не бывает финальных попыток, поэтому CanReqAfter всегда не 0
			CanRetryAt: timestamppb.New(res.CanRetryAt()),
			ExpireAt:   timestamppb.New(res.ExpireAt()),
		}, nil
	}
}
