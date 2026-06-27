package kafapi

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	assert "github.com/selyukovn/go-wm-assert"
)

// Handle
//
// Парсит kafka.Message и вызывает соответствующий типу/версии сообщения метод обработки из ServiceInterface.
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - ErrorDecoding
//   - ErrorMapping
//   - ErrorUnsupported
//   - ErrorHandling
func Handle(service ServiceInterface, ctx context.Context, kMsg *kafka.Message) error {
	assert.NotNilDeepMust(service)
	assert.NotNilDeepMust(ctx)
	assert.NotNilDeepMust(kMsg)

	meta, data, err := parse(kMsg)
	if err != nil {
		return err
	}

	err = handle(service, ctx, meta, data)
	if err != nil {
		return err
	}

	return nil
}
