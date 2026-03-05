package cfm

import (
	"context"
	"github.com/selyukovn/go-std"
)

type ClientInterface interface {
	// CreateForEmail
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	//   - std.ErrorRuntime
	CreateForEmail(ctx context.Context, traceId string, email std.Email) (CreateForEmailResult, error)

	// Request
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	//   - std.ErrorNotFound
	//   - ErrorFinished
	//   - ErrorNoAttemptsLeft
	//   - ErrorRequestsFrequency
	//   - std.ErrorRuntime
	Request(ctx context.Context, traceId string, cfmId string) (RequestResult, error)

	// Confirm
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	//   - std.ErrorNotFound
	//   - ErrorFinished
	//   - std.ErrorUnprocessable -- если не была запрошена
	//   - std.ErrorRuntime
	Confirm(ctx context.Context, traceId string, cfmId string, code string) (ConfirmResult, error)
}
