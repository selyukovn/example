package action_request

import (
	"context"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/cfm"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Const
// ---------------------------------------------------------------------------------------------------------------------

var FactoryNil = Factory{}

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
func NewFactory(idGenerator IdGeneratorInterface) Factory {
	assert.NotNilDeepMust(idGenerator)

	return Factory{
		idGenerator: idGenerator,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// CreateSignIn
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (f Factory) CreateSignIn(
	ctx context.Context,
	accId account.Id,
	cfmId cfm.Id,
	now time.Time,
) (*SignIn, error) {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(accId.IsNil())
	assert.FalseMust(cfmId.IsNil())
	assert.FalseMust(now.IsZero())

	id, err := f.idGenerator.Generate(ctx)
	switch err.(type) {
	case nil:
	case std.ErrorRuntime:
		return nil, std.WrapErrorToRuntime(err, f, "CreateSignIn")
	default:
		panic(err)
	}

	signIn := &SignIn{
		id:          id,
		accId:       accId,
		cfmId:       cfmId,
		requestedAt: now,
	}

	return signIn, nil
}

// ---------------------------------------------------------------------------------------------------------------------
