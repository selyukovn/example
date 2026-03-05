package cfm

import (
	"context"
	"github.com/selyukovn/go-std"
)

type ServiceInterface interface {
	// CreateForEmail
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorRuntime
	CreateForEmail(ctx context.Context, email std.Email) (ServiceResultCreate, error)

	// Request
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorNotFound
	// 	- ErrorFinished
	// 	- ErrorNoAttemptsLeft
	// 	- ErrorRequestsFrequency
	// 	- std.ErrorRuntime
	Request(ctx context.Context, cfmId Id) (ServiceResultRequest, error)

	// Confirm
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorNotFound
	// 	- ErrorFinished
	// 	- std.ErrorUnprocessable -- если не была запрошена
	// 	- std.ErrorRuntime
	Confirm(ctx context.Context, cfmId Id, code Code) (ServiceResultConfirm, error)
}
