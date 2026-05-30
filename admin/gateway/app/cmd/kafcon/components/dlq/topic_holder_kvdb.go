package dlq

import (
	"context"
	"errors"
	"example/admin/gateway/internal/infra/kvdb"
	"fmt"
	"github.com/avast/retry-go/v5"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ TopicHolderInterface = TopicHolderKvDb{}

type TopicHolderKvDb struct {
	kvDb kvdb.KvDbInterface
}

func topicHolderKvDbMakeKey(topic string) string {
	return fmt.Sprintf("DlqTopicHolder/%s", topic)
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewTopicHolderKvDb
//
// Паникует при нулевых аргументах.
func NewTopicHolderKvDb(kvDb kvdb.KvDbInterface) TopicHolderKvDb {
	assert.NotNilDeepMust(kvDb)

	return TopicHolderKvDb{
		kvDb: kvDb,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Hold
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorAlreadyDone
//   - std.ErrorRuntime
func (l TopicHolderKvDb) Hold(ctx context.Context, topic string) error {
	assert.NotNilDeepMust(ctx)
	assert.Str().NotEmpty().Must(topic)

	_m_ := "Hold"

	key := topicHolderKvDbMakeKey(topic)

	err := l.kvDb.Insert(ctx, key, []byte{})
	switch err.(type) {
	case nil:
		return nil
	case std.ErrorAlreadyDone:
		return err
	case std.ErrorRuntime:
		return std.WrapErrorToRuntime(err, l, _m_)
	default:
		panic(err)
	}
}

// UnHold
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorAlreadyDone
//   - std.ErrorRuntime
func (l TopicHolderKvDb) UnHold(ctx context.Context, topic string) error {
	assert.NotNilDeepMust(ctx)
	assert.Str().NotEmpty().Must(topic)

	_m_ := "UnHold"

	key := topicHolderKvDbMakeKey(topic)

	err := l.kvDb.Delete(ctx, key)
	switch err.(type) {
	case nil:
		return nil
	case std.ErrorNotFound:
		return std.NewErrorAlreadyDoneFf("топик не удерживается")
	case std.ErrorRuntime:
		return std.WrapErrorToRuntime(err, l, _m_)
	default:
		panic(err)
	}
}

// WaitTillOnHold
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (l TopicHolderKvDb) WaitTillOnHold(ctx context.Context, topic string) error {
	assert.NotNilDeepMust(ctx)
	assert.Str().NotEmpty().Must(topic)

	_m_ := "WaitTillOnHold"

	key := topicHolderKvDbMakeKey(topic)

	if err := retry.New(
		retry.Context(ctx),
		retry.Attempts(0 /* 0 -- бесконечно */),
		retry.Delay(time.Second /* 1 секунда -- примерное значение */),
		retry.DelayType(retry.FixedDelay),
	).Do(func() error {
		_, err := l.kvDb.Get(ctx, key)
		switch err.(type) {
		case nil:
			return errors.ErrUnsupported // любая ошибка для запуска повтора
		case std.ErrorNotFound:
			return nil
		case std.ErrorRuntime:
			return retry.Unrecoverable(err)
		default:
			panic(err)
		}
	}); err != nil {
		return std.WrapErrorToRuntime(err, l, _m_)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
