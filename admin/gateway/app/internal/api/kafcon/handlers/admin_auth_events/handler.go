package admin_auth_events

import (
	"context"
	"example/admin/gateway/internal/api/kafcon/handlers/admin_auth_events/kafapi"
	"example/admin/gateway/internal/api/kafcon/kernel"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
)

func newHandler(service kafapi.ServiceInterface) kernel.HandlerInterface {
	return kernel.FnHandler(func(ctx context.Context, kMsg *kafka.Message) error {
		_o_ := TopicName
		_m_ := "newHandler"

		err := kafapi.Handle(service, ctx, kMsg)
		switch vErr := err.(type) {
		case nil, kafapi.ErrorDecoding, kafapi.ErrorMapping:
			return vErr
		case kafapi.ErrorUnsupported:
			// Внимание!
			// Сообщение может быть неопознанным из-за неактуальности локального `kafapi`-пакета
			// (например, при добавлении нового типа сообщения или изменении версии известного сообщения),
			// что скорее свидетельствует об отсутствии необходимости в обработке такого сообщения данным сервисом.
			// Ошибка разработчика менее вероятна (при условии поддержки обратной совместимости сообщений),
			// поэтому итоговый сценарий -- игнорирование с предупреждением.
			logger.WarnFf(ctx, "Неопознанное сообщение (id=%d): %s", vErr.Meta.Id, vErr.Error())
			return nil
		case kafapi.ErrorHandling:
			return std.WrapErrorToRuntime(vErr, _o_, _m_, "Handle", "ErrorHandling")
		default:
			panic(vErr)
		}
	})
}
