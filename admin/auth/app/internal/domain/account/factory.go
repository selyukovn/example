package account

import (
	"context"
	"github.com/selyukovn/go-events"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

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
	assert.NotNilDeepMust(idGenerator)

	return &Factory{
		idGenerator: idGenerator,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Create
//
// Паникует при нулевых аргументах:
//   - ctx
//   - email
//   - now
//   - evs
//
// Ошибки:
//   - std.ErrorRuntime
func (f *Factory) Create(
	ctx context.Context,
	email std.Email,
	ipWhitelist IpWhitelist,
	now time.Time,
	evs *event.Collection,
) (*Account, error) {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(email.IsNil())
	assert.FalseMust(now.IsZero())
	assert.NotNilDeepMust(evs)

	id, err := f.idGenerator.Generate(ctx)
	switch err.(type) {
	case nil:
	case std.ErrorRuntime:
		return nil, std.WrapErrorToRuntime(err, f, "Create")
	default:
		panic(err)
	}

	acc := &Account{
		id:          id,
		email:       email,
		ipWhitelist: ipWhitelist,
	}

	evs.Add(NewEventCreated(now, id, email))

	return acc, nil
}

// ---------------------------------------------------------------------------------------------------------------------
