package redis

import (
	"context"
	"errors"
	"example/admin/gateway/internal/infra/kvdb"
	"github.com/redis/go-redis/v9"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

const infTtl = 0

var _ kvdb.KvDbInterface = KvDb{}

type KvDb struct {
	r *redis.Client
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func New(client *redis.Client) KvDb {
	assert.NotNilDeepMust(client)

	return KvDb{
		r: client,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Insert
//
// Ошибки:
//   - std.ErrorAlreadyDone
//   - std.ErrorRuntime
func (k KvDb) Insert(ctx context.Context, key string, value []byte) error {
	_, err := k.r.SetArgs(ctx, key, value, redis.SetArgs{
		TTL:  infTtl,
		Mode: string(redis.NX),
	}).Result()

	if errors.Is(err, redis.Nil) {
		return std.NewErrorAlreadyDoneFf("key already exists")
	} else if err != nil {
		return std.WrapErrorToRuntime(err, k, "Insert")
	}

	return nil
}

// Update
//
// Ошибки:
//   - std.ErrorNotFound
//   - std.ErrorRuntime
func (k KvDb) Update(ctx context.Context, key string, value []byte) error {
	isExist, err := k.r.SetXX(ctx, key, value, infTtl).Result()

	if err != nil {
		return std.WrapErrorToRuntime(err, k, "Update")
	} else if !isExist {
		return std.NewErrorNotFoundFf("key not found")
	}

	return nil
}

// Delete
//
// Ошибки:
//   - std.ErrorNotFound
//   - std.ErrorRuntime
func (k KvDb) Delete(ctx context.Context, key string) error {
	// DEL не работает с шаблонами ключей, поэтому ничего экранировать не надо.
	n, err := k.r.Del(ctx, key).Result()

	if err != nil {
		return std.WrapErrorToRuntime(err, k, "Delete")
	} else if n == 0 {
		return std.NewErrorNotFoundFf("key not found")
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
func (k KvDb) Get(ctx context.Context, key string) ([]byte, error) {
	v, err := k.r.Get(ctx, key).Bytes()

	if errors.Is(err, redis.Nil) {
		return nil, std.NewErrorNotFoundFf("key not found")
	} else if err != nil {
		return nil, std.WrapErrorToRuntime(err, k, "Get")
	}

	return v, nil
}

// ---------------------------------------------------------------------------------------------------------------------
