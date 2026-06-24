package handlers

import (
	"context"
	"example/admin/cfm/internal/api/grpc/kernel"
	"example/admin/cfm/internal/api/grpc/pb"
	"example/admin/cfm/internal/domain/cfm"
	"example/admin/cfm/internal/opera/use_cases/request"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewRequest(ucRequest request.Command) func(ctx context.Context, req *pb.RequestRequest) (*pb.RequestResponse, error) {
	return func(ctx context.Context, req *pb.RequestRequest) (*pb.RequestResponse, error) {
		cfmId, err := cfm.IdFromString(req.CfmId)
		if err != nil {
			logger.DebugFf(ctx, err.Error())
			return nil, kernel.ErrorInvalidArgument("кривой id")
		}

		// --

		res, err := ucRequest.Execute(ctx, cfmId)
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
		case cfm.ErrorNoAttemptsLeft:
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorNoAttemptsLeftDetail{
				CfmId: vErr.CfmId().String(),
			})
		case cfm.ErrorRequestsFrequency:
			return nil, kernel.ErrorFailedPrecondition(&pb.ErrorRequestsFrequencyDetail{
				CfmId:              vErr.CfmId().String(),
				CanReqAttemptsLeft: uint32(vErr.CanReqAttemptsLeft()),
				CanReqAfter:        timestamppb.New(vErr.CanReqAfter()),
			})
		case std.ErrorRuntime:
			logger.ErrorFf(ctx, err.Error())
			return nil, kernel.ErrorInternal()
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
