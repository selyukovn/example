package loggable

import (
	"context"
	"example/admin/gateway/internal/api/kafcon/handlers/admin_auth_events/kafapi"
	"github.com/selyukovn/example_gopkg/processing"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ kafapi.ServiceInterface = Decorator{}

type Decorator struct {
	origin kafapi.ServiceInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewDecorator(origin kafapi.ServiceInterface) Decorator {
	assert.NotNilDeepMust(origin)

	return Decorator{
		origin: origin,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (d Decorator) enrichCtx(ctx context.Context, meta *kafapi.Meta) context.Context {
	operationId := meta.OperationId
	ctx = processing.EnrichCtx(ctx, operationId)
	ctx = logger.AddAttrToCtx(ctx, "processing.OperationId", operationId)
	return ctx
}

func (d Decorator) AccountCreatedV1(ctx context.Context, meta *kafapi.Meta, data *kafapi.DataAccountCreatedV1) error {
	ctx = d.enrichCtx(ctx, meta)
	return d.origin.AccountCreatedV1(ctx, meta, data)
}

func (d Decorator) AccountDeactivatedV1(ctx context.Context, meta *kafapi.Meta, data *kafapi.DataAccountDeactivatedV1) error {
	ctx = d.enrichCtx(ctx, meta)
	return d.origin.AccountDeactivatedV1(ctx, meta, data)
}

func (d Decorator) IpWhitelistChangedV1(ctx context.Context, meta *kafapi.Meta, data *kafapi.DataIpWhitelistChangedV1) error {
	ctx = d.enrichCtx(ctx, meta)
	return d.origin.IpWhitelistChangedV1(ctx, meta, data)
}

func (d Decorator) SessionCreatedV1(ctx context.Context, meta *kafapi.Meta, data *kafapi.DataSessionCreatedV1) error {
	ctx = d.enrichCtx(ctx, meta)
	return d.origin.SessionCreatedV1(ctx, meta, data)
}

func (d Decorator) SessionClosedV1(ctx context.Context, meta *kafapi.Meta, data *kafapi.DataSessionClosedV1) error {
	ctx = d.enrichCtx(ctx, meta)
	return d.origin.SessionClosedV1(ctx, meta, data)
}

// ---------------------------------------------------------------------------------------------------------------------
