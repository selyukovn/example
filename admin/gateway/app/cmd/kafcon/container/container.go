package container

import (
	"database/sql"
	"example/admin/gateway/cmd/kafcon/components/dlq"
	adapt_infra_cache_loggable "example/admin/gateway/internal/adapt/infra/cache/loggable"
	adapt_infra_clients_auth_cachable "example/admin/gateway/internal/adapt/infra/clients/auth/cachable"
	infra_cache "example/admin/gateway/internal/infra/cache"
	infra_cache_redis "example/admin/gateway/internal/infra/cache/redis"
	"github.com/redis/go-redis/v9"
	"github.com/selyukovn/go-txr"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

type Container = struct {
	AuthServiceCacher adapt_infra_clients_auth_cachable.Cacher
	Dlq               Dlq
}

type Dlq = struct {
	Storage dlq.StorageInterface
}

func New(
	redisCacheClient *redis.Client,
	sqlDb *sql.DB,
	sqlDbFnIsDeadlockError func(error) bool,
	sqlDbFnIsDuplicateKeyError func(error) bool,
) *Container {
	assert.NotNilDeepMust(redisCacheClient)
	assert.NotNilDeepMust(sqlDb)
	assert.NotNilDeepMust(sqlDbFnIsDeadlockError)
	assert.NotNilDeepMust(sqlDbFnIsDuplicateKeyError)

	var cache infra_cache.CacheInterface
	cache = infra_cache_redis.New(redisCacheClient)
	cache = adapt_infra_cache_loggable.NewDecorator(cache, true)

	authServiceCacher := adapt_infra_clients_auth_cachable.NewCacher(cache)

	sqlTxr := txr.NewTxrImplSql(sqlDb, 2, 50*time.Millisecond, sqlDbFnIsDeadlockError)

	dlqStorage := dlq.NewStorageSQL(sqlTxr)

	return &Container{
		AuthServiceCacher: authServiceCacher,
		Dlq: Dlq{
			Storage: dlqStorage,
		},
	}
}
