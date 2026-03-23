package sign_in_confirm

import (
	"context"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/cfm"
	"example/admin/auth/internal/opera/domain_facades"
	goroutiner "github.com/selyukovn/go-routiner"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Command struct {
	grt          goroutiner.Goroutiner
	accDomFac    domain_facades.AccountDomFac
	actReqDomFac domain_facades.ActionRequestDomFac
	cfmService   cfm.ServiceInterface
	sessDomFac   domain_facades.SessionDomFac
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewCommand
//
// Паникует при нулевых аргументах.
func NewCommand(
	grt goroutiner.Goroutiner,
	accDomFac domain_facades.AccountDomFac,
	actReqDomFac domain_facades.ActionRequestDomFac,
	cfmService cfm.ServiceInterface,
	sessDomFac domain_facades.SessionDomFac,
) Command {
	assert.NotZeroMust(grt)
	assert.Cmp[domain_facades.AccountDomFac]().NotEq(domain_facades.AccountDomFacNil).Must(accDomFac)
	assert.Cmp[domain_facades.ActionRequestDomFac]().NotEq(domain_facades.ActionRequestDomFacNil).Must(actReqDomFac)
	assert.NotNilDeepMust(cfmService)
	assert.Cmp[domain_facades.SessionDomFac]().NotEq(domain_facades.SessionDomFacNil).Must(sessDomFac)

	return Command{
		grt:          grt,
		accDomFac:    accDomFac,
		actReqDomFac: actReqDomFac,
		cfmService:   cfmService,
		sessDomFac:   sessDomFac,
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
//   - std.ErrorNotFound -- SignIn не найден
//   - account.ErrorDeactivated
//   - account.ErrorIpWhitelist
//   - cfm.ErrorFinished
//   - std.ErrorUnprocessable -- если уже есть сессия
//   - std.ErrorRuntime
func (c Command) Execute(args Args) (Result, error) {
	assert.FalseMust(args.IsNil())

	ctx := args.Ctx()
	cl := args.Client()
	signInId := args.SignInId()
	code := args.Code()

	// находим sign in
	accId, cfmId, err := c.actReqDomFac.CheckSignIn(ctx, signInId)
	switch err.(type) {
	case nil:
	case std.ErrorNotFound:
		return resultError(err)
	case std.ErrorRuntime:
		return resultError(std.WrapErrorToRuntime(err, c, "actReqDomFac", "CheckSignIn"))
	default:
		panic(err)
	}

	// Проверяем доступ и параллельно наличие сессии.
	// Если есть созданная сессия, то нет смысла обращаться к конфирмации.
	errs := c.grt.
		Batch(ctx).
		Add(func(ctx context.Context) error {
			err := c.accDomFac.CheckAccess(ctx, cl, accId)
			switch err.(type) {
			case nil:
			case account.ErrorDeactivated, account.ErrorIpWhitelist:
				return err
			case std.ErrorNotFound, std.ErrorRuntime:
				// ErrorNotFound -- тоже бага: id есть, а аккаунта нет
				return std.WrapErrorToRuntime(err, c, "accDomFac", "CheckAccess")
			default:
				panic(err)
			}
			return nil
		}).
		Add(func(ctx context.Context) error {
			isSessionExist, err := c.sessDomFac.HasBySignInRequest(ctx, signInId)
			switch err.(type) {
			case nil:
			case std.ErrorRuntime:
				logger.ErrorFf(ctx, "%T не удалось проверить сессию для SignIn %q: %#v", c, signInId, err)
				// Нет return'а, потому что ветвь сценария не является основной -- лога достаточно.
			default:
				panic(err)
			}
			if isSessionExist {
				// todo : возможно, стоит возвращать cfm.ErrorFinished вместо std.ErrorUnprocessable, т.к. суть одна, но...
				return std.NewErrorUnprocessableFf("Сессия уже создана из SignIn-запроса %q", signInId)
			}
			return nil
		}).
		Wait()

	// Ошибку доступа (при наличии) обязательно нужно отдать первой,
	// иначе технически можно будет сканировать наличие сессий.
	if errs[0] != nil {
		return resultError(errs[0])
	} else if errs[1] != nil {
		return resultError(errs[1])
	}

	// подтверждаем конфирмацию
	cfmRes, err := c.cfmService.Confirm(ctx, cfmId, code)
	switch err.(type) {
	case nil:
	case std.ErrorNotFound:
		return resultError(std.WrapErrorToRuntime(err, c, "cfmService"))
	case cfm.ErrorFinished:
		return resultError(err)
	case std.ErrorUnprocessable, std.ErrorRuntime:
		// ErrorUnprocessable -- не может быть, чтобы тут конфирмация еще не была запрошена -- бага
		return resultError(std.WrapErrorToRuntime(err, c, "cfmService"))
	default:
		panic(err)
	}

	// неудачная попытка подтверждения
	// Последняя или нет -- уже фронт разберется.
	if !cfmRes.IsFinishedAsPassed() {
		return resultFail(cfmRes.FailsLeft())
	}

	// создаем сессию
	sessId, sessExpAt, err := c.sessDomFac.Create(ctx, cl, accId, signInId)
	switch err.(type) {
	case nil:
	case std.ErrorAlreadyDone, std.ErrorRuntime:
		// ErrorAlreadyDone -- тоже бага, т.к. проверяли же наличие сессии ранее.
		return resultError(std.WrapErrorToRuntime(err, c, "sessDomFac", "Create"))
	default:
		panic(err)
	}

	return resultSuccess(sessId, sessExpAt)
}

// ---------------------------------------------------------------------------------------------------------------------
