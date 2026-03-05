package event_storage

import (
	"context"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/session"
	_ "github.com/go-sql-driver/mysql"
)

type RepositoryInterface interface {
	// AddAccountCreated
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	//   - std.ErrorRuntime
	AddAccountCreated(ctx context.Context, e account.EventCreated) error

	// AddAccountDeactivated
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	//   - std.ErrorRuntime
	AddAccountDeactivated(ctx context.Context, e account.EventDeactivated) error

	// AddAccountIpWhitelistChanged
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	//   - std.ErrorRuntime
	AddAccountIpWhitelistChanged(ctx context.Context, e account.EventIpWhitelistChanged) error

	// AddSessionCreated
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	//   - std.ErrorRuntime
	AddSessionCreated(ctx context.Context, e session.EventCreated) error

	// AddSessionClosed
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	//   - std.ErrorRuntime
	AddSessionClosed(ctx context.Context, e session.EventClosed) error
}
