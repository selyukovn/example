package memory

import (
	"github.com/selyukovn/go-std"
	"time"
)

// Ticker
// ---------------------------------------------------------------------------------------------------------------------

const tickerPeriod = 1 * time.Minute

func calcTickerPointPrev(t time.Time) unixSeconds {
	return unixSeconds(t.
		Add(-tickerPeriod).
		Add(-time.Duration(t.Nanosecond()) * time.Nanosecond).
		Add(-time.Duration(t.Second()) * time.Second).
		Unix(),
	)
}

func calcTickerPointNext(t time.Time) unixSeconds {
	return unixSeconds(t.
		Add(tickerPeriod).
		Add(-time.Duration(t.Nanosecond()) * time.Nanosecond).
		Add(-time.Duration(t.Second()) * time.Second).
		Unix(),
	)
}

// Cache
// ---------------------------------------------------------------------------------------------------------------------

// time.Time не годится в качестве ключа, но time.Time.Unix() -- самое оно.
type unixSeconds int64

type cacheCore struct {
	m                 *std.SyncMap[string, cacheItem]
	tickerPointToKeys *std.SyncMap[unixSeconds, map[string]struct{}]
}

type cacheItem struct {
	data  []byte
	expAt time.Time
}

// ---------------------------------------------------------------------------------------------------------------------
