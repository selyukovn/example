package container

import (
	"context"
	"example/admin/gateway/cmd/http/components/processing"
	adapt_infra_cache "example/admin/gateway/internal/adapt/infra/cache"
	adapt_infra_clients_auth "example/admin/gateway/internal/adapt/infra/clients/auth"
	adapt_infra_clients_auth_cacher "example/admin/gateway/internal/adapt/infra/clients/auth/cacher"
	infra_cache "example/admin/gateway/internal/infra/cache"
	infra_cache_memory "example/admin/gateway/internal/infra/cache/memory"
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
) (
	// TODO : РЕФАКТОРИНГ : мемори-кеш требует запуска тикера, который нужно уметь аккуратно останавливать,
	// 		но это никак не должно быть связано с контейнером.
	*Container,
	infra_cache_memory.Ticker,
) {
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
	var authCache infra_cache.CacheInterface
	authCache, authCacheTicker := infra_cache_memory.New()
	authCache = adapt_infra_cache.NewDecoratorLoggable(authCache, true)
	sAuth = adapt_infra_clients_auth.NewDecoratorCacheable(
		sAuth,
		adapt_infra_clients_auth_cacher.NewCacher(authCache),
	)
	sAuth = adapt_infra_clients_auth.NewDecoratorLoggable(sAuth)

	// -----------------------------------------------------------------------------------------------------------------

	// Контейнер -- структура потенциально "растущая" (будут добавляться новые сервисы и т.д.).
	// Поэтому лучше сразу использовать контейнер через указатель.
	return &Container{
		Services: Services{
			Auth: sAuth,
		},
	}, authCacheTicker
}
