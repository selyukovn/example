package kafcon

import (
	"context"
	"example/admin/gateway/cmd/kafcon/bundles/admin_auth_events"
	"example/admin/gateway/cmd/kafcon/container"
	"example/admin/gateway/cmd/kafcon/kernel"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// ---------------------------------------------------------------------------------------------------------------------
// Const
// ---------------------------------------------------------------------------------------------------------------------

const (
	AdminAuthEventsTopic = admin_auth_events.TopicName
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type topicRouter struct {
	m map[string]kernel.TopicHandlerInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func newRouter(ctr *container.Container) topicRouter {
	m := make(map[string]kernel.TopicHandlerInterface)

	admin_auth_events.Register(m, ctr)

	return topicRouter{
		m: m,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (r topicRouter) handle(ctx context.Context, kMsg *kafka.Message) error {
	topicName := *kMsg.TopicPartition.Topic

	topicHandler, isExist := r.m[topicName]
	if !isExist {
		panic(fmt.Errorf("%T не знает, как обрабатывать сообщения из топика %q", r, topicName))
	}

	return topicHandler.Handle(ctx, kMsg)
}

// ---------------------------------------------------------------------------------------------------------------------
