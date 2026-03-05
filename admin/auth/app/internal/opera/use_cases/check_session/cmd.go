package check_session

import (
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/opera/components"
	"example/admin/auth/internal/opera/domain_facades"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Command struct {
	logger     components.LoggerInterface
	accDomFac  *domain_facades.AccountDomFac
	sessDomFac *domain_facades.SessionDomFac
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewCommand
//
// Паникует при нулевых аргументах.
func NewCommand(
	logger components.LoggerInterface,
	accDomFac *domain_facades.AccountDomFac,
	sessDomFac *domain_facades.SessionDomFac,
) *Command {
	assert.NotNilDeepMust(logger)
	assert.NotNilDeepMust(accDomFac)
	assert.NotNilDeepMust(sessDomFac)

	return &Command{
		logger:     logger,
		accDomFac:  accDomFac,
		sessDomFac: sessDomFac,
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
//   - std.ErrorNotFound -- сессия не найдена
//   - account.ErrorDeactivated
//   - account.ErrorIpWhitelist
//   - session.ErrorClosed
//   - std.ErrorRuntime
func (c *Command) Execute(args Args) (Result, error) {
	assert.FalseMust(args.IsNil())

	ctx := args.Ctx()
	cl := args.Client()
	sessId := args.SessId()

	var err error

	// находим аккаунт ид
	accId, sessExpAt, err := c.sessDomFac.GetAccIdAndExpAt(ctx, sessId)
	switch err.(type) {
	case nil:
	case std.ErrorNotFound:
		return resultError(err)
	case std.ErrorRuntime:
		return resultError(std.WrapErrorToRuntime(err, c, "Execute", "sessDomFac", "GetAccIdAndExpAt"))
	default:
		panic(err)
	}

	// проверяем аккаунт
	err = c.accDomFac.CheckAccess(ctx, cl, accId)
	switch err.(type) {
	case nil:
	case account.ErrorDeactivated, account.ErrorIpWhitelist:
		return resultError(err)
	case std.ErrorNotFound, std.ErrorRuntime:
		// NotFound -- тоже бага: id есть, а аккаунта нет
		return resultError(std.WrapErrorToRuntime(err, c, "execNormalCase", "accDomFac", "CheckAccess"))
	default:
		panic(err)
	}

	return resultSuccess(accId, sessExpAt)
}

// ---------------------------------------------------------------------------------------------------------------------
