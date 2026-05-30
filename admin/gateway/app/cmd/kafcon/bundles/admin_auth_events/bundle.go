package admin_auth_events

import (
	"context"
	"example/admin/gateway/cmd/common/components/processing"
	"example/admin/gateway/cmd/kafcon/bundles/admin_auth_events/handlers"
	"example/admin/gateway/cmd/kafcon/bundles/admin_auth_events/kafapi"
	"example/admin/gateway/cmd/kafcon/components/dlq"
	"example/admin/gateway/cmd/kafcon/container"
	"example/admin/gateway/cmd/kafcon/kernel"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
	"strconv"
)

// ---------------------------------------------------------------------------------------------------------------------
// Register
// ---------------------------------------------------------------------------------------------------------------------

const TopicName = "admin-auth-events"

func Register(topicRouting map[string]kernel.TopicHandlerInterface, ctr *container.Container) {
	assert.NotNilDeepMust(topicRouting)
	assert.NotNilDeepMust(ctr)

	var topicHandler kernel.TopicHandlerInterface
	topicHandler = newTopicHandlerDefault(ctr)
	topicHandler = newTopicHandlerDecoratorDlq(
		topicHandler,
		ctr.Dlq.GroupTracker,
		ctr.Dlq.Storage,
		ctr.Dlq.TopicHolder,
	)

	topicRouting[TopicName] = topicHandler
}

// ---------------------------------------------------------------------------------------------------------------------
// Default
// ---------------------------------------------------------------------------------------------------------------------

func newTopicHandlerDefault(ctr *container.Container) kernel.TopicHandlerInterface {
	service := sDecoratorBoundary{
		origin: sDefault{
			sessionClosedV1: handlers.NewSessionClosedV1(ctr),
		},
	}

	return kernel.FnTopicHandler(func(ctx context.Context, kMsg *kafka.Message) error {
		_o_ := TopicName
		_m_ := "newTopicHandlerDefault"

		err := kafapi.Handle(service, ctx, kMsg)
		switch vErr := err.(type) {
		case nil:
			return nil
		case kafapi.ErrorDecoding, kafapi.ErrorMapping:
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

// Default: Service Default
// ---------------------------------------------------------------------------------------------------------------------

var _ kafapi.ServiceInterface = sDefault{}

type sDefault struct {
	kafapi.EmptyService
	sessionClosedV1 func(context.Context, *kafapi.Meta, *kafapi.DataSessionClosedV1) error
}

func (s sDefault) SessionClosedV1(ctx context.Context, meta *kafapi.Meta, msg *kafapi.DataSessionClosedV1) error {
	return s.sessionClosedV1(ctx, meta, msg)
}

// Default: Service Decorator Boundary
// ---------------------------------------------------------------------------------------------------------------------

var _ kafapi.ServiceInterface = sDecoratorBoundary{}

type sDecoratorBoundary struct {
	origin kafapi.ServiceInterface
}

func (s sDecoratorBoundary) enrichCtx(ctx context.Context, meta *kafapi.Meta) context.Context {
	operationId := meta.OperationId
	ctx = processing.EnrichCtx(ctx, operationId)
	ctx = logger.AddAttrToCtx(ctx, "processing.OperationId", operationId)
	/* см. */ _ = processing.OperationId
	return ctx
}

func (s sDecoratorBoundary) AccountCreatedV1(ctx context.Context, meta *kafapi.Meta, data *kafapi.DataAccountCreatedV1) error {
	ctx = s.enrichCtx(ctx, meta)

	logger.InfoFf(ctx, "%T/%T - start: %+v (%+v)", s, data, data, meta)
	err := s.origin.AccountCreatedV1(ctx, meta, data)
	logger.InfoFf(ctx, "%T/%T - end: %#v", s, data, err)

	return err
}

func (s sDecoratorBoundary) AccountDeactivatedV1(ctx context.Context, meta *kafapi.Meta, data *kafapi.DataAccountDeactivatedV1) error {
	ctx = s.enrichCtx(ctx, meta)

	logger.InfoFf(ctx, "%T/%T - start: %+v (%+v)", s, data, data, meta)
	err := s.origin.AccountDeactivatedV1(ctx, meta, data)
	logger.InfoFf(ctx, "%T/%T - end: %#v", s, data, err)

	return err
}

func (s sDecoratorBoundary) IpWhitelistChangedV1(ctx context.Context, meta *kafapi.Meta, data *kafapi.DataIpWhitelistChangedV1) error {
	ctx = s.enrichCtx(ctx, meta)

	logger.InfoFf(ctx, "%T/%T - start: %+v (%+v)", s, data, data, meta)
	err := s.origin.IpWhitelistChangedV1(ctx, meta, data)
	logger.InfoFf(ctx, "%T/%T - end: %#v", s, data, err)

	return err
}

func (s sDecoratorBoundary) SessionCreatedV1(ctx context.Context, meta *kafapi.Meta, data *kafapi.DataSessionCreatedV1) error {
	ctx = s.enrichCtx(ctx, meta)

	msgMasked := data
	msgMasked.SessionId = std.MaskStrNotFirstLast(msgMasked.SessionId)
	logger.InfoFf(ctx, "%T/%T - start: %+v (%+v)", s, data, msgMasked, meta)
	err := s.origin.SessionCreatedV1(ctx, meta, data)
	logger.InfoFf(ctx, "%T/%T - end: %#v", s, data, err)

	return err
}

func (s sDecoratorBoundary) SessionClosedV1(ctx context.Context, meta *kafapi.Meta, data *kafapi.DataSessionClosedV1) error {
	ctx = s.enrichCtx(ctx, meta)

	msgMasked := data
	msgMasked.SessionId = std.MaskStrNotFirstLast(msgMasked.SessionId)
	logger.InfoFf(ctx, "%T/%T - start: %+v (%+v)", s, data, msgMasked, meta)
	err := s.origin.SessionClosedV1(ctx, meta, data)
	logger.InfoFf(ctx, "%T/%T - end: %#v", s, data, err)

	return err
}

// ---------------------------------------------------------------------------------------------------------------------
// Decorator DLQ
// ---------------------------------------------------------------------------------------------------------------------

func newTopicHandlerDecoratorDlq(
	origin kernel.TopicHandlerInterface,
	dlqGroupTracker dlq.GroupTrackerInterface,
	dlqStorage dlq.StorageInterface,
	dlqTopicHolder dlq.TopicHolderInterface,
) kernel.TopicHandlerInterface {
	fnGetMeta := func(ctx context.Context, kMsg *kafka.Message) (
		rId string,
		rGroupId string,
		rErr error,
	) {
		// Внимание!
		//
		// `kafapi` сервиса `auth` не разделяет парсинг метаданных и обработку сообщения
		// и не предоставляет возможности вклиниться в эти этапы для вставки кода DLQ -- и, в общем-то, не обязан.
		// Метаданные передаются в теле сообщения из-за особенностей продюсера, поэтому и нет смысла в явном разделении.
		//
		// Для `dlq.NewTopicHandlerDecoratorFn` метаданные должны быть известны до вызова обработчика сообщения
		// для предотвращения нормальной обработки и перенаправления сообщения в DLQ при обнаружении отравленной группы.
		// Такой подход обеспечивает сохранение порядка сообщений в каждой из групп --
		// это общий механизм, не зависящий от `kafapi` сервиса `auth` или кода обработки любого другого топика.
		//
		// Трюк с `kafapi`-сервисом для извлечения метаданных решает проблему несовместимости механизмов,
		// однако, приводит к повторному декодированию сообщения при вызове основного обработчика -- нужна оптимизация.
		// todo : избавиться от повторного декодирования при извлечении метаданных из auth.kafapi для DLQ.
		trickyService := new(sTrickyServiceToGetMetadataForDlq)
		err := kafapi.Handle(trickyService, ctx, kMsg)
		switch vErr := err.(type) {
		case nil:
			rId = strconv.FormatUint(uint64(trickyService.Meta.Id), 10)
			rGroupId = trickyService.Meta.GroupId
			rErr = nil
			return
		case kafapi.ErrorDecoding:
			rId = ""
			rGroupId = ""
			rErr = vErr
			return
		case kafapi.ErrorMapping:
			rId = strconv.FormatUint(uint64(vErr.Meta.Id), 10)
			rGroupId = vErr.Meta.GroupId
			rErr = nil
			return
		case kafapi.ErrorUnsupported:
			rId = strconv.FormatUint(uint64(vErr.Meta.Id), 10)
			rGroupId = vErr.Meta.GroupId
			rErr = nil
			return
		default:
			panic(vErr)
		}
		return
	}

	fnIsDlqCase := func(_ context.Context, err error) bool {
		switch err.(type) {
		case kafapi.ErrorDecoding, kafapi.ErrorMapping:
			return true
		default:
			return false
		}
	}

	fnDecorate := dlq.NewTopicHandlerDecoratorFn(
		dlqGroupTracker,
		dlqStorage,
		dlqTopicHolder,
		TopicName,
		fnGetMeta,
		fnIsDlqCase,
	)

	return fnDecorate(origin)
}

// Decorator DLQ: Tricky Service to Parse Metadata
// ---------------------------------------------------------------------------------------------------------------------

var _ kafapi.ServiceInterface = new(sTrickyServiceToGetMetadataForDlq)

type sTrickyServiceToGetMetadataForDlq struct {
	Meta *kafapi.Meta
}

func (s *sTrickyServiceToGetMetadataForDlq) handle(meta *kafapi.Meta) error {
	s.Meta = meta
	return nil
}

func (s *sTrickyServiceToGetMetadataForDlq) AccountCreatedV1(_ context.Context, meta *kafapi.Meta, _ *kafapi.DataAccountCreatedV1) error {
	return s.handle(meta)
}
func (s *sTrickyServiceToGetMetadataForDlq) AccountDeactivatedV1(_ context.Context, meta *kafapi.Meta, _ *kafapi.DataAccountDeactivatedV1) error {
	return s.handle(meta)
}
func (s *sTrickyServiceToGetMetadataForDlq) IpWhitelistChangedV1(_ context.Context, meta *kafapi.Meta, _ *kafapi.DataIpWhitelistChangedV1) error {
	return s.handle(meta)
}
func (s *sTrickyServiceToGetMetadataForDlq) SessionCreatedV1(_ context.Context, meta *kafapi.Meta, _ *kafapi.DataSessionCreatedV1) error {
	return s.handle(meta)
}
func (s *sTrickyServiceToGetMetadataForDlq) SessionClosedV1(_ context.Context, meta *kafapi.Meta, _ *kafapi.DataSessionClosedV1) error {
	return s.handle(meta)
}

// ---------------------------------------------------------------------------------------------------------------------
