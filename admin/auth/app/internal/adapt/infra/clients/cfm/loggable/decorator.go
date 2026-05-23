package loggable

import (
	"context"
	"example/admin/auth/internal/infra/clients/cfm"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ cfm.ClientInterface = Decorator{}

type Decorator struct {
	origin cfm.ClientInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewDecorator
//
// Паникует при нулевых аргументах.
func NewDecorator(origin cfm.ClientInterface) Decorator {
	assert.NotNilDeepMust(origin)

	return Decorator{
		origin: origin,
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
func (d Decorator) CreateForEmail(ctx context.Context, email std.Email) (
	rRes cfm.CreateForEmailResult,
	rErr error,
) {
	logger.InfoFf(ctx, "%T.%s - start(%q)", d, "CreateForEmail", email)
	defer func() { logger.InfoFf(ctx, "%T.%s - end(%#v, %#v = %s)", d, "CreateForEmail", rRes, rErr, rErr) }()

	return d.origin.CreateForEmail(ctx, email)
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
func (d Decorator) Request(ctx context.Context, cfmId string) (
	rRes cfm.RequestResult,
	rErr error,
) {
	logger.InfoFf(ctx, "%T.%s - start(%q)", d, "Request", cfmId)
	defer func() { logger.InfoFf(ctx, "%T.%s - end(%#v, %#v = %s)", d, "Request", rRes, rErr, rErr) }()

	return d.origin.Request(ctx, cfmId)
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
func (d Decorator) Confirm(ctx context.Context, cfmId string, code string) (
	rRes cfm.ConfirmResult,
	rErr error,
) {
	logger.InfoFf(ctx, "%T.%s - start(%q, %q)", d, "Confirm", cfmId, std.MaskStrNotFirstLast(code))
	defer func() { logger.InfoFf(ctx, "%T.%s - end(%#v, %#v = %s)", d, "Confirm", rRes, rErr, rErr) }()

	return d.origin.Confirm(ctx, cfmId, code)
}

// ---------------------------------------------------------------------------------------------------------------------
