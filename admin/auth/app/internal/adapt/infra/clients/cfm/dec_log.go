package cfm

import (
	"context"
	"example/admin/auth/internal/infra/clients/cfm"
	"example/admin/auth/internal/infra/logger"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type DecoratorLoggable struct {
	origin cfm.ClientInterface
	logger *logger.Logger
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewDecoratorLoggable
//
// Паникует при нулевых аргументах.
func NewDecoratorLoggable(origin cfm.ClientInterface, logger *logger.Logger) *DecoratorLoggable {
	assert.NotNilDeepMust(origin)
	assert.NotNilDeepMust(logger)

	return &DecoratorLoggable{
		origin: origin,
		logger: logger,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// CreateForEmail
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (d *DecoratorLoggable) CreateForEmail(ctx context.Context, traceId string, email std.Email) (
	rRes cfm.CreateForEmailResult,
	rErr error,
) {
	d.logger.CtxInfoFf(ctx, "%T.%s - start(%q)", d, "CreateForEmail", email)
	defer func() { d.logger.CtxInfoFf(ctx, "%T.%s - end(%#v, %#v = %s)", d, "CreateForEmail", rRes, rErr, rErr) }()

	return d.origin.CreateForEmail(ctx, traceId, email)
}

// Request
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - cfm.ErrorFinished
//   - cfm.ErrorNoAttemptsLeft
//   - cfm.ErrorRequestsFrequency
//   - std.ErrorRuntime
func (d *DecoratorLoggable) Request(ctx context.Context, traceId string, cfmId string) (
	rRes cfm.RequestResult,
	rErr error,
) {
	d.logger.CtxInfoFf(ctx, "%T.%s - start(%q)", d, "Request", cfmId)
	defer func() { d.logger.CtxInfoFf(ctx, "%T.%s - end(%#v, %#v = %s)", d, "Request", rRes, rErr, rErr) }()

	return d.origin.Request(ctx, traceId, cfmId)
}

// Confirm
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - cfm.ErrorFinished
//   - std.ErrorUnprocessable -- если не была запрошена
//   - std.ErrorRuntime
func (d *DecoratorLoggable) Confirm(ctx context.Context, traceId string, cfmId string, code string) (
	rRes cfm.ConfirmResult,
	rErr error,
) {
	d.logger.CtxInfoFf(
		ctx,
		"%T.%s - start(%q, %q)",
		d, "Confirm", cfmId, std.MaskStrNotFirstLast(code),
	)
	defer func() { d.logger.CtxInfoFf(ctx, "%T.%s - end(%#v, %#v = %s)", d, "Confirm", rRes, rErr, rErr) }()

	return d.origin.Confirm(ctx, traceId, cfmId, code)
}

// ---------------------------------------------------------------------------------------------------------------------
