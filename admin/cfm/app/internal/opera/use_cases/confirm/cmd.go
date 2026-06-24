package confirm

import (
	"context"
	"example/admin/cfm/internal/domain/cfm"
	"example/admin/cfm/internal/domain/cfm/code"
	"example/admin/cfm/internal/opera/domain_facades"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Command struct {
	cfmDomFac domain_facades.CfmDomFac
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewCommand
//
// Паникует при нулевых аргументах.
func NewCommand(
	cfmDomFac domain_facades.CfmDomFac,
) Command {
	assert.Cmp[domain_facades.CfmDomFac]().NotEq(domain_facades.CfmDomFacNil).Must(cfmDomFac)

	return Command{
		cfmDomFac: cfmDomFac,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Execute
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - cfm.ErrorFinished
//   - std.ErrorUnprocessable -- если еще не запрашивалась (request.NewCommand)
//   - std.ErrorRuntime
func (c Command) Execute(ctx context.Context, cfmId cfm.Id, cfmCode code.Code) (Result, error) {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(cfmId.IsNil())
	assert.FalseMust(cfmCode.IsNil())

	finishedAt, isAsPassed, failsLeft, err := c.cfmDomFac.Confirm(ctx, cfmId, cfmCode)
	switch err.(type) {
	case nil:
	case std.ErrorNotFound, cfm.ErrorFinished, std.ErrorUnprocessable:
		return resultError(err)
	case std.ErrorRuntime:
		return resultError(std.WrapErrorToRuntime(err, c, "Execute", "Confirm"))
	default:
		panic(err)
	}

	if !finishedAt.IsZero() {
		return resultSuccessLast(cfmId, finishedAt, isAsPassed)
	}

	return resultSuccessCanAgain(cfmId, failsLeft)
}

// ---------------------------------------------------------------------------------------------------------------------
