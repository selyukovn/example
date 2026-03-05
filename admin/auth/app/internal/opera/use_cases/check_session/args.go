package check_session

import (
	"context"
	"example/admin/auth/internal/domain/client"
	"example/admin/auth/internal/domain/session"
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
	ctx    context.Context
	cl     client.Client
	sessId session.Id
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewArgs
//
// Паникует при нулевых аргументах.
func NewArgs(
	ctx context.Context,
	cl client.Client,
	sessId session.Id,
) Args {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(cl.IsNil())
	assert.FalseMust(sessId.IsNil())

	return Args{
		ctx:    ctx,
		cl:     cl,
		sessId: sessId,
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

func (a Args) SessId() session.Id {
	return a.sessId
}

// ---------------------------------------------------------------------------------------------------------------------
