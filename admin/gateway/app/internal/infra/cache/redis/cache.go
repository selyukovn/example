package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Cache struct {
	r *redis.Client
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func New(client *redis.Client) Cache {
	assert.NotNilDeepMust(client)

	return Cache{
		r: client,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Set
//
// Ошибки:
//   - std.ErrorRuntime
func (c Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	err := c.r.Set(ctx, key, value, ttl).Err()

	if err != nil {
		return std.WrapErrorToRuntime(err, c, "Set")
	}

	return nil
}

// Unset
//
// Ошибки:
//   - std.ErrorRuntime
func (c Cache) Unset(ctx context.Context, key string) error {
	err := c.r.Del(ctx, key).Err()

	if err != nil {
		return std.WrapErrorToRuntime(err, c, "Unset")
	}

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
func (c Cache) Get(ctx context.Context, key string) ([]byte, error) {
	v, err := c.r.Get(ctx, key).Bytes()

	if errors.Is(err, redis.Nil) {
		return nil, std.NewErrorNotFoundFf("Key does not exist")
	} else if err != nil {
		return nil, std.WrapErrorToRuntime(err, c, "Get")
	}

	return v, nil
}

// ---------------------------------------------------------------------------------------------------------------------
