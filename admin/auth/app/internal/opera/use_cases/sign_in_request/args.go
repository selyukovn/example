package sign_in_request

import (
	"context"
	"example/admin/auth/internal/domain/client"
	"github.com/selyukovn/go-std"
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
	ctx   context.Context
	cl    client.Client
	email std.Email
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewArgs
//
// Паникует при нулевых аргументах.
func NewArgs(ctx context.Context, cl client.Client, email std.Email) Args {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[client.Client]().NotEq(client.ClientNil).Must(cl)
	assert.Cmp[std.Email]().NotEq(std.EmailNil).Must(email)

	return Args{
		ctx:   ctx,
		cl:    cl,
		email: email,
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

func (a Args) Email() std.Email {
	return a.email
}

// ---------------------------------------------------------------------------------------------------------------------
