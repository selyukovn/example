package kafcon

import (
	"context"
	"errors"
	"example/admin/gateway/internal/api/kafcon/kernel"
	"github.com/avast/retry-go/v5"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
	"strings"
	"sync/atomic"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Consumer struct {
	kCns        *kafka.Consumer
	kCnsRbErr   error
	kCnsStopErr error
	stopCalled  *atomic.Bool
	stopped     *atomic.Bool
	topic       string
	clientId    string
	handler     kernel.HandlerInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

const atLeastOnceEnableAutoCommit = true

func NewConsumer(
	service string,
	topic string,
	idInGroup string,
	brokerHostPorts []string,
	handler kernel.HandlerInterface,
) *Consumer {
	assert.Str().NotEmpty().Must(service)
	assert.Str().Word().Must(topic)
	assert.Str().NotEmpty().Must(idInGroup)
	assert.SliceCmp[[]string, string]().LenMin(1).Uniques().CustomElementEach("each", func(s string) bool {
		assert.Str().NotEmpty(). /* todo : UrlHostPort(). */ Must(s)
		return true
	}).Must(brokerHostPorts)
	assert.NotNilDeepMust(handler)

	// todo : принимать в аргументах?
	sessionTimeoutMs := 45_000   // 45 секунд -- значение по умолчанию
	maxPollIntervalMs := 300_000 // 300 секунд -- значение по умолчанию

	groupId := service + "-" + topic
	clientId := service + "-" + topic + "-" + idInGroup

	// https://github.com/confluentinc/librdkafka/blob/master/CONFIGURATION.md
	kCfg := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(brokerHostPorts, ","),

		// https://github.com/confluentinc/librdkafka/blob/master/INTRODUCTION.md#next-generation-consumer-group-protocol-kip-848
		// По умолчанию и так classic, но объявлено явно, т.к. в v2.12.0 и появился новый протокол "consumer".
		"group.protocol": "classic",

		"group.id":  groupId,
		"client.id": clientId,

		// В общем случае пусть лучше ошибка будет, чем непредсказуемые последствия.
		// А для точной настройки лучше использовать callback в `Subscribe()` -- см. `onReBalance()`
		"auto.offset.reset": "error",

		// Из описания `max.run.interval.ms`:
		// "... It is recommended to set enable.auto.offset.store=false for long-time processing applications
		// and then explicitly store offsets (using offsets_store()) after message processing,
		// to make sure offsets are not auto-committed prior to processing has finished. ..."
		"enable.auto.offset.store": false,

		// Вызов коммита вручную после каждого store-offset-вызова не имеет смысла (к тому же вызов блокирующий),
		// поскольку при перезапуске всеравно будут повторно вычитаны сообщения с момента последнего успешного коммита
		// с разницей лишь в количестве таких сообщений: вся автокоммитная пачка или одно при ручном вызове.
		// Поэтому проще использовать фоновый автоматический коммит (по умолчанию).
		"enable.auto.commit": atLeastOnceEnableAutoCommit,

		"session.timeout.ms":   sessionTimeoutMs,
		"max.poll.interval.ms": maxPollIntervalMs,
	}

	kCns, err := kafka.NewConsumer(kCfg)
	if err != nil {
		panic(err) // todo : вернуть ошибку?
	}

	return &Consumer{
		kCns:        kCns,
		kCnsRbErr:   nil,
		kCnsStopErr: nil,
		stopCalled:  new(atomic.Bool),
		stopped:     new(atomic.Bool),
		topic:       topic,
		clientId:    clientId,
		handler:     handler,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func kRetryErrNet[T any](call string, ctx context.Context, fn func() (T, error)) (T, error) {
	return retry.NewWithData[T](
		retry.Context(ctx),
		retry.Attempts(0 /* 0 -- бесконечно */),
		retry.MaxDelay(3*time.Second /* макс. 3 секунды -- примерное значение */),
		retry.RetryIf(func(err error) bool {
			kErr, isKafkaError := err.(kafka.Error)

			// возможно, есть еще сетевые или другие временные ошибки, но пока не попадались.
			if isKafkaError && (kErr.Code() == kafka.ErrResolve ||
				kErr.Code() == kafka.ErrTransport ||
				kErr.Code() == kafka.ErrAllBrokersDown) {
				return true
			}

			return false
		}),
		retry.OnRetry(func(n uint, err error) {
			logger.WarnFf(ctx, "%s не удался в %d-й раз: Err=%#v(%s)", call, n, err, err)
		}),
	).Do(fn)
}

func (c *Consumer) Start(ctx context.Context) error {
	assert.NotNilDeepMust(ctx)

	ctx = logger.AddAttrToCtx(ctx, "consumer", c.Id())

	if c.stopped.Load() || c.stopCalled.Load() {
		return std.NewErrorUnprocessableFf("Консьюмер уже остановлен")
	}

	defer func() {
		c.stopped.Store(true)
	}()

	defer func() {
		c.kCnsStopErr = c.kCns.Close()
	}()

	return c.run(ctx)
}

func (c *Consumer) run(ctx context.Context) error {
	_m_ := "run"

	// Внимание!
	// Уведомления о ре-балансировке поставляются последовательно вместе с сообщениями (см. `kafka.Poll()`),
	// а НЕ асинхронно, как может показаться при использовании `ReadMessage` и указании `rebalanceCb` в `Subscribe`.
	// Явное использование `Poll` могло бы это подчеркнуть, однако смысл такого подхода теряется,
	// поскольку `rebalanceCb` всеравно придется указывать в `Subscribe`, т.к. других способов это сделать нет,
	// а если этого не сделать, то при закрытии консьюмера не будет обработано финальное уведомление изъятия партиции,
	// и, например, сохранение оффсетов во внешнем хранилище не будет выполнено.
	// Поэтому `Subscribe` вызывается здесь, а не в точке инициализации консьюмера.
	if err := c.kCns.Subscribe(c.topic, func(_ *kafka.Consumer, e kafka.Event) error {
		// см. `kafka.handleRebalanceEvent()` -- ошибка из ребалансировки никуда не возвращается из `rebalanceCb`.
		// Но прервать обработку таки может быть необходимым.
		c.kCnsRbErr = c.onReBalance(ctx, e)
		return c.kCnsRbErr
	}); err != nil {
		return err
	}

	// --

	for {
		iterationId := strings.Replace(uuid.Must(uuid.NewRandom()).String(), "-", "", -1)
		ctx := logger.AddAttrToCtx(ctx, "consumer.iteration", iterationId)

		if c.stopCalled.Load() {
			logger.DebugFf(ctx, "Получен сигнал остановки")
			return nil
		}

		logger.DebugFf(ctx, "ReadMessage...")
		kMsg, err := kRetryErrNet[*kafka.Message]("ReadMessage", ctx, func() (*kafka.Message, error) {
			// `ReadMessage(-1)` для ожидания появления сообщения использовать было бы удобнее и, пожалуй, экономичнее,
			// однако в таком случае корректное завершение обработки исключено,
			// поскольку невозможно будет отреагировать на сигнал завершения работы до появления сообщения.
			_ = c.kCns.Close // -- не прерывает ожидание, поэтому таймауты обязательны.
			// Таймаут не может быть большим, чтобы хоть сколь-нибудь своевременно получать сигнал к остановке.
			// Вынесение чтения сообщений в отдельный поток не улучшит ситуацию, но добавит бессмысленной сложности.
			return c.kCns.ReadMessage(3 * time.Second /* 3 секунды -- примерное значение */)
		})

		if c.stopCalled.Load() {
			logger.DebugFf(ctx, "Получен сигнал остановки")
			return nil
		}

		kErr, isKErr := err.(kafka.Error)

		if c.kCnsRbErr != nil {
			if isKErr && !kErr.IsTimeout() {
				return errors.Join(c.kCnsRbErr, kErr)
			}
			return c.kCnsRbErr
		}

		if isKErr && kErr.IsTimeout() {
			logger.DebugFf(ctx, "ReadMessage: Timeout")
			continue
		}

		logger.DebugFf(ctx, "ReadMessage: KMsg=%#v; Err=%#v(%s)", kMsg, err, err)

		if err != nil {
			return err
		}

		// --

		logger.DebugFf(ctx, "Обработка сообщения...")
		err = c.handler.Handle(ctx, kMsg)
		if err != nil {
			logger.DebugFf(ctx, "Сообщение НЕ обработано: Err=%#v(%s)", err, err)
			return std.WrapErrorToRuntime(err, c, _m_, "handle")
		} else {
			logger.DebugFf(ctx, "Сообщение обработано")
		}

		// --

		// Внимание!
		// Четкое описание возвращаемых из `StoreMessage` ошибок отсутствует, тесты пакета-обертки тоже говорят мало,
		// а код `librdkafka` написан на C -- поэтому todo : приходится разбирать ошибки опытным путем....
		// Из кода обертки ясно, что при обращении к закрытому консьюмеру возвращается ошибка с кодом `ErrState`,
		// но на закрытый консьюмер здесь наткнуться нельзя, т.к. он закрывается после выхода из метода.
		logger.DebugFf(ctx, "Сохранение оффсета...")
		if _, err = kRetryErrNet[[]kafka.TopicPartition]("StoreMessage", ctx, func() ([]kafka.TopicPartition, error) {
			return c.kCns.StoreMessage(kMsg)
		}); err != nil {
			logger.DebugFf(ctx, "Оффсет НЕ сохранен: Err=%#v(%s)", err, err)
			return std.WrapErrorToRuntime(err, c, _m_, "StoreMessage")
		} else {
			logger.DebugFf(ctx, "Оффсет сохранен")
		}

		// В `at-least-once` семантике нет смысла коммитить оффсет вручную -- см. конфиг "enable.auto.commit".
		_ = atLeastOnceEnableAutoCommit
	}
}

func (c *Consumer) onReBalance(ctx context.Context, e kafka.Event) error {
	_m_ := "onReBalance"

	logger.InfoFf(ctx, "Уведомление о ре-балансировке: %#v", e)
	defer logger.DebugFf(ctx, "Уведомление о ре-балансировке принято")

	// todo : получение/сохранение оффсетов во внешнем хранилище ???

	switch eV := e.(type) {
	case kafka.AssignedPartitions:
		// Коды оффсетов: OffsetInvalid=-1001, OffsetBeginning=-2, OffsetStored=-1000, OffsetEnd=-1

		logger.DebugFf(ctx, "Получение ранее закоммиченных оффсетов...")
		cTps, err := kRetryErrNet[[]kafka.TopicPartition]("Committed", ctx, func() ([]kafka.TopicPartition, error) {
			return c.kCns.Committed(eV.Partitions, 3_000 /* 3 секунды -- примерное значение*/)
		})
		if err != nil {
			// Можно перечитать все с начала или ждать новых сообщений,
			// но если не удалось получить ранее закоммиченные оффсеты -- то ошибка, вероятно, серьезная,
			// поэтому лучше прервать обработку.
			logger.ErrorFf(ctx, "Ошибка получения ранее закоммиченных оффсетов: Err=%#v(%s)", err, err)
			return std.WrapErrorToRuntime(err, c, _m_, "Committed")
		} else {
			logger.InfoFf(ctx, "Ранее закоммиченные оффсеты: %#v", cTps)
		}

		for i, eTp := range eV.Partitions {
			if eTp.Offset == kafka.OffsetInvalid {

				cTpOffset := kafka.OffsetInvalid
				for _, cTp := range cTps {
					// топик -- указатель, сравнивать нужно значения!
					if *cTp.Topic == *eTp.Topic && cTp.Partition == eTp.Partition {
						cTpOffset = cTp.Offset
						break
					}
				}

				eTp.Offset = std.Ternary(cTpOffset == kafka.OffsetInvalid, kafka.OffsetBeginning, cTpOffset)
				eV.Partitions[i] = eTp
			}
		}

		logger.DebugFf(ctx, "Установка оффсетов: %#v", eV.Partitions)
		err = c.kCns.Assign(eV.Partitions)
		if err != nil {
			logger.ErrorFf(ctx, "Ошибка установки оффсетов: Err=%#v(%s)", err, err)
			return std.WrapErrorToRuntime(err, c, _m_, "Assign")
		} else {
			logger.DebugFf(ctx, "Оффсеты установлены")
		}
	case kafka.RevokedPartitions:
		_ = c.kCns.Close    // -- отмечает консьюмер завершенным уже после обработки финальной ре-балансировки,
		_ = c.kCns.IsClosed // -- поэтому здесь проверка закрытия консьюмера не сработает.

		// https://github.com/confluentinc/confluent-kafka-go/blob/v2.12.0/examples/consumer_rebalance_example/consumer_rebalance_example.go#L183
		// "... there can be cases where the assignment is lost involuntarily.
		// In this case, the partition might already be owned by another consumer,
		// and operations including committing offsets may not work. ..."
		if c.kCns.AssignmentLost() {
			// В семантике `at-least-once` это нормально -- другой консьюмер обработает и сохранит.
			logger.WarnFf(ctx, "Обработка изъятия партиции прервана -- партиции утрачены")
			return nil
		}

		logger.DebugFf(ctx, "Коммит накопившихся оффсетов...")
		if _, err := kRetryErrNet[[]kafka.TopicPartition]("Commit", ctx, func() ([]kafka.TopicPartition, error) {
			return c.kCns.Commit()
		}); err != nil {
			kErr, isKErr := err.(kafka.Error)
			if isKErr && kErr.Code() == kafka.ErrNoOffset {
				logger.DebugFf(ctx, "Нечего коммитить -- нет накопленных оффсетов")
			} else {
				// В семантике `at-least-once` это нормально -- другой консьюмер обработает и сохранит.
				logger.WarnFf(ctx, "Накопившиеся оффсеты не закоммичены: Err=%#v(%s)", err, err)
			}
		} else {
			logger.DebugFf(ctx, "Накопившиеся оффсеты закоммичены")
		}
	default:
		panic(std.NewErrorRuntimeFf("Неизвестное событие ре-балансировки топика %q: %#v", c.topic, e))
	}

	return nil
}

func (c *Consumer) Stop(ctx context.Context) error {
	assert.NotNilDeepMust(ctx)

	if c.stopped.Load() {
		return std.NewErrorAlreadyDoneFf("Уже остановлен")
	}

	c.stopCalled.Store(true)

	logger.DebugFf(ctx, "Ждем завершения остановки...")
	for !c.stopped.Load() {
		// 100 ms -- время ожидания событий при закрытии kafka-консьюмера --
		// соответственно, использовать меньше значение смысла нет, а ограничение сверху примерное.
		_ = c.kCns.Close
		time.Sleep(100 * time.Millisecond)
	}

	logger.DebugFf(ctx, "Остановка завершена")

	return c.kCnsStopErr
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (c *Consumer) Id() string {
	return c.clientId
}

// ---------------------------------------------------------------------------------------------------------------------
