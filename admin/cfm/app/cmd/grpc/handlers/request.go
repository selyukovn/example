package handlers

import (
	"context"
	"example/admin/cfm/cmd/common/container"
	"example/admin/cfm/cmd/grpc/helpers"
	"example/admin/cfm/cmd/grpc/pb"
	"example/admin/cfm/internal/domain/cfm"
	"example/admin/cfm/internal/opera/use_cases/request"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewRequest(ctr *container.Container) func(ctx context.Context, req *pb.RequestRequest) (*pb.RequestResponse, error) {
	return func(ctx context.Context, req *pb.RequestRequest) (*pb.RequestResponse, error) {
		cfmId, err := cfm.IdFromString(req.CfmId)
		if err != nil {
			logger.DebugFf(ctx, err.Error())
			return nil, helpers.ErrorInvalidArgument("кривой id")
		}

		// --

		res, err := ctr.UseCases.Request.Execute(request.NewArgs(ctx, cfmId))
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
		case cfm.ErrorNoAttemptsLeft:
			return nil, helpers.ErrorFailedPrecondition(&pb.ErrorNoAttemptsLeftDetail{
				CfmId: vErr.CfmId().String(),
			})
		case cfm.ErrorRequestsFrequency:
			return nil, helpers.ErrorFailedPrecondition(&pb.ErrorRequestsFrequencyDetail{
				CfmId:              vErr.CfmId().String(),
				CanReqAttemptsLeft: uint32(vErr.CanReqAttemptsLeft()),
				CanReqAfter:        timestamppb.New(vErr.CanReqAfter()),
			})
		case std.ErrorRuntime:
			logger.ErrorFf(ctx, err.Error())
			return nil, helpers.ErrorInternal()
		default:
			panic(err)
		}

		// can again
		if res.CanReqAgain() {
			return &pb.RequestResponse{
				CfmId:              res.CfmId().String(),
				CanReqAgain:        true,
				CanReqAttemptsLeft: uint32(res.CanReqAttemptsLeft()),
				CanReqAfter:        timestamppb.New(res.CanReqAfter()),
			}, nil
		}

		// last
		return &pb.RequestResponse{
			CfmId:              res.CfmId().String(),
			CanReqAgain:        false,
			CanReqAttemptsLeft: 0,
			CanReqAfter:        nil,
		}, nil
	}
}
