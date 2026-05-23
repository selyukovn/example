package loggable

import (
	"context"
	"example/admin/gateway/internal/infra/cache"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ cache.CacheInterface = Decorator{}

type Decorator struct {
	origin   cache.CacheInterface
	maskKeys bool
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewDecorator(origin cache.CacheInterface, maskKeys bool) Decorator {
	return Decorator{
		origin:   origin,
		maskKeys: maskKeys,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (d Decorator) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	keyForLog := std.Ternary(d.maskKeys, std.MaskStrNotFirstLast(key), key)
	logger.InfoFf(ctx, "%T.%s(%q); ttl=%s", d, "Set", keyForLog, ttl.String())

	err := d.origin.Set(ctx, key, value, ttl)

	if err != nil {
		keyForLog := std.Ternary(d.maskKeys, std.MaskStrNotFirstLast(key), key)
		logger.ErrorFf(ctx, "%T.%s(%q)- ERROR: %#v = %s", d, "Set", keyForLog, err, err)
	}

	return err
}

func (d Decorator) Unset(ctx context.Context, key string) error {
	keyForLog := std.Ternary(d.maskKeys, std.MaskStrNotFirstLast(key), key)
	logger.InfoFf(ctx, "%T.%s(%q)", d, "Unset", keyForLog)

	err := d.origin.Unset(ctx, key)

	if err != nil {
		keyForLog := std.Ternary(d.maskKeys, std.MaskStrNotFirstLast(key), key)
		logger.ErrorFf(ctx, "%T.%s(%q)- ERROR: %#v = %s", d, "Unset", keyForLog, err, err)
	}

	return err
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

// Get
//
// Ошибки:
//   - std.ErrorNotFound
//   - std.ErrorRuntime
func (d Decorator) Get(ctx context.Context, key string) ([]byte, error) {
	keyForLog := std.Ternary(d.maskKeys, std.MaskStrNotFirstLast(key), key)

	rRes, rErr := d.origin.Get(ctx, key)
	switch rErr.(type) {
	case nil:
		logger.InfoFf(ctx, "%T.%s(%q)- hit", d, "Get", keyForLog)
	case std.ErrorNotFound:
		logger.InfoFf(ctx, "%T.%s(%q)- miss", d, "Get", keyForLog)
	case std.ErrorRuntime:
		logger.ErrorFf(ctx, "%T.%s(%q)- ERROR: %#v = %s", d, "Get", keyForLog, rErr, rErr)
	default:
		panic(rErr)
	}

	return rRes, rErr
}

// ---------------------------------------------------------------------------------------------------------------------
