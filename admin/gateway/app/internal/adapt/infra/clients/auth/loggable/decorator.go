package loggable

import (
	"context"
	"example/admin/gateway/internal/infra/clients/auth"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
	"net/netip"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ auth.ClientInterface = Decorator{}

type Decorator struct {
	origin auth.ClientInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewDecorator
//
// Паникует при нулевых аргументах.
func NewDecorator(origin auth.ClientInterface) Decorator {
	assert.NotNilDeepMust(origin)

	return Decorator{
		origin: origin,
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
func (d Decorator) SignInRequest(
	ctx context.Context,
	fromIp netip.Addr,
	fromUserAgent string,
	email std.Email,
) (
	rRes auth.SignInRequestResult,
	rErr error,
) {
	logger.InfoFf(ctx, "%T.%s - start(%q, %q, %q)", d, "SignInRequest", fromIp, fromUserAgent, email)
	defer func() { logger.InfoFf(ctx, "%T.%s - end(%#v, %#v = %s)", d, "SignInRequest", rRes, rErr, rErr) }()

	return d.origin.SignInRequest(ctx, fromIp, fromUserAgent, email)
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
func (d Decorator) SignInRequestRetry(
	ctx context.Context,
	fromIp netip.Addr,
	fromUserAgent string,
	signInId string,
) (
	rRes auth.SignInRequestRetryResult,
	rErr error,
) {
	logger.InfoFf(ctx, "%T.%s - start(%q, %q, %q)", d, "SignInRequestRetry", fromIp, fromUserAgent, signInId)
	defer func() {
		logger.InfoFf(ctx, "%T.%s - end(%#v, %#v = %s)", d, "SignInRequestRetry", rRes, rErr, rErr)
	}()

	return d.origin.SignInRequestRetry(ctx, fromIp, fromUserAgent, signInId)
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
func (d Decorator) SignInConfirm(
	ctx context.Context,
	fromIp netip.Addr,
	fromUserAgent string,
	signInId string,
	code string,
) (
	rRes auth.SignInConfirmResult,
	rErr error,
) {
	logger.InfoFf(
		ctx,
		"%T.%s - start(%q, %q, %q, %q)",
		d, "SignInConfirm", fromIp, fromUserAgent, signInId, std.MaskStrNotFirstLast(code),
	)
	defer func() {
		rResMasked := rRes
		rResMasked.SessionId = std.MaskStrNotFirstLast(rResMasked.SessionId)
		logger.InfoFf(ctx, "%T.%s - end(%#v, %#v = %s)", d, "SignInConfirm", rResMasked, rErr, rErr)
	}()

	return d.origin.SignInConfirm(ctx, fromIp, fromUserAgent, signInId, code)
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
func (d Decorator) SignOut(
	ctx context.Context,
	fromIp netip.Addr,
	fromUserAgent string,
	sessionId string,
) (
	rErr error,
) {
	logger.InfoFf(
		ctx,
		"%T.%s - start(%q, %q, %q)",
		d, "SignOut", fromIp, fromUserAgent, std.MaskStrNotFirstLast(sessionId),
	)
	defer func() { logger.InfoFf(ctx, "%T.%s - end(%#v = %s)", d, "SignOut", rErr, rErr) }()

	return d.origin.SignOut(ctx, fromIp, fromUserAgent, sessionId)
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
func (d Decorator) CheckSession(
	ctx context.Context,
	fromIp netip.Addr,
	fromUserAgent string,
	sessionId string,
) (
	rRes auth.CheckSessionResult,
	rErr error,
) {
	logger.InfoFf(
		ctx,
		"%T.%s - start(%q, %q, %q)",
		d, "CheckSession", fromIp, fromUserAgent, std.MaskStrNotFirstLast(sessionId),
	)
	defer func() { logger.InfoFf(ctx, "%T.%s - end(%#v, %#v = %s)", d, "CheckSession", rRes, rErr, rErr) }()

	return d.origin.CheckSession(ctx, fromIp, fromUserAgent, sessionId)
}

// ---------------------------------------------------------------------------------------------------------------------
