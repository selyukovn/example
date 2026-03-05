package handlers

import (
	"context"
	"example/admin/cfm/cmd/common/container"
	"example/admin/cfm/cmd/grpc/helpers"
	"example/admin/cfm/cmd/grpc/pb"
	"example/admin/cfm/internal/opera/use_cases/create_for_email"
	"github.com/selyukovn/go-std"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewCreateForEmail(ctr *container.Container) func(ctx context.Context, req *pb.CreateForEmailRequest) (*pb.CreateForEmailResponse, error) {
	return func(ctx context.Context, req *pb.CreateForEmailRequest) (*pb.CreateForEmailResponse, error) {
		email, err := std.EmailFromString(req.Email)
		if err != nil {
			ctr.Logger.CtxDebugFf(ctx, err.Error())
			return nil, helpers.ErrorInvalidArgument("кривой email")
		}

		// --

		res, err := ctr.UseCases.CreateForEmail.Execute(create_for_email.NewArgs(ctx, email))
		switch err.(type) {
		case nil:
		case std.ErrorRuntime:
			ctr.Logger.CtxErrorFf(ctx, err.Error())
			return nil, helpers.ErrorInternal()
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
