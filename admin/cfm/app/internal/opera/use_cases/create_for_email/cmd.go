package create_for_email

import (
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
//   - std.ErrorRuntime
func (c *Command) Execute(args Args) (Result, error) {
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
