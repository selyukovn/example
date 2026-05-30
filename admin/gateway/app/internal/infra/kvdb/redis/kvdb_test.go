package redis

import (
	"bytes"
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/selyukovn/go-std"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testMakeKvDb() (KvDb, *redis.Client) {
	r := redis.NewClient(&redis.Options{Addr: std.Must(miniredis.Run()).Addr()})
	return New(r), r
}

// ---------------------------------------------------------------------------------------------------------------------
// Insert
// ---------------------------------------------------------------------------------------------------------------------

func Test_Insert(t *testing.T) {
	ctx := context.TODO()
	key := "hello"

	t.Run("new", func(t *testing.T) {
		kvDb, r := testMakeKvDb()
		expVal := []byte("world")

		err := kvDb.Insert(ctx, key, expVal)

		actVal := std.Must(r.Get(ctx, key).Bytes())
		assert.NoError(t, err)
		assert.True(t, bytes.Equal(actVal, expVal))
	})

	t.Run("already-exists", func(t *testing.T) {
		kvDb, r := testMakeKvDb()
		val := []byte("world")
		r.Set(ctx, key, val, 0)

		err := kvDb.Insert(ctx, key, val)

		assert.IsType(t, std.ErrorAlreadyDone{}, err)
	})

	t.Run("error", func(t *testing.T) {
		kvDb, r := testMakeKvDb()
		_ = r.Close()

		err := kvDb.Insert(ctx, key, []byte("world"))

		assert.IsType(t, std.ErrorRuntime{}, err)
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------------------------------------------------

func Test_Update(t *testing.T) {
	ctx := context.TODO()
	key := "hello"

	t.Run("not existed", func(t *testing.T) {
		kvDb, _ := testMakeKvDb()

		err := kvDb.Update(ctx, key, []byte("world"))

		assert.IsType(t, std.ErrorNotFound{}, err)
	})

	t.Run("existed", func(t *testing.T) {
		kvDb, r := testMakeKvDb()
		expVal := []byte("world")
		r.Set(ctx, key, expVal, 0)

		err := kvDb.Update(ctx, key, expVal)

		actVal := std.Must(r.Get(ctx, key).Bytes())
		assert.NoError(t, err)
		assert.True(t, bytes.Equal(actVal, expVal))
	})

	t.Run("error", func(t *testing.T) {
		kvDb, r := testMakeKvDb()
		_ = r.Close()

		err := kvDb.Update(ctx, key, []byte("world"))

		assert.IsType(t, std.ErrorRuntime{}, err)
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------------------------------------------------

func Test_Delete(t *testing.T) {
	ctx := context.TODO()
	key := "hello"

	t.Run("not existed", func(t *testing.T) {
		kvDb, _ := testMakeKvDb()

		err := kvDb.Delete(ctx, key)

		assert.IsType(t, std.ErrorNotFound{}, err)
	})

	t.Run("existed", func(t *testing.T) {
		kvDb, r := testMakeKvDb()
		r.Set(ctx, key, []byte("world"), 0)

		err := kvDb.Delete(ctx, key)

		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		kvDb, r := testMakeKvDb()
		_ = r.Close()

		err := kvDb.Delete(ctx, key)

		assert.IsType(t, std.ErrorRuntime{}, err)
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// Get
// ---------------------------------------------------------------------------------------------------------------------

func Test_Get(t *testing.T) {
	ctx := context.TODO()

	t.Run("not existed", func(t *testing.T) {
		kvDb, _ := testMakeKvDb()
		key := "hello"

		val, err := kvDb.Get(ctx, key)

		assert.Nil(t, val)
		assert.IsType(t, std.ErrorNotFound{}, err)
	})

	t.Run("existed", func(t *testing.T) {
		kvDb, r := testMakeKvDb()
		tCases := map[string][]byte{
			"key1": []byte("value1"),
			"key2": []byte("value2"),
			"key3": []byte(""),
			"key4": {0},
			"key5": {},
			"key6": nil,
		}

		for key := range tCases {
			expVal := tCases[key]
			r.Set(ctx, key, expVal, 0)

			actVal, err := kvDb.Get(ctx, key)

			assert.NoError(t, err)
			assert.True(t, bytes.Equal(actVal, expVal))
		}
	})

	t.Run("error", func(t *testing.T) {
		kvDb, r := testMakeKvDb()
		key := "hello"
		_ = r.Close()

		val, err := kvDb.Get(ctx, key)

		assert.Nil(t, val)
		assert.IsType(t, std.ErrorRuntime{}, err)
	})
}

// ---------------------------------------------------------------------------------------------------------------------
