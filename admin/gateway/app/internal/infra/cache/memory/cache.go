package memory

import (
	"context"
	"example/admin/gateway/internal/infra/cache"
	"github.com/selyukovn/go-std"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ cache.CacheInterface = Cache{}

type Cache struct {
	core cacheCore
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func newCache(core cacheCore) Cache {
	return Cache{core: core}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Set
//
// Ошибки:
//   - std.ErrorRuntime
func (c Cache) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	expAt := time.Now().Add(ttl)

	tickerPoint := calcTickerPointNext(expAt)

	c.core.tickerPointToKeys.Modify(tickerPoint, func(v map[string]struct{}) map[string]struct{} {
		if v == nil {
			v = make(map[string]struct{})
		}
		v[key] = struct{}{}
		return v
	})

	c.core.m.Set(key, cacheItem{
		data:  value,
		expAt: expAt,
	})

	return nil
}

// Unset
//
// Ошибки:
//   - std.ErrorRuntime
func (c Cache) Unset(_ context.Context, key string) error {
	c.core.m.Delete(key)
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

// Get
//
// Ошибки:
//   - std.ErrorNotFound
//   - std.ErrorRuntime
func (c Cache) Get(_ context.Context, key string) ([]byte, error) {
	item, exists := c.core.m.Get(key)
	if !exists {
		return nil, std.NewErrorNotFoundFf("%T: ключ не найден", c)
	}

	if item.expAt.Before(time.Now()) {
		return nil, std.NewErrorNotFoundFf("%T: ключ не найден", c)
	}

	return item.data, nil
}

// ---------------------------------------------------------------------------------------------------------------------
