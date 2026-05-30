package container

import (
	"database/sql"
	"example/admin/gateway/cmd/kafcon/components/dlq"
	adapt_infra_cache_loggable "example/admin/gateway/internal/adapt/infra/cache/loggable"
	adapt_infra_clients_auth_cachable "example/admin/gateway/internal/adapt/infra/clients/auth/cachable"
	infra_cache "example/admin/gateway/internal/infra/cache"
	infra_cache_redis "example/admin/gateway/internal/infra/cache/redis"
	infra_kvdb "example/admin/gateway/internal/infra/kvdb"
	infra_kvdb_redis "example/admin/gateway/internal/infra/kvdb/redis"
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
	GroupTracker dlq.GroupTrackerInterface
	Storage      dlq.StorageInterface
	TopicHolder  dlq.TopicHolderInterface
}

func New(
	redisCacheClient *redis.Client,
	redisKvDbClient *redis.Client,
	sqlDb *sql.DB,
	sqlDbFnIsDeadlockError func(error) bool,
	sqlDbFnIsDuplicateKeyError func(error) bool,
) *Container {
	assert.NotNilDeepMust(redisCacheClient)
	assert.NotNilDeepMust(redisKvDbClient)
	assert.NotNilDeepMust(sqlDb)
	assert.NotNilDeepMust(sqlDbFnIsDeadlockError)
	assert.NotNilDeepMust(sqlDbFnIsDuplicateKeyError)

	var cache infra_cache.CacheInterface
	cache = infra_cache_redis.New(redisCacheClient)
	cache = adapt_infra_cache_loggable.NewDecorator(cache, true)

	authServiceCacher := adapt_infra_clients_auth_cachable.NewCacher(cache)

	var kvDb infra_kvdb.KvDbInterface
	kvDb = infra_kvdb_redis.New(redisKvDbClient)

	sqlTxr := txr.NewTxrImplSql(sqlDb, 2, 50*time.Millisecond, sqlDbFnIsDeadlockError)

	groupTracker := dlq.NewGroupTrackerKvDb(kvDb)
	dlqStorage := dlq.NewStorageSQL(sqlTxr)
	dlqTopicHolder := dlq.NewTopicHolderKvDb(kvDb)

	return &Container{
		AuthServiceCacher: authServiceCacher,
		Dlq: Dlq{
			GroupTracker: groupTracker,
			Storage:      dlqStorage,
			TopicHolder:  dlqTopicHolder,
		},
	}
}
