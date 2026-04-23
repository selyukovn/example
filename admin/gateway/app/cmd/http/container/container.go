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
	"fmt"
	goroutiner "github.com/selyukovn/go-routiner"
	"github.com/selyukovn/go-std/logger"
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
	authCache, authCacheTicker := infra_cache_memory.New()
	// Точка запуска тикера может быть любой, поскольку ему не требуется graceful-shutdown.
	goroutiner.New(goroutiner.MwPanicToError(func(pv any, ds []byte, ctx context.Context) error {
		logger.PanicFf(ctx, pv, ds, "goroutiner.authCacheTicker")
		return fmt.Errorf("panic: %#v; stack: %s", pv, ds)
	})).SingleAsync(context.Background(), func(ctx context.Context) error {
		authCacheTicker.Start()
		return nil
	})
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
