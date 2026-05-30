package dlq

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/selyukovn/go-std"
)

func Decorate(
	storage StorageInterface,
	topic string,
	fnGetGroupId func(context.Context, *kafka.Message) (string, error),
	fnHandle func(context.Context, *kafka.Message) error,
	fnHandleErrIsDlq func(context.Context, error) bool,
) func(context.Context, *kafka.Message) error {
	return func(ctx context.Context, kMsg *kafka.Message) error {
		_o_ := "dlq"
		_m_ := "Decorate"

		isAnyPoisoned, err := storage.IsGroupPoisoned(ctx, topic, AnyGroup)
		if err != nil {
			return std.WrapErrorToRuntime(err, _o_, _m_, "IsGroupPoisoned", AnyGroup)
		} else if isAnyPoisoned {
			err = storage.Add(ctx, topic, AnyGroup, kMsg)
			if err != nil {
				return std.WrapErrorToRuntime(err, _o_, _m_, "IsGroupPoisoned", AnyGroup, "Add")
			}
			return nil
		}

		groupId, err := fnGetGroupId(ctx, kMsg)
		if err != nil {
			err := storage.Add(ctx, topic, AnyGroup, kMsg)
			if err != nil {
				return std.WrapErrorToRuntime(err, _o_, _m_, "fnGetGroupId", AnyGroup, "Add")
			}
			return nil
		}

		isGroupPoisoned, err := storage.IsGroupPoisoned(ctx, topic, groupId)
		if err != nil {
			return std.WrapErrorToRuntime(err, _o_, _m_, "IsGroupPoisoned", "Group")
		} else if isGroupPoisoned {
			err = storage.Add(ctx, topic, groupId, kMsg)
			if err != nil {
				return std.WrapErrorToRuntime(err, _o_, _m_, "IsGroupPoisoned", "Group", "Add")
			}
			return nil
		}

		// --

		handleErr := fnHandle(ctx, kMsg)

		if fnHandleErrIsDlq(ctx, handleErr) {
			err := storage.Add(ctx, topic, groupId, kMsg)
			if err != nil {
				return std.WrapErrorToRuntime(err, _o_, _m_, "fnHandleErrIsDlq", "Add")
			}
			return nil
		}

		return handleErr
	}
}
