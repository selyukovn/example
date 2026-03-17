package auth

import (
	"context"
	"example/admin/gateway/internal/infra/clients/auth"
	"example/admin/gateway/internal/infra/logger"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"net/netip"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ auth.ClientInterface = &DecoratorLoggable{}

type DecoratorLoggable struct {
	origin auth.ClientInterface
	logger *logger.Logger
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewDecoratorLoggable
//
// Паникует при нулевых аргументах.
func NewDecoratorLoggable(origin auth.ClientInterface, logger *logger.Logger) *DecoratorLoggable {
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

// SignInRequest
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - auth.ErrorValidation
//   - auth.ErrorAccountAccessDenied
//   - std.ErrorRuntime
func (d *DecoratorLoggable) SignInRequest(
	ctx context.Context,
	traceId string,
	fromIp netip.Addr,
	fromUserAgent string,
	email std.Email,
) (
	rRes auth.SignInRequestResult,
	rErr error,
) {
	d.logger.CtxInfoFf(ctx, "%T.%s - start(%q, %q, %q)", d, "SignInRequest", fromIp, fromUserAgent, email)
	defer func() { d.logger.CtxInfoFf(ctx, "%T.%s - end(%#v, %#v = %s)", d, "SignInRequest", rRes, rErr, rErr) }()

	return d.origin.SignInRequest(ctx, traceId, fromIp, fromUserAgent, email)
}

// SignInRequestRetry
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - auth.ErrorValidation
//   - auth.ErrorAccountAccessDenied
//   - auth.ErrorSignInFinished
//   - auth.ErrorNoAttemptsLeft
//   - auth.ErrorRequestsFrequency
//   - std.ErrorUnprocessable -- сессия уже существует
//   - std.ErrorRuntime
func (d *DecoratorLoggable) SignInRequestRetry(
	ctx context.Context,
	traceId string,
	fromIp netip.Addr,
	fromUserAgent string,
	signInId string,
) (
	rRes auth.SignInRequestRetryResult,
	rErr error,
) {
	d.logger.CtxInfoFf(ctx, "%T.%s - start(%q, %q, %q)", d, "SignInRequestRetry", fromIp, fromUserAgent, signInId)
	defer func() {
		d.logger.CtxInfoFf(ctx, "%T.%s - end(%#v, %#v = %s)", d, "SignInRequestRetry", rRes, rErr, rErr)
	}()

	return d.origin.SignInRequestRetry(ctx, traceId, fromIp, fromUserAgent, signInId)
}

// SignInConfirm
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - auth.ErrorValidation
//   - auth.ErrorAccountAccessDenied
//   - auth.ErrorSignInFinished
//   - std.ErrorUnprocessable
//   - std.ErrorRuntime
func (d *DecoratorLoggable) SignInConfirm(
	ctx context.Context,
	traceId string,
	fromIp netip.Addr,
	fromUserAgent string,
	signInId string,
	code string,
) (
	rRes auth.SignInConfirmResult,
	rErr error,
) {
	d.logger.CtxInfoFf(
		ctx,
		"%T.%s - start(%q, %q, %q, %q)",
		d, "SignInConfirm", fromIp, fromUserAgent, signInId, std.MaskStrNotFirstLast(code),
	)
	defer func() {
		rResMasked := rRes
		rResMasked.SessionId = std.MaskStrNotFirstLast(rResMasked.SessionId)
		d.logger.CtxInfoFf(ctx, "%T.%s - end(%#v, %#v = %s)", d, "SignInConfirm", rResMasked, rErr, rErr)
	}()

	return d.origin.SignInConfirm(ctx, traceId, fromIp, fromUserAgent, signInId, code)
}

// SignOut
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - auth.ErrorValidation
//   - auth.ErrorAccountAccessDenied
//   - std.ErrorAlreadyDone
//   - std.ErrorRuntime
func (d *DecoratorLoggable) SignOut(
	ctx context.Context,
	traceId string,
	fromIp netip.Addr,
	fromUserAgent string,
	sessionId string,
) (
	rErr error,
) {
	d.logger.CtxInfoFf(
		ctx,
		"%T.%s - start(%q, %q, %q)",
		d, "SignOut", fromIp, fromUserAgent, std.MaskStrNotFirstLast(sessionId),
	)
	defer func() { d.logger.CtxInfoFf(ctx, "%T.%s - end(%#v = %s)", d, "SignOut", rErr, rErr) }()

	return d.origin.SignOut(ctx, traceId, fromIp, fromUserAgent, sessionId)
}

// CheckSession
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - auth.ErrorValidation
//   - auth.ErrorAccountAccessDenied
//   - auth.ErrorSessionClosed
//   - std.ErrorRuntime
func (d *DecoratorLoggable) CheckSession(
	ctx context.Context,
	traceId string,
	fromIp netip.Addr,
	fromUserAgent string,
	sessionId string,
) (
	rRes auth.CheckSessionResult,
	rErr error,
) {
	d.logger.CtxInfoFf(
		ctx,
		"%T.%s - start(%q, %q, %q)",
		d, "CheckSession", fromIp, fromUserAgent, std.MaskStrNotFirstLast(sessionId),
	)
	defer d.logger.CtxInfoFf(ctx, "%T.%s - end(%#v, %#v = %s)", d, "CheckSession", rRes, rErr, rErr)

	return d.origin.CheckSession(ctx, traceId, fromIp, fromUserAgent, sessionId)
}

// ---------------------------------------------------------------------------------------------------------------------
