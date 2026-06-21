package kernel

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// ---------------------------------------------------------------------------------------------------------------------

type HandlerInterface interface {
	Handle(context.Context, *kafka.Message) error
}

// ---------------------------------------------------------------------------------------------------------------------

type FnHandler func(context.Context, *kafka.Message) error

func (f FnHandler) Handle(ctx context.Context, kMsg *kafka.Message) error {
	return f(ctx, kMsg)
}

// ---------------------------------------------------------------------------------------------------------------------
