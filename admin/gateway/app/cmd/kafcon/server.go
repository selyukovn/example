package kafcon

import (
	"context"
	"errors"
	"example/admin/gateway/cmd/kafcon/container"
	"github.com/avast/retry-go/v5"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
	"strings"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Server struct {
	mu             *sync.Mutex
	config         *kafka.ConfigMap
	topicName      string
	actualConsumer *kafka.Consumer
	router         topicRouter
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

const enableAutoCommit = true

// NewServer
//
// Паникует при нулевых аргументах.
func NewServer(
	ownerServiceName string,
	brokerHostPorts []string,
	topicName string,
	idInGroup string,
	ctr *container.Container,
) *Server {
	assert.Str().Word().Must(ownerServiceName)
	assert.SliceCmp[[]string, string]().LenMin(1).Uniques().CustomElementEach("each", func(s string) bool {
		assert.Str().NotEmpty(). /* todo : UrlHostPort(). */ Must(s)
		return true
	}).Must(brokerHostPorts)
	assert.Str().Word().Must(topicName)
	assert.Str().NotEmpty().Must(idInGroup)
	assert.NotNilDeepMust(ctr)

	groupId := ownerServiceName + "-" + topicName
	clientId := groupId + "-" + idInGroup

	// https://github.com/confluentinc/librdkafka/blob/master/CONFIGURATION.md
	config := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(brokerHostPorts, ","),

		// https://github.com/confluentinc/librdkafka/blob/master/INTRODUCTION.md#next-generation-consumer-group-protocol-kip-848
		// По умолчанию и так classic, но объявлено явно, т.к. в v2.12.0 и появился новый протокол "consumer".
		"group.protocol": "classic",

		"group.id":  groupId,
		"client.id": clientId,

		// В общем случае пусть лучше ошибка будет, чем непредсказуемые последствия.
		// А для точной настройки лучше использовать callback'и (см. Subscribe()).
		// Однако, на текущем этапе использование callback без внешнего хранилища равносильно "smallest" значению.
		// todo : "auto.offset.reset: error", внешнее хранилище оффсетов и callback в Subscribe().
		"auto.offset.reset": "smallest",

		// Из описания `max.poll.interval.ms`:
		// "... It is recommended to set enable.auto.offset.store=false for long-time processing applications
		// and then explicitly store offsets (using offsets_store()) after message processing,
		// to make sure offsets are not auto-committed prior to processing has finished. ..."
		"enable.auto.offset.store": false,

		// Вызов коммита вручную после каждого store-offset-вызова не имеет смысла (к тому же вызов блокирующий),
		// поскольку при перезапуске всеравно будут повторно вычитаны сообщения с момента последнего успешного коммита
		// с разницей лишь в количестве таких сообщений: вся автокоммитная пачка или одно при ручном вызове.
		// Поэтому проще использовать фоновый автоматический коммит (по умолчанию).
		"enable.auto.commit": enableAutoCommit,
	}

	// Мьютекс гарантирует, что инициализация консьюмера и каждая итерация обработки сообщений будут выполнены,
	// даже если `Server.Stop()` был вызван в процессе их выполнения.
	// Это упрощает код: не нужно постоянно проверять закрытие консьюмера и обрабатывать соответствующие ошибки.
	mu := new(sync.Mutex)

	return &Server{
		mu:             mu,
		config:         config,
		topicName:      topicName,
		actualConsumer: nil,
		router:         newRouter(ctr),
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (s *Server) Start(ctx context.Context) error {
	clientId := assert.Str().NotEmpty().MustGet(std.Must(s.config.Get("client.id", "")).(string))

	ctx = logger.AddAttrToCtx(ctx, "server", clientId)

	err := s.openConsumer(ctx)
	if err != nil {
		return err
	}

	return s.loop(ctx)
}

// Отдельный метод нужен для использования `defer mutex.Unlock()`.
func (s *Server) openConsumer(ctx context.Context) error {
	_m_ := "openConsumer"

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.actualConsumer != nil && !s.actualConsumer.IsClosed() {
		return nil
	}

	// --

	newConsumer, err := kafka.NewConsumer(s.config)
	if err != nil {
		return std.WrapErrorToRuntime(err, s, _m_, "NewConsumer")
	}

	if err = newConsumer.Subscribe(s.topicName, func(c *kafka.Consumer, e kafka.Event) error {
		clientId := assert.Str().NotEmpty().MustGet(std.Must(s.config.Get("client.id", "")).(string))
		logger.InfoFf(ctx, "Ребалансировка %q: событие %#v", clientId, e)
		return nil
	}); err != nil {
		return std.WrapErrorToRuntime(err, s, _m_, "Subscribe")
	}

	s.actualConsumer = newConsumer

	return nil
}

func (s *Server) loop(ctx context.Context) error {
	for {
		iterationId := strings.Replace(uuid.Must(uuid.NewRandom()).String(), "-", "", -1)
		ctx := logger.AddAttrToCtx(ctx, "server.iteration", iterationId)

		err, toBreak := s.iteration(ctx)

		if err != nil {
			logger.ErrorFf(ctx, err.Error())
		}

		if toBreak {
			stopErr := s.stop(ctx)
			if stopErr != nil {
				return errors.Join(err, stopErr)
			}

			return err
		}
	}
}

func kRetry[T any](call string, ctx context.Context, fn func() (T, error)) (T, error) {
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
			logger.WarnFf(ctx, "%s не удался в %d-й раз: %s (err=%#v)", call, n, err, err)
		}),
	).Do(fn)
}

// Отдельный метод нужен для использования `defer mutex.Unlock()`.
func (s *Server) iteration(ctx context.Context) (resultError error, toBreak bool) {
	_m_ := "iteration"

	kMsg, err := kRetry[*kafka.Message]("ReadMessage", ctx, func() (*kafka.Message, error) {
		return s.actualConsumer.ReadMessage(-1 /* -1 == блокируется до появления сообщения. todo : таймаут ??? */)
	})

	// Мьютекс должен быть закрыт после `ReadMessage`,
	// поскольку `ReadMessage` блокирует выполнение до появления сообщения или на указанное время.
	// Если закрыть мьютекс до, то `Server.Stop()` не сможет закрыть консьюмер, ожидая открытия мьютекса.
	s.mu.Lock()
	defer s.mu.Unlock()

	// `ReadMessage` все еще может напороться на закрытый консьюмер,
	// или консьюмер может закрыться до закрытия мьютекса (во время или после `ReadMessage`).
	if s.actualConsumer.IsClosed() {
		return nil, true
	}

	if err != nil {
		err = std.WrapErrorToRuntime(err, s, _m_, "ReadMessage")
		return err, true
	}

	// Внимание!
	// Гарантии доставки (за счет дедупликации, т.к. `StoreMessage` вызывается после обработки -- имеем at-least-once),
	// дедупликация, повторная обработка при ошибках, dead-letter-queue и т.д. --
	// все эти операции зависят от конкретного потребителя, т.е. от специфического бандла или даже хендлера в бандле,
	// поэтому должны быть реализованы НЕ здесь (на уровне консьюмера), а на уровне бандлов (например, в middleware).
	// Кроме того, для подобных операций нужны соответствующие метаданные (например, идентификатор группы сообщений),
	// которые каждый продюсер устанавливает по-своему -- что снова переводит ответственность на бандлы.
	//
	// К слову о получении метаданных:
	// можно было бы стандартизировать заголовки, но это сложно из-за разнородных продюсеров (debezium, рукотворные, ...)
	// и необходимости в едином "источнике" стандарта (например, импортируемый пакет или контролируемая копипаста),
	// а также бессмысленно из-за уже упомянутой зависимости от специфического потребителя, т.е. бандла/хендлера.
	//
	// Касаемо повторной обработки в случае возникновения ошибки:
	// если не сохранить оффсет (и не коммитить), при перезапуске консьюмера чтение начнется со старой позиции.
	// Однако, в текущем потоке прочитанное сообщение не будет прочитано заново --
	// т.е. `ReadMessage` всеравно прочитает следующее сообщение, даже если текущий оффсет не был зафиксирован.
	// Поэтому повторная обработка должна быть организована на уровне конкретных хендлеров.
	//
	// todo : параллельная обработка разных групп сообщений из одной партиции ???
	// Проще решать количеством консьюмеров и партиций, а не на уровне кода, поскольку сообщения считываются по одному.
	// Что если следующее сообщение из группы А уже будет обработано, а обработка текущего из группы Б провалится?

	err = s.router.handle(ctx, kMsg)
	if err != nil {
		err = std.WrapErrorToRuntime(err, s, _m_, "handle")
		return err, true
	}

	if _, err = kRetry[[]kafka.TopicPartition]("StoreMessage", ctx, func() ([]kafka.TopicPartition, error) {
		return s.actualConsumer.StoreMessage(kMsg)
	}); err != nil {
		err = std.WrapErrorToRuntime(err, s, _m_, "StoreMessage")
		return err, true
	}

	// В текущей ситуации нет смысла коммитить оффсет вручную -- см. конфиг "enable.auto.commit".
	_ = enableAutoCommit

	return nil, false
}

func (s *Server) stop(ctx context.Context) error {
	// todo : возможно, есть смысл ограничить по времени

	if s.actualConsumer == nil || s.actualConsumer.IsClosed() {
		return nil
	}

	// Теоретически в момент вызова `Stop()` в буфере могут находиться неотправленные оффсеты,
	// т.е. логично было бы вызвать коммит перед закрытием,
	// однако, финальный коммит и так выполняется librdkafka'ой внутри вызова `consumer.Close()`.
	// https://github.com/confluentinc/librdkafka/blob/master/INTRODUCTION.md#high-level-kafkaconsumer

	return s.actualConsumer.Close()
}

func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.stop(ctx)
}

// ---------------------------------------------------------------------------------------------------------------------
