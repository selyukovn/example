package container

import (
	adapt_infra_clients_auth "example/admin/gateway/internal/adapt/infra/clients/auth"
	infra_clients_auth "example/admin/gateway/internal/infra/clients/auth"
	infra_clients_auth_grpc "example/admin/gateway/internal/infra/clients/auth/grpc"
	assert "github.com/selyukovn/go-wm-assert"
)

type Container struct {
	Services Services
}

type Services = struct {
	Auth infra_clients_auth.ClientInterface
}

func New(
	appCfmApiGrpcBaseUrl string,
	appCfmApiGrpcApiKey string,
) *Container {
	assert.NotNilDeepMust(appCfmApiGrpcBaseUrl)
	assert.NotNilDeepMust(appCfmApiGrpcApiKey)

	// -----------------------------------------------------------------------------------------------------------------

	// auth
	var sAuth infra_clients_auth.ClientInterface
	sAuth = infra_clients_auth_grpc.NewClientGrpcMust(appCfmApiGrpcBaseUrl, appCfmApiGrpcApiKey)
	sAuth = adapt_infra_clients_auth.NewDecoratorLoggable(sAuth)

	// -----------------------------------------------------------------------------------------------------------------

	return &Container{
		Services: Services{
			Auth: sAuth,
		},
	}
}
