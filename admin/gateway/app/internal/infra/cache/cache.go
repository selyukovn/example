package cache

import (
	"context"
	"time"
)

type CacheInterface interface {
	// Set
	//
	// Ошибки:
	// 	- std.ErrorRuntime
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Unset
	//
	// Ошибки:
	// 	- std.ErrorRuntime
	Unset(ctx context.Context, key string) error

	// Get
	//
	// Ошибки:
	// 	- std.ErrorNotFound
	// 	- std.ErrorRuntime
	Get(ctx context.Context, key string) ([]byte, error)
}
