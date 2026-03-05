package sign_in_request_retry

import (
	"context"
	"example/admin/auth/internal/domain/action_request"
	"example/admin/auth/internal/domain/client"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Const
// ---------------------------------------------------------------------------------------------------------------------

var ArgsNil = Args{}

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Args struct {
	ctx      context.Context
	cl       client.Client
	signInId action_request.Id
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewArgs
//
// Паникует при нулевых аргументах.
func NewArgs(ctx context.Context, cl client.Client, signInId action_request.Id) Args {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[client.Client]().NotEq(client.ClientNil).Must(cl)
	assert.Cmp[action_request.Id]().NotEq(action_request.IdNil).Must(signInId)

	return Args{
		ctx:      ctx,
		cl:       cl,
		signInId: signInId,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (a Args) IsNil() bool {
	return a == ArgsNil
}

func (a Args) Ctx() context.Context {
	return a.ctx
}

func (a Args) Client() client.Client {
	return a.cl
}

func (a Args) SignInId() action_request.Id {
	return a.signInId
}

// ---------------------------------------------------------------------------------------------------------------------
