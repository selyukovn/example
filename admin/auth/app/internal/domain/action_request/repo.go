package action_request

import (
	"context"
)

// Использование any с последующим утверждением типов --
// более удобный в применении и более явно отражающий модель за счет различных Get-методов вариант
// по сравнению с "1 тип <--> 1 репозиторий" (в том числе с обобщениями).

type RepositoryInterface interface {
	// Add
	//
	// Паникует, если:
	// 	- ctx == nil
	// 	- actReq == nil или не входит в набор типов: *SignIn
	//
	// Ошибки:
	// 	- std.ErrorAlreadyDone -- если с таким id уже существует
	// 	- std.ErrorRuntime
	Add(ctx context.Context, actReq any) error

	// Update
	//
	// Паникует, если:
	// 	- ctx == nil
	// 	- actReq == nil или не входит в набор типов: *SignIn
	//
	// Ошибки:
	// 	- std.ErrorRuntime
	Update(ctx context.Context, actReq any) error

	// GetSignIn
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorNotFound
	// 	- std.ErrorRuntime
	GetSignIn(ctx context.Context, actReqId Id) (*SignIn, error)
}
