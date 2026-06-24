package handlers

import (
	"context"
	"example/admin/cfm/internal/api/grpc/kernel"
	"example/admin/cfm/internal/api/grpc/pb"
	"example/admin/cfm/internal/domain/cfm"
	"example/admin/cfm/internal/domain/cfm/code"
	"example/admin/cfm/internal/opera/use_cases/confirm"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewConfirm(ucConfirm confirm.Command) func(ctx context.Context, req *pb.ConfirmRequest) (*pb.ConfirmResponse, error) {
	return func(ctx context.Context, req *pb.ConfirmRequest) (*pb.ConfirmResponse, error) {
		cfmId, err := cfm.IdFromString(req.CfmId)
		if err != nil {
			logger.DebugFf(ctx, err.Error())
			return nil, kernel.ErrorInvalidArgument("кривой id")
		}

		cCode, err := code.CodeFromString(req.Code)
		if err != nil {
			logger.DebugFf(ctx, err.Error())
			return nil, kernel.ErrorInvalidArgument("кривой code")
		}

		// --

		res, err := ucConfirm.Execute(ctx, cfmId, cCode)
		switch vErr := err.(type) {
		case nil:
		case std.ErrorNotFound:
			return nil, kernel.ErrorNotFound()
		case cfm.ErrorFinished:
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorFinishedDetail{
				CfmId:       vErr.CfmId().String(),
				FinishedAt:  timestamppb.New(vErr.FinishedAt()),
				IsAsExpired: vErr.IsAsExpired(),
				IsAsFailed:  vErr.IsAsFailed(),
				IsAsPassed:  vErr.IsAsPassed(),
			})
		case std.ErrorUnprocessable:
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorNotRequestedDetail{
				CfmId: cfmId.String(),
			})
		case std.ErrorRuntime:
			logger.ErrorFf(ctx, err.Error())
			return nil, kernel.ErrorInternal()
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
