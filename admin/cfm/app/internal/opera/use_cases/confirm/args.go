package confirm

import (
	"context"
	"example/admin/cfm/internal/domain/cfm"
	"example/admin/cfm/internal/domain/cfm/code"
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
	ctx     context.Context
	cfmId   cfm.Id
	cfmCode code.Code
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewArgs
//
// Паникует при нулевых аргументах.
func NewArgs(ctx context.Context, cfmId cfm.Id, cfmCode code.Code) Args {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(cfmId.IsNil())
	assert.FalseMust(cfmCode.IsNil())

	return Args{
		ctx:     ctx,
		cfmId:   cfmId,
		cfmCode: cfmCode,
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

func (a Args) CfmId() cfm.Id {
	return a.cfmId
}

func (a Args) CfmCode() code.Code {
	return a.cfmCode
}

// ---------------------------------------------------------------------------------------------------------------------
