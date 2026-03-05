package container

import (
	adapt_infra_clients_auth "example/admin/bff/internal/adapt/infra/clients/auth"
	infra_clients_auth "example/admin/bff/internal/infra/clients/auth"
	infra_clients_auth_grpc "example/admin/bff/internal/infra/clients/auth/grpc"
	infra_logger "example/admin/bff/internal/infra/logger"
	assert "github.com/selyukovn/go-wm-assert"
	"io"
)

type Container struct {
	Logger   *infra_logger.Logger
	Services Services
}

type Services = struct {
	Auth infra_clients_auth.ClientInterface
}

func New(
	logIo io.Writer,
	isDebug bool,
	appCfmApiGrpcBaseUrl string,
	appCfmApiGrpcApiKey string,
) *Container {
	assert.NotNilDeepMust(logIo)
	assert.NotNilDeepMust(appCfmApiGrpcBaseUrl)
	assert.NotNilDeepMust(appCfmApiGrpcApiKey)

	// -----------------------------------------------------------------------------------------------------------------

	// logger
	infraLogger := infra_logger.NewLogger(logIo, isDebug)

	// -----------------------------------------------------------------------------------------------------------------

	// auth
	var sAuth infra_clients_auth.ClientInterface
	sAuth = infra_clients_auth_grpc.NewClientGrpcMust(appCfmApiGrpcBaseUrl, appCfmApiGrpcApiKey)
	sAuth = adapt_infra_clients_auth.NewDecoratorLoggable(sAuth, infraLogger)

	// -----------------------------------------------------------------------------------------------------------------

	return &Container{
		Logger: infraLogger,
		Services: Services{
			Auth: sAuth,
		},
	}
}
