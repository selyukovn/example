package loggable

import (
	"context"
	"example/admin/gateway/internal/api/kafcon/handlers/admin_auth_events/kafapi"
	"github.com/selyukovn/go-std"
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

func (d Decorator) AccountCreatedV1(ctx context.Context, meta *kafapi.Meta, data *kafapi.DataAccountCreatedV1) error {
	logger.InfoFf(ctx, "%T/%T - start: %+v (%+v)", d, data, data, meta)
	err := d.origin.AccountCreatedV1(ctx, meta, data)
	logger.InfoFf(ctx, "%T/%T - end: %#v", d, data, err)
	return err
}

func (d Decorator) AccountDeactivatedV1(ctx context.Context, meta *kafapi.Meta, data *kafapi.DataAccountDeactivatedV1) error {
	logger.InfoFf(ctx, "%T/%T - start: %+v (%+v)", d, data, data, meta)
	err := d.origin.AccountDeactivatedV1(ctx, meta, data)
	logger.InfoFf(ctx, "%T/%T - end: %#v", d, data, err)
	return err
}

func (d Decorator) IpWhitelistChangedV1(ctx context.Context, meta *kafapi.Meta, data *kafapi.DataIpWhitelistChangedV1) error {
	logger.InfoFf(ctx, "%T/%T - start: %+v (%+v)", d, data, data, meta)
	err := d.origin.IpWhitelistChangedV1(ctx, meta, data)
	logger.InfoFf(ctx, "%T/%T - end: %#v", d, data, err)
	return err
}

func (d Decorator) SessionCreatedV1(ctx context.Context, meta *kafapi.Meta, data *kafapi.DataSessionCreatedV1) error {
	msgMasked := data
	msgMasked.SessionId = std.MaskStrNotFirstLast(msgMasked.SessionId)
	logger.InfoFf(ctx, "%T/%T - start: %+v (%+v)", d, data, msgMasked, meta)
	err := d.origin.SessionCreatedV1(ctx, meta, data)
	logger.InfoFf(ctx, "%T/%T - end: %#v", d, data, err)
	return err
}

func (d Decorator) SessionClosedV1(ctx context.Context, meta *kafapi.Meta, data *kafapi.DataSessionClosedV1) error {
	msgMasked := data
	msgMasked.SessionId = std.MaskStrNotFirstLast(msgMasked.SessionId)
	logger.InfoFf(ctx, "%T/%T - start: %+v (%+v)", d, data, msgMasked, meta)
	err := d.origin.SessionClosedV1(ctx, meta, data)
	logger.InfoFf(ctx, "%T/%T - end: %#v", d, data, err)
	return err
}

// ---------------------------------------------------------------------------------------------------------------------
