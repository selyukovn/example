package kafcon

import (
	"context"
	"example/admin/gateway/cmd/kafcon/bundles/admin_auth_events"
	"example/admin/gateway/cmd/kafcon/container"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type DlqProcessor struct {
	ctr    *container.Container
	routes map[string]struct {
		FnGetGroupId func(context.Context, *kafka.Message) (string, error)
		FnHandle     func(context.Context, *kafka.Message) error
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewDlqProcessor(ctr *container.Container) DlqProcessor {
	assert.NotNilDeepMust(ctr)

	routes := map[string]struct {
		FnGetGroupId func(context.Context, *kafka.Message) (string, error)
		FnHandle     func(context.Context, *kafka.Message) error
	}{
		admin_auth_events.TopicName: {
			FnGetGroupId: admin_auth_events.FnGetGroupId,
			FnHandle:     admin_auth_events.NewTopicHandlerDefault(ctr).Handle,
		},
	}

	return DlqProcessor{
		ctr:    ctr,
		routes: routes,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (d DlqProcessor) Process(ctx context.Context, topic string, groupId string) error {
	assert.NotNilDeepMust(ctx)
	assert.Str().NotEmpty().Must(groupId)

	_m_ := "Process"

	route, isKnownTopic := d.routes[topic]
	if !isKnownTopic {
		panic(fmt.Errorf("%T не знает, как обрабатывать сообщения из топика %q", d, topic))
	}

	for item := range d.ctr.Dlq.Storage.GetMessages(ctx, topic, groupId, 0) {
		if item.Err != nil {
			return std.WrapErrorToRuntime(item.Err, d, _m_, "GetMessages")
		} else if item.KMsg == nil {
			return nil
		} else {
			logger.DebugFf(ctx, "Сообщение: %#v", item.KMsg)
		}

		err := route.FnHandle(ctx, item.KMsg)
		if err != nil {
			return std.WrapErrorToRuntime(err, d, _m_, "Handle")
		} else {
			logger.DebugFf(ctx, "Сообщение обработано")
		}

		logger.DebugFf(ctx, "Удаление сообщения из DLQ...")
		err = d.ctr.Dlq.Storage.Remove(ctx, item.KMsg)
		if err != nil {
			return std.WrapErrorToRuntime(err, d, _m_, "Remove")
		} else {
			logger.DebugFf(ctx, "Cообщение удалено из DLQ")
		}
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
