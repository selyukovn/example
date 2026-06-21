package kafcon

import (
	"context"
	"example/admin/gateway/internal/api/kafcon/components/dlq"
	"example/admin/gateway/internal/api/kafcon/kernel"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type DlqProcessor struct {
	storage dlq.StorageInterface
	handler kernel.HandlerInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewDlqProcessor(storage dlq.StorageInterface, handler kernel.HandlerInterface) DlqProcessor {
	return DlqProcessor{
		storage: storage,
		handler: handler,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (d DlqProcessor) Process(ctx context.Context, topic string, groupId string) error {
	assert.NotNilDeepMust(ctx)
	assert.Str().NotEmpty().Must(groupId)

	_m_ := "Process"

	for item := range d.storage.GetMessages(ctx, topic, groupId, 0) {
		if item.Err != nil {
			return std.WrapErrorToRuntime(item.Err, d, _m_, "GetMessages")
		} else if item.KMsg == nil {
			return nil
		} else {
			logger.DebugFf(ctx, "Сообщение: %#v", item.KMsg)
		}

		err := d.handler.Handle(ctx, item.KMsg)
		if err != nil {
			return std.WrapErrorToRuntime(err, d, _m_, "Handle")
		} else {
			logger.DebugFf(ctx, "Сообщение обработано")
		}

		logger.DebugFf(ctx, "Удаление сообщения из DLQ...")
		err = d.storage.Remove(ctx, item.KMsg)
		if err != nil {
			return std.WrapErrorToRuntime(err, d, _m_, "Remove")
		} else {
			logger.DebugFf(ctx, "Cообщение удалено из DLQ")
		}
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
