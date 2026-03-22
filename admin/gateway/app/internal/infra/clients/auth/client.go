package auth

import (
	"context"
	"github.com/selyukovn/go-std"
	"net/netip"
)

type ClientInterface interface {
	// SignInRequest
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	//   - std.ErrorNotFound
	//   - ErrorValidation
	//   - ErrorAccountAccessDenied
	//   - std.ErrorRuntime
	SignInRequest(
		ctx context.Context,
		fromIp netip.Addr,
		fromUserAgent string,
		email std.Email,
	) (SignInRequestResult, error)

	// SignInRequestRetry
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	//   - std.ErrorNotFound
	//   - ErrorValidation
	//   - ErrorAccountAccessDenied
	//   - ErrorSignInFinished
	//   - ErrorNoAttemptsLeft
	//   - ErrorRequestsFrequency
	//   - std.ErrorUnprocessable -- сессия уже существует
	//   - std.ErrorRuntime
	SignInRequestRetry(
		ctx context.Context,
		fromIp netip.Addr,
		fromUserAgent string,
		signInId string,
	) (SignInRequestRetryResult, error)

	// SignInConfirm
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	//   - std.ErrorNotFound
	//   - ErrorValidation
	//   - ErrorAccountAccessDenied
	//   - ErrorSignInFinished
	//   - std.ErrorUnprocessable
	//   - std.ErrorRuntime
	SignInConfirm(
		ctx context.Context,
		fromIp netip.Addr,
		fromUserAgent string,
		signInId string,
		code string,
	) (SignInConfirmResult, error)

	// SignOut
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	//   - std.ErrorNotFound
	//   - ErrorValidation
	//   - ErrorAccountAccessDenied
	//   - std.ErrorAlreadyDone
	//   - std.ErrorRuntime
	SignOut(
		ctx context.Context,
		fromIp netip.Addr,
		fromUserAgent string,
		sessionId string,
	) error

	// CheckSession
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	//   - std.ErrorNotFound
	//   - ErrorValidation
	//   - ErrorAccountAccessDenied
	//   - ErrorSessionClosed
	//   - std.ErrorRuntime
	CheckSession(
		ctx context.Context,
		fromIp netip.Addr,
		fromUserAgent string,
		sessionId string,
	) (CheckSessionResult, error)
}
