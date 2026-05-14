package container

import (
	adapt_infra_cache "example/admin/gateway/internal/adapt/infra/cache"
	adapt_infra_clients_auth_cacher "example/admin/gateway/internal/adapt/infra/clients/auth/cacher"
	infra_cache "example/admin/gateway/internal/infra/cache"
	infra_cache_redis "example/admin/gateway/internal/infra/cache/redis"
	"github.com/redis/go-redis/v9"
)

type Container = struct {
	AuthServiceCacher adapt_infra_clients_auth_cacher.Cacher
}

func New(
	redisCacheClient *redis.Client,
) *Container {

	var cache infra_cache.CacheInterface
	cache = infra_cache_redis.New(redisCacheClient)
	cache = adapt_infra_cache.NewDecoratorLoggable(cache, true)

	authServiceCacher := adapt_infra_clients_auth_cacher.NewCacher(cache)

	return &Container{
		AuthServiceCacher: authServiceCacher,
	}
}
