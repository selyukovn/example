package cfm

import (
	"context"
	"time"
)

type RepositoryInterface interface {
	// Add
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorAlreadyDone -- если с таким id уже существует
	// 	- std.ErrorRuntime
	Add(ctx context.Context, c *Cfm) error

	// Update
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorRuntime
	Update(ctx context.Context, c *Cfm) error

	// ----------------------------------------------------------------

	// GetById
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorNotFound
	// 	- std.ErrorRuntime
	GetById(ctx context.Context, id Id) (*Cfm, error)

	// GetIdsOfGoingToExpire
	//
	// Cм. Cfm.IsFinished()
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	// 	- std.ErrorRuntime
	GetIdsOfGoingToExpire(ctx context.Context, now time.Time, limit uint) ([]Id, error)
}
