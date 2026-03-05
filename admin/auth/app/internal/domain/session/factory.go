package session

import (
	"context"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/action_request"
	"example/admin/auth/internal/domain/client"
	"github.com/selyukovn/go-events"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Const
// ---------------------------------------------------------------------------------------------------------------------

const sessionTtl = time.Hour * 24 * 7

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Factory struct {
	idGenerator IdGeneratorInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewFactory
//
// Паникует при нулевых аргументах.
func NewFactory(idGenerator IdGeneratorInterface) *Factory {
	assert.Cmp[IdGeneratorInterface]().NotEq(nil).Must(idGenerator)

	return &Factory{
		idGenerator: idGenerator,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Create
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (f *Factory) Create(
	ctx context.Context,
	accId account.Id,
	signInId action_request.Id,
	cl client.Client,
	now time.Time,
	evs *event.Collection,
) (*Session, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[account.Id]().NotEq(account.IdNil).Must(accId)
	assert.Cmp[action_request.Id]().NotEq(action_request.IdNil).Must(signInId)
	assert.Cmp[client.Client]().NotEq(client.ClientNil).Must(cl)
	assert.Time().NotZero().Must(now)
	assert.Cmp[*event.Collection]().NotEq(nil).Must(evs)

	id, err := f.idGenerator.Generate(ctx)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, f, "Create")
	}

	expireAt := now.Add(sessionTtl)

	s := &Session{
		id:            id,
		accId:         accId,
		signInId:      signInId,
		initialClient: cl,
		initiatedAt:   now,
		expireAt:      expireAt,
		closedAt:      time.Time{},
	}

	evs.Add(NewEventCreated(now, id, accId))

	return s, err
}

// ---------------------------------------------------------------------------------------------------------------------
