package request

import (
	"context"
	"example/admin/cfm/internal/domain/cfm"
	"example/admin/cfm/internal/opera/domain_facades"
	goroutiner "github.com/selyukovn/go-routiner"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Command struct {
	grt       *goroutiner.Goroutiner
	cfmDomFac *domain_facades.CfmDomFac
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewCommand
//
// Паникует при нулевых аргументах.
func NewCommand(
	grt *goroutiner.Goroutiner,
	cfmDomFac *domain_facades.CfmDomFac,
) *Command {
	assert.NotNilDeepMust(grt)
	assert.NotNilDeepMust(cfmDomFac)

	return &Command{
		grt:       grt,
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
//   - cfm.ErrorNoAttemptsLeft
//   - cfm.ErrorRequestsFrequency
//   - std.ErrorRuntime
func (c *Command) Execute(args Args) (Result, error) {
	assert.FalseMust(args.IsNil())

	ctx := args.Ctx()
	cfmId := args.CfmId()

	// запоминаем
	cCode, email, canReqAgain, canReqAttemptsLeft, canReqAfter, err := c.cfmDomFac.Request(ctx, cfmId)
	switch err.(type) {
	case nil:
	case std.ErrorNotFound,
		cfm.ErrorFinished,
		cfm.ErrorNoAttemptsLeft,
		cfm.ErrorRequestsFrequency:
		return resultError(err)
	case std.ErrorRuntime:
		return resultError(std.WrapErrorToRuntime(err, c, "Execute", "Request"))
	default:
		panic(err)
	}

	// отправляем
	c.grt.SingleAsync(ctx, func(ctx context.Context) error {
		err = c.cfmDomFac.SendToEmail(ctx, cCode, email)
		switch err.(type) {
		case nil:
		case std.ErrorRuntime:
			return std.WrapErrorToRuntime(err, c, "Execute", "SendToEmail")
		default:
			panic(err)
		}
		return nil
	})

	if !canReqAgain {
		return resultSuccessLast(cfmId)
	}

	return resultSuccessCanAgain(cfmId, canReqAttemptsLeft, canReqAfter)
}

// ---------------------------------------------------------------------------------------------------------------------
