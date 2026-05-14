package admin_auth_events

import (
	"context"
	"example/admin/gateway/cmd/common/components/processing"
	"example/admin/gateway/cmd/kafcon/bundles/admin_auth_events/handlers"
	"example/admin/gateway/cmd/kafcon/bundles/admin_auth_events/kafapi"
	"example/admin/gateway/cmd/kafcon/container"
	"example/admin/gateway/cmd/kafcon/kernel"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Register
// ---------------------------------------------------------------------------------------------------------------------

const TopicName = "admin-auth-events"

func Register(topicRouting map[string]kernel.TopicHandlerInterface, ctr *container.Container) {
	assert.NotNilDeepMust(topicRouting)
	assert.NotNilDeepMust(ctr)

	service := sDecoratorBoundary{
		origin: sDefault{
			sessionClosedV1: handlers.NewSessionClosedV1(ctr),
		},
	}

	topicRouting[TopicName] = kernel.FnTopicHandler(func(ctx context.Context, kMsg *kafka.Message) error {
		err := kafapi.Handle(service, ctx, kMsg)
		switch vErr := err.(type) {
		case nil:
			// pass
		case kafapi.ErrorDecoding:
			// todo : ???
			return vErr
		case kafapi.ErrorMapping:
			// todo : ???
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
			// todo : if err ???
			return vErr
		default:
			panic(vErr)
		}

		return nil
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// Default
// ---------------------------------------------------------------------------------------------------------------------

var _ kafapi.ServiceInterface = sDefault{}

type sDefault struct {
	kafapi.EmptyService
	sessionClosedV1 func(context.Context, *kafapi.Meta, *kafapi.DataSessionClosedV1) error
}

func (s sDefault) SessionClosedV1(ctx context.Context, meta *kafapi.Meta, msg *kafapi.DataSessionClosedV1) error {
	return s.sessionClosedV1(ctx, meta, msg)
}

// ---------------------------------------------------------------------------------------------------------------------
// Decorator Boundary
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
