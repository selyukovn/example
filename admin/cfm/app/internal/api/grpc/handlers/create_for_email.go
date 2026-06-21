package handlers

import (
	"context"
	"example/admin/cfm/internal/api/grpc/kernel"
	"example/admin/cfm/internal/api/grpc/pb"
	"example/admin/cfm/internal/opera/use_cases/create_for_email"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewCreateForEmail(ucCreateForEmail create_for_email.Command) func(ctx context.Context, req *pb.CreateForEmailRequest) (*pb.CreateForEmailResponse, error) {
	return func(ctx context.Context, req *pb.CreateForEmailRequest) (*pb.CreateForEmailResponse, error) {
		email, err := std.EmailFromString(req.Email)
		if err != nil {
			logger.DebugFf(ctx, err.Error())
			return nil, kernel.ErrorInvalidArgument("кривой email")
		}

		// --

		res, err := ucCreateForEmail.Execute(create_for_email.NewArgs(ctx, email))
		switch err.(type) {
		case nil:
		case std.ErrorRuntime:
			logger.ErrorFf(ctx, err.Error())
			return nil, kernel.ErrorInternal()
		default:
			panic(err)
		}

		return &pb.CreateForEmailResponse{
			Email:    res.Email().String(),
			CfmId:    res.CfmId().String(),
			ExpireAt: timestamppb.New(res.ExpireAt()),
		}, nil
	}
}
