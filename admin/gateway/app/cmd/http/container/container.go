package container

import (
	"context"
	"example/admin/gateway/cmd/common/components/processing"
	adapt_infra_cache "example/admin/gateway/internal/adapt/infra/cache"
	adapt_infra_clients_auth "example/admin/gateway/internal/adapt/infra/clients/auth"
	adapt_infra_clients_auth_cacher "example/admin/gateway/internal/adapt/infra/clients/auth/cacher"
	infra_cache "example/admin/gateway/internal/infra/cache"
	infra_cache_redis "example/admin/gateway/internal/infra/cache/redis"
	infra_clients_auth "example/admin/gateway/internal/infra/clients/auth"
	infra_clients_auth_grpc "example/admin/gateway/internal/infra/clients/auth/grpc"
	"github.com/redis/go-redis/v9"
	assert "github.com/selyukovn/go-wm-assert"
)

type Container = struct {
	Services Services
}

type Services = struct {
	Auth infra_clients_auth.ClientInterface
}

func New(
	redisCacheClient *redis.Client,
	appCfmApiGrpcBaseUrl string,
	appCfmApiGrpcApiKey string,
) *Container {
	assert.Str().NotEmpty().Must(appCfmApiGrpcBaseUrl)
	assert.Str().NotEmpty().Must(appCfmApiGrpcApiKey)

	// -----------------------------------------------------------------------------------------------------------------

	// <auth-service>
	var sAuth infra_clients_auth.ClientInterface
	sAuth = infra_clients_auth_grpc.NewClientGrpcMust(
		appCfmApiGrpcBaseUrl,
		appCfmApiGrpcApiKey,
		func(ctx context.Context) string {
			return processing.OperationId(ctx)
		},
	)
	// +cache
	var authCache infra_cache.CacheInterface
	authCache = infra_cache_redis.New(redisCacheClient)
	authCache = adapt_infra_cache.NewDecoratorLoggable(authCache, true)
	sAuth = adapt_infra_clients_auth.NewDecoratorCacheable(
		sAuth,
		adapt_infra_clients_auth_cacher.NewCacher(authCache),
	)
	// +logger
	sAuth = adapt_infra_clients_auth.NewDecoratorLoggable(sAuth)
	// </auth-service>

	// -----------------------------------------------------------------------------------------------------------------

	// Контейнер -- структура потенциально "растущая" (будут добавляться новые сервисы и т.д.).
	// Поэтому лучше сразу использовать контейнер через указатель.
	return &Container{
		Services: Services{
			Auth: sAuth,
		},
	}
}
