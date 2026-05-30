package dlq

import (
	"context"
	"example/admin/gateway/internal/infra/kvdb"
	"fmt"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ GroupTrackerInterface = GroupTrackerKvDb{}

type GroupTrackerKvDb struct {
	kvDb kvdb.KvDbInterface
}

func groupTrackerKvDbMakeKey(topic string, groupId string) string {
	return fmt.Sprintf("DlqGroupTracker/%s/%s", topic, groupId)
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewGroupTrackerKvDb
//
// Паникует при нулевых аргументах.
func NewGroupTrackerKvDb(kvDb kvdb.KvDbInterface) GroupTrackerKvDb {
	assert.NotNilDeepMust(kvDb)

	return GroupTrackerKvDb{kvDb: kvDb}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// SetLastHandled
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (g GroupTrackerKvDb) SetLastHandled(ctx context.Context, topic string, groupId string, msgId string) error {
	assert.NotNilDeepMust(ctx)
	assert.Str().NotEmpty().Must(topic)
	assert.Str().NotEmpty().Must(groupId)
	assert.Str().NotEmpty().Must(msgId)

	_m_ := "SetLastHandled"

	key := groupTrackerKvDbMakeKey(topic, groupId)
	value := []byte(msgId)

	err := g.kvDb.Update(ctx, key, value)
	switch err.(type) {
	case nil:
	case std.ErrorNotFound:
		err = g.kvDb.Insert(ctx, key, value)
		switch err.(type) {
		case nil:
		case std.ErrorAlreadyDone:
			err = g.kvDb.Update(ctx, key, value)
			switch err.(type) {
			case nil:
			case std.ErrorNotFound:
				return std.WrapErrorToRuntime(err, g, _m_, "Update", "Insert", "Update", "ErrorNotFound")
			case std.ErrorRuntime:
				return std.WrapErrorToRuntime(err, g, _m_, "Update", "Insert", "Update", "ErrorRuntime")
			default:
				panic(err)
			}
		case std.ErrorRuntime:
			return std.WrapErrorToRuntime(err, g, _m_, "Update", "Insert", "ErrorRuntime")
		default:
			panic(err)
		}
	case std.ErrorRuntime:
		return std.WrapErrorToRuntime(err, g, _m_, "Update")
	default:
		panic(err)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
