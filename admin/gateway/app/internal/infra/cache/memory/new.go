package memory

import (
	"github.com/selyukovn/go-std"
)

func New() (Cache, Ticker) {
	core := cacheCore{
		m:                 std.NewSyncMap[string, cacheItem](),
		tickerPointToKeys: std.NewSyncMap[unixSeconds, map[string]struct{}](),
	}

	return newCache(core), newTicker(core)
}
