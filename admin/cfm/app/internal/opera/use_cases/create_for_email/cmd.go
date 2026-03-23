package create_for_email

import (
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
//   - std.ErrorRuntime
func (c Command) Execute(args Args) (Result, error) {
	assert.FalseMust(args.IsNil())

	ctx := args.Ctx()
	email := args.Email()

	cfmId, expireAt, err := c.cfmDomFac.CreateForEmail(ctx, email)
	switch err.(type) {
	case nil:
	case std.ErrorRuntime:
		return resultError(std.WrapErrorToRuntime(err, c, "Execute"))
	default:
		panic(err)
	}

	return resultSuccess(cfmId, email, expireAt)
}

// ---------------------------------------------------------------------------------------------------------------------
