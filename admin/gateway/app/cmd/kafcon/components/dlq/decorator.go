package dlq

import (
	"context"
	"example/admin/gateway/cmd/kafcon/kernel"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
)

func NewTopicHandlerDecoratorFn(
	groupTracker GroupTrackerInterface,
	storage StorageInterface,
	topicHolder TopicHolderInterface,
	topic string,
	fnGetMeta func(ctx context.Context, kMsg *kafka.Message) (rId string, rGroupId string, rErr error),
	fnIsDlqCase func(ctx context.Context, err error) bool,
) func(kernel.TopicHandlerInterface) kernel.TopicHandlerInterface {
	return func(next kernel.TopicHandlerInterface) kernel.TopicHandlerInterface {
		return kernel.FnTopicHandler(func(ctx context.Context, kMsg *kafka.Message) error {
			_o_ := "dlq"
			_m_ := "NewTopicHandlerDecoratorFn"

			// Внимание!
			// При запуске обработки исцеленной группы из DLQ необходима блокировка обработки всего топика,
			// иначе при достаточно плотном поступлении сообщений порядок их обработки может быть нарушен.
			// Предположим, что в DLQ уже есть несколько сообщений некоторой группы, причина перенаправления устранена,
			// был запущен DLQ-обработчик для этой группы, а консьюмер продолжает получать новые сообщения.
			// Без блокировки новые сообщения продолжат перенаправляться в DLQ, пока с группы не снята метка перенаправления.
			// С момента обработки последнего сообщения из DLQ до момента снятия метки перенаправления
			// новое сообщение может быть добавлено в DLQ -- и станет потерянным, а порядок сообщений в группе нарушится,
			// поскольку DLQ-обработка уже завершена, а обработка новых сообщений возвращена в нормальное русло.
			// Блокировка партиции, а не топика, не даст необходимых гарантий,
			// поскольку продюсер может начать отправку той же группы в другую партицию при ее добавлении.
			err := topicHolder.WaitTillOnHold(ctx, topic)
			if err != nil {
				return std.WrapErrorToRuntime(err, _o_, _m_, "WaitTillOnHold")
			}

			// --

			// Если из полученного ранее отравленного сообщения не удалось извлечь идентификатор группы,
			// то весь топик считается отравленным, чтобы не нарушить порядок сообщений.
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

			// Внимание!
			// Если метаданные сообщения неизвестны, в DLQ придется перенаправлять все сообщения топика,
			// иначе будет нарушен порядок обработки сообщений одной из групп.
			// Это может быть накладно, но остановка консьюмера не позволит исправить данные перед повторной обработкой,
			// а достаточно малый retention основного топика может привести и к потере этого сообщения.
			// Также сообщения из неизвестной группы могут начать отправляться в новую партицию при ее добавлении,
			// поэтому для соблюдения порядка полученных сообщений перенаправление должно быть на уровне топика.
			mId, mGroupId, err := fnGetMeta(ctx, kMsg)
			if err != nil {
				err := storage.Add(ctx, topic, AnyGroup, kMsg)
				if err != nil {
					return std.WrapErrorToRuntime(err, _o_, _m_, "Parse")
				}
				return nil // !!!
			}

			isGroupPoisoned, err := storage.IsGroupPoisoned(ctx, topic, mGroupId)
			if err != nil {
				return std.WrapErrorToRuntime(err, _o_, _m_, "IsGroupPoisoned", "Group")
			} else if isGroupPoisoned {
				err = storage.Add(ctx, topic, mGroupId, kMsg)
				if err != nil {
					return std.WrapErrorToRuntime(err, _o_, _m_, "IsGroupPoisoned", "Group", "Add")
				}
				return nil
			}

			// --

			rErr := next.Handle(ctx, kMsg)
			if fnIsDlqCase(ctx, rErr) {
				err := storage.Add(ctx, topic, mGroupId, kMsg)
				if err != nil {
					return std.WrapErrorToRuntime(err, _o_, _m_, "fnIsDlqCase", "Add")
				}
				return nil // !!!
			}

			// --

			err = groupTracker.SetLastHandled(ctx, topic, mGroupId, mId)
			if err != nil {
				logger.ErrorFf(ctx, std.WrapErrorToRuntime(err, _o_, _m_, "SetLastHandled").Error())
			}

			return rErr
		})
	}
}
