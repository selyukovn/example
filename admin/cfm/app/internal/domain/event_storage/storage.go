package event_storage

import (
	"context"
	"example/admin/cfm/internal/domain/cfm"
	"fmt"
	"github.com/selyukovn/go-events"
	assert "github.com/selyukovn/go-wm-assert"
)

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
func NewStorage(repo RepositoryInterface) *Storage {
	assert.Cmp[RepositoryInterface]().NotEq(nil).Must(repo)

	return &Storage{repo: repo}
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
func (s *Storage) Store(ctx context.Context, evs *event.Collection) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*event.Collection]().NotEq(nil).Must(evs)

	for _, e := range evs.All() {
		switch ev := e.(type) {

		// cfm
		case cfm.EventFinished:
			return s.repo.AddCfmFinished(ctx, ev)

		// not registered
		default:
			panic(fmt.Errorf("%T.Store не знает, как сохранить событие %T", s, ev))
		}
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
