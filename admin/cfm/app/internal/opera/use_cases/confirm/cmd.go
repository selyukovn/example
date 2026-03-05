package confirm

import (
	"example/admin/cfm/internal/domain/cfm"
	"example/admin/cfm/internal/opera/components"
	"example/admin/cfm/internal/opera/domain_facades"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Command struct {
	logger    components.LoggerInterface
	cfmDomFac *domain_facades.CfmDomFac
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewCommand
//
// Паникует при нулевых аргументах.
func NewCommand(
	logger components.LoggerInterface,
	cfmDomFac *domain_facades.CfmDomFac,
) *Command {
	assert.NotNilDeepMust(logger)
	assert.NotNilDeepMust(cfmDomFac)

	return &Command{
		logger:    logger,
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
func (c *Command) Execute(args Args) (Result, error) {
	assert.FalseMust(args.IsNil())

	ctx := args.Ctx()
	cfmId := args.CfmId()
	cfmCode := args.CfmCode()

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
