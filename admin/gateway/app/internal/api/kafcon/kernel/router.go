package kernel

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Router struct {
	m map[string]HandlerInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewRouter() Router {
	return Router{
		m: make(map[string]HandlerInterface),
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (r Router) Register(topic string, handler HandlerInterface) {
	r.m[topic] = handler
}

func (r Router) Handle(ctx context.Context, kMsg *kafka.Message) error {
	topic := *kMsg.TopicPartition.Topic

	handler, isExist := r.m[topic]
	if !isExist {
		panic(fmt.Errorf("%T не знает, как обрабатывать сообщения из топика %q", r, topic))
	}

	return handler.Handle(ctx, kMsg)
}

// ---------------------------------------------------------------------------------------------------------------------
