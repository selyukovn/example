package session

import (
	"context"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/action_request"
	"time"
)

type RepositoryInterface interface {
	// Add
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorAlreadyDone -- если с таким id или signInId уже существует
	// 	- std.ErrorRuntime
	Add(ctx context.Context, s *Session) error

	// Update
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorRuntime
	Update(ctx context.Context, s *Session) error

	// GetById
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorNotFound
	// 	- std.ErrorRuntime
	GetById(ctx context.Context, id Id) (*Session, error)

	// GetIdsOfGoingToExpire
	//
	// Cм. Session.IsClosed
	//
	// Паникует при нулевых аргументах:
	// 	- ctx
	//
	// Ошибки:
	// 	- std.ErrorRuntime
	GetIdsOfGoingToExpire(ctx context.Context, now time.Time, limit uint) ([]Id, error)

	// HasBySignInRequest
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	//   - std.ErrorRuntime
	HasBySignInRequest(ctx context.Context, signInRequestId action_request.Id) (bool, error)

	// GetAccIdById
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorNotFound
	// 	- std.ErrorRuntime
	GetAccIdById(ctx context.Context, id Id) (account.Id, error)

	// GetAccIdAndExpAtById
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorNotFound
	// 	- std.ErrorRuntime
	GetAccIdAndExpAtById(ctx context.Context, id Id) (account.Id, time.Time, error)
}
