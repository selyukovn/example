package handlers

import (
	"context"
	"example/admin/auth/internal/api/grpc/kernel"
	"example/admin/auth/internal/api/grpc/kernel_extra"
	"example/admin/auth/internal/api/grpc/pb"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/opera/use_cases/sign_in_request"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewSignInRequest(ucSignInRequest sign_in_request.Command) func(ctx context.Context, req *pb.SignInRequestRequest) (*pb.SignInRequestResponse, error) {
	return func(ctx context.Context, req *pb.SignInRequestRequest) (*pb.SignInRequestResponse, error) {
		cl, err := kernel_extra.ParseClient(req.FromIp, req.FromUserAgent)
		if err != nil {
			logger.DebugFf(ctx, err.Error())
			return nil, kernel.ErrorInvalidArgument("кривой client")
		}

		email, err := std.EmailFromString(req.Email)
		if err != nil {
			logger.DebugFf(ctx, err.Error())
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorValidationDetail{
				Field:   "Email",
				Message: err.Error(),
			})
		}

		// --

		res, err := ucSignInRequest.Execute(sign_in_request.NewArgs(ctx, cl, email))
		switch err.(type) {
		case nil:
		case std.ErrorNotFound:
			return nil, kernel.ErrorNotFound()
		case account.ErrorDeactivated, account.ErrorIpWhitelist:
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorAccountAccessDeniedDetail{})
		case std.ErrorRuntime:
			logger.ErrorFf(ctx, err.Error())
			return nil, kernel.ErrorInternal()
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
