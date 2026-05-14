package kernel

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// ---------------------------------------------------------------------------------------------------------------------

type TopicHandlerInterface interface {
	Handle(context.Context, *kafka.Message) error
}

// ---------------------------------------------------------------------------------------------------------------------

type FnTopicHandler func(context.Context, *kafka.Message) error

func (f FnTopicHandler) Handle(ctx context.Context, kMsg *kafka.Message) error {
	return f(ctx, kMsg)
}

// ---------------------------------------------------------------------------------------------------------------------
