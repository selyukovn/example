package cfm

import (
	"example/admin/cfm/internal/api/grpc/handlers/cfm/actions"
	"example/admin/cfm/internal/api/grpc/handlers/cfm/pb"
	"example/admin/cfm/internal/opera/use_cases/confirm"
	"example/admin/cfm/internal/opera/use_cases/create_for_email"
	"example/admin/cfm/internal/opera/use_cases/request"
	assert "github.com/selyukovn/go-wm-assert"
	"google.golang.org/grpc"
)

func Register(
	grpcServer *grpc.Server,
	ucCreateForEmail create_for_email.Command,
	ucRequest request.Command,
	ucConfirm confirm.Command,
	sMws []func(pb.CfmServiceServer) pb.CfmServiceServer,
) {
	assert.NotNilDeepMust(grpcServer)

	var service pb.CfmServiceServer = newServiceDefault(
		actions.NewCreateForEmail(ucCreateForEmail),
		actions.NewRequest(ucRequest),
		actions.NewConfirm(ucConfirm),
	)

	for i := len(sMws) - 1; i >= 0; i-- {
		service = sMws[i](service)
	}

	pb.RegisterCfmServiceServer(grpcServer, service)
}
