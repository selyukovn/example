package container

import (
	adapt_infra_cache_loggable "example/admin/gateway/internal/adapt/infra/cache/loggable"
	adapt_infra_clients_auth_cachable "example/admin/gateway/internal/adapt/infra/clients/auth/cachable"
	infra_cache "example/admin/gateway/internal/infra/cache"
	infra_cache_redis "example/admin/gateway/internal/infra/cache/redis"
	"github.com/redis/go-redis/v9"
)

type Container = struct {
	AuthServiceCacher adapt_infra_clients_auth_cachable.Cacher
}

func New(
	redisCacheClient *redis.Client,
) *Container {

	var cache infra_cache.CacheInterface
	cache = infra_cache_redis.New(redisCacheClient)
	cache = adapt_infra_cache_loggable.NewDecorator(cache, true)

	authServiceCacher := adapt_infra_clients_auth_cachable.NewCacher(cache)

	return &Container{
		AuthServiceCacher: authServiceCacher,
	}
}
