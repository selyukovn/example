package handlers

import (
	"context"
	"example/admin/cfm/cmd/common/container"
	"example/admin/cfm/cmd/grpc/helpers"
	"example/admin/cfm/cmd/grpc/pb"
	"example/admin/cfm/internal/domain/cfm"
	"example/admin/cfm/internal/domain/cfm/code"
	"example/admin/cfm/internal/opera/use_cases/confirm"
	"github.com/selyukovn/go-std"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewConfirm(ctr *container.Container) func(ctx context.Context, req *pb.ConfirmRequest) (*pb.ConfirmResponse, error) {
	return func(ctx context.Context, req *pb.ConfirmRequest) (*pb.ConfirmResponse, error) {
		cfmId, err := cfm.IdFromString(req.CfmId)
		if err != nil {
			ctr.Logger.CtxDebugFf(ctx, err.Error())
			return nil, helpers.ErrorInvalidArgument("кривой id")
		}

		cCode, err := code.CodeFromString(req.Code)
		if err != nil {
			ctr.Logger.CtxDebugFf(ctx, err.Error())
			return nil, helpers.ErrorInvalidArgument("кривой code")
		}

		// --

		res, err := ctr.UseCases.Confirm.Execute(confirm.NewArgs(ctx, cfmId, cCode))
		switch vErr := err.(type) {
		case nil:
		case std.ErrorNotFound:
			return nil, helpers.ErrorNotFound()
		case cfm.ErrorFinished:
			return nil, helpers.ErrorFailedPrecondition(&pb.ErrorFinishedDetail{
				CfmId:       vErr.CfmId().String(),
				FinishedAt:  timestamppb.New(vErr.FinishedAt()),
				IsAsExpired: vErr.IsAsExpired(),
				IsAsFailed:  vErr.IsAsFailed(),
				IsAsPassed:  vErr.IsAsPassed(),
			})
		case std.ErrorUnprocessable:
			return nil, helpers.ErrorFailedPrecondition(&pb.ErrorNotRequestedDetail{
				CfmId: cfmId.String(),
			})
		case std.ErrorRuntime:
			ctr.Logger.CtxErrorFf(ctx, err.Error())
			return nil, helpers.ErrorInternal()
		default:
			panic(err)
		}

		// can again
		if !res.IsFinished() {
			return &pb.ConfirmResponse{
				CfmId:              res.CfmId().String(),
				FinishedAt:         nil,
				IsFinishedAsFailed: false,
				IsFinishedAsPassed: false,
				FailsLeft:          uint32(res.FailsLeft()),
			}, nil
		}

		return &pb.ConfirmResponse{
			CfmId:              res.CfmId().String(),
			FinishedAt:         timestamppb.New(res.FinishedAt()),
			IsFinishedAsFailed: res.IsFinishedAsFailed(),
			IsFinishedAsPassed: res.IsFinishedAsPassed(),
			FailsLeft:          0,
		}, nil
	}
}
