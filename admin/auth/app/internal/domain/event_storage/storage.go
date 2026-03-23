package event_storage

import (
	"context"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/session"
	"fmt"
	"github.com/selyukovn/go-events"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Const
// ---------------------------------------------------------------------------------------------------------------------

var StorageNil = Storage{}

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Storage struct {
	repo RepositoryInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewStorage
//
// Паникует при нулевых аргументах.
func NewStorage(repo RepositoryInterface) Storage {
	assert.Cmp[RepositoryInterface]().NotEq(nil).Must(repo)

	return Storage{repo: repo}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Store
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (s Storage) Store(ctx context.Context, evs *event.Collection) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*event.Collection]().NotEq(nil).Must(evs)

	for _, e := range evs.All() {
		switch ev := e.(type) {

		// account
		case account.EventCreated:
			return s.repo.AddAccountCreated(ctx, ev)
		case account.EventDeactivated:
			return s.repo.AddAccountDeactivated(ctx, ev)
		case account.EventIpWhitelistChanged:
			return s.repo.AddAccountIpWhitelistChanged(ctx, ev)

		// session
		case session.EventCreated:
			return s.repo.AddSessionCreated(ctx, ev)
		case session.EventClosed:
			return s.repo.AddSessionClosed(ctx, ev)

		// not registered
		default:
			panic(fmt.Errorf("%T.Store не знает, как сохранить событие %T", s, ev))
		}
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
