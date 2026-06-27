package dlq

import (
	"context"
	"example/admin/gateway/internal/api/kafcon/components/dlq"
	"example/admin/gateway/internal/api/kafcon/handlers/admin_auth_events"
	"example/admin/gateway/internal/api/kafcon/handlers/admin_auth_events/kafapi"
	"example/admin/gateway/internal/api/kafcon/kernel"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	assert "github.com/selyukovn/go-wm-assert"
)

func NewDecorator(handler kernel.HandlerInterface, dlqStorage dlq.StorageInterface) kernel.HandlerInterface {
	assert.NotNilDeepMust(handler)
	assert.NotNilDeepMust(dlqStorage)

	return kernel.FnHandler(dlq.Decorate(
		dlqStorage,
		admin_auth_events.TopicName,
		getGroupId,
		handler.Handle,
		isErrorDlq,
	))
}

func isErrorDlq(_ context.Context, err error) bool {
	switch err.(type) {
	case kafapi.ErrorDecoding, kafapi.ErrorMapping:
		return true
	default:
		return false
	}
}

func getGroupId(ctx context.Context, kMsg *kafka.Message) (string, error) {
	// Внимание!
	//
	// `kafapi` сервиса `auth` не разделяет парсинг метаданных и обработку сообщения
	// и не предоставляет возможности вклиниться в эти этапы для вставки кода DLQ -- и, в общем-то, не обязан.
	// Метаданные передаются в теле сообщения из-за особенностей продюсера, поэтому и нет смысла в явном разделении.
	//
	// Для `dlq.Decorate` метаданные должны быть известны до вызова обработчика сообщения
	// для предотвращения нормальной обработки и перенаправления сообщения в DLQ при обнаружении отравленной группы.
	// Такой подход обеспечивает сохранение порядка сообщений в каждой из групп --
	// это общий механизм, не зависящий от `kafapi` сервиса `auth` или кода обработки любого другого топика.
	//
	// Трюк с `kafapi`-сервисом для извлечения метаданных решает проблему несовместимости механизмов,
	// однако, приводит к повторному декодированию сообщения при вызове основного обработчика -- нужна оптимизация.
	// todo : избавиться от повторного декодирования при извлечении метаданных из auth.kafapi для DLQ.
	sgm := new(sDlqGetMeta)
	err := kafapi.Handle(sgm, ctx, kMsg)
	switch vErr := err.(type) {
	case nil:
		return sgm.Meta.GroupId, nil
	case kafapi.ErrorDecoding:
		return "", vErr
	case kafapi.ErrorMapping:
		return vErr.Meta.GroupId, nil
	case kafapi.ErrorUnsupported:
		return vErr.Meta.GroupId, nil
	default:
		panic(vErr)
	}
}
