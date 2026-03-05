package sign_out

import (
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/session"
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
	assert.Cmp[components.LoggerInterface]().NotEq(nil).Must(logger)
	assert.Cmp[*domain_facades.AccountDomFac]().NotEq(nil).Must(accDomFac)
	assert.Cmp[*domain_facades.SessionDomFac]().NotEq(nil).Must(sessDomFac)

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
//   - std.ErrorAlreadyDone -- если уже закрыта
//   - std.ErrorRuntime
func (c *Command) Execute(args Args) error {
	assert.FalseMust(args.IsNil())

	ctx := args.Ctx()
	cl := args.Client()
	sessId := args.SessId()

	var err error

	// находим аккаунт ид
	accId, err := c.sessDomFac.GetAccId(ctx, sessId)
	switch err.(type) {
	case nil:
	case std.ErrorNotFound:
		return err
	case std.ErrorRuntime:
		return std.WrapErrorToRuntime(err, c, "Execute", "sessDomFac", "GetAccId")
	default:
		panic(err)
	}

	// проверяем аккаунт
	err = c.accDomFac.CheckAccess(ctx, cl, accId)
	switch err.(type) {
	case nil:
	case account.ErrorDeactivated, account.ErrorIpWhitelist:
		return err
	case std.ErrorNotFound, std.ErrorRuntime:
		// NotFound -- тоже бага: id есть, а аккаунта нет
		return std.WrapErrorToRuntime(err, c, "execNormalCase", "accDomFac", "CheckAccess")
	default:
		panic(err)
	}

	// завершаем сессию
	err = c.sessDomFac.Close(ctx, sessId)
	switch err.(type) {
	case nil:
	case std.ErrorNotFound:
		return std.WrapErrorToRuntime(err, c, "Execute", "session")
	case session.ErrorClosed:
		return std.NewErrorAlreadyDoneFf("Сессия %q уже закрыта: %v", sessId, err)
	case std.ErrorRuntime:
		return std.WrapErrorToRuntime(err, c, "Execute", "session")
	default:
		panic(err)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
