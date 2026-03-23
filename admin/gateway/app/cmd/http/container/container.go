package container

import (
	"context"
	"example/admin/gateway/cmd/http/components/processing"
	adapt_infra_clients_auth "example/admin/gateway/internal/adapt/infra/clients/auth"
	infra_clients_auth "example/admin/gateway/internal/infra/clients/auth"
	infra_clients_auth_grpc "example/admin/gateway/internal/infra/clients/auth/grpc"
	assert "github.com/selyukovn/go-wm-assert"
)

type Container = struct {
	Services Services
}

type Services = struct {
	Auth infra_clients_auth.ClientInterface
}

func New(
	appCfmApiGrpcBaseUrl string,
	appCfmApiGrpcApiKey string,
) *Container {
	assert.Str().NotEmpty().Must(appCfmApiGrpcBaseUrl)
	assert.Str().NotEmpty().Must(appCfmApiGrpcApiKey)

	// -----------------------------------------------------------------------------------------------------------------

	// auth
	var sAuth infra_clients_auth.ClientInterface
	sAuth = infra_clients_auth_grpc.NewClientGrpcMust(
		appCfmApiGrpcBaseUrl,
		appCfmApiGrpcApiKey,
		func(ctx context.Context) string {
			return processing.OperationId(ctx)
		},
	)
	sAuth = adapt_infra_clients_auth.NewDecoratorLoggable(sAuth)

	// -----------------------------------------------------------------------------------------------------------------

	// Контейнер -- структура потенциально "растущая" (будут добавляться новые сервисы и т.д.).
	// Поэтому лучше сразу использовать контейнер через указатель.
	return &Container{
		Services: Services{
			Auth: sAuth,
		},
	}
}
