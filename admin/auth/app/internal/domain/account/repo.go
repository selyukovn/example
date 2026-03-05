package account

import (
	"context"
	"github.com/selyukovn/go-std"
)

type RepositoryInterface interface {
	// Add
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorAlreadyDone -- если с таким id или email уже существует
	// 	- std.ErrorRuntime
	Add(ctx context.Context, a *Account) error

	// Update
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorRuntime
	Update(ctx context.Context, a *Account) error

	// GetById
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorNotFound
	// 	- std.ErrorRuntime
	GetById(ctx context.Context, id Id) (*Account, error)

	// GetByEmail
	//
	// Ошибки:
	// 	- std.ErrorNotFound
	// 	- std.ErrorRuntime
	GetByEmail(ctx context.Context, email std.Email) (*Account, error)

	// GetEmailById
	//
	// Ошибки:
	// 	- std.ErrorNotFound
	// 	- std.ErrorRuntime
	GetEmailById(ctx context.Context, id Id) (std.Email, error)
}
