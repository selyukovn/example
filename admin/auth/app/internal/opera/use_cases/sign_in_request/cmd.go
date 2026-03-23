package sign_in_request

import (
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/cfm"
	"example/admin/auth/internal/opera/domain_facades"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Command struct {
	accDomFac    domain_facades.AccountDomFac
	actReqDomFac domain_facades.ActionRequestDomFac
	cfmService   cfm.ServiceInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewCommand
//
// Паникует при нулевых аргументах.
func NewCommand(
	accDomFac domain_facades.AccountDomFac,
	actReqDomFac domain_facades.ActionRequestDomFac,
	cfmService cfm.ServiceInterface,
) Command {
	assert.Cmp[domain_facades.AccountDomFac]().NotEq(domain_facades.AccountDomFacNil).Must(accDomFac)
	assert.Cmp[domain_facades.ActionRequestDomFac]().NotEq(domain_facades.ActionRequestDomFacNil).Must(actReqDomFac)
	assert.NotNilDeepMust(cfmService)

	return Command{
		accDomFac:    accDomFac,
		actReqDomFac: actReqDomFac,
		cfmService:   cfmService,
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
//   - std.ErrorNotFound -- аккаунт не найден
//   - account.ErrorDeactivated
//   - account.ErrorIpWhitelist
//   - std.ErrorRuntime
func (c Command) Execute(args Args) (Result, error) {
	assert.FalseMust(args.IsNil())

	ctx := args.Ctx()
	cl := args.Client()
	accEmail := args.Email()

	var err error

	// проверяем аккаунт
	accId, err := c.accDomFac.CanSignIn(ctx, cl, accEmail)
	switch err.(type) {
	case nil:
	case std.ErrorNotFound, account.ErrorDeactivated, account.ErrorIpWhitelist:
		return resultError(err)
	case std.ErrorRuntime:
		return resultError(std.WrapErrorToRuntime(err, c, "Execute", "accDomFac", "CanSignIn"))
	default:
		panic(err)
	}

	// создаем сервисную конфирмацию
	cfmResCreate, err := c.cfmService.CreateForEmail(ctx, accEmail)
	switch err.(type) {
	case nil:
	case std.ErrorRuntime:
		return resultError(std.WrapErrorToRuntime(err, c, "Execute", "cfmService", "CreateForEmail"))
	default:
		panic(err)
	}

	// создаем SignIn-запрос
	actReqId, err := c.actReqDomFac.CreateSignIn(ctx, accId, cfmResCreate.CfmId())
	switch err.(type) {
	case nil:
	case std.ErrorRuntime:
		return resultError(std.WrapErrorToRuntime(err, c, "Execute", "actReqDomFac", "CreateSignIn"))
	default:
		panic(err)
	}

	// запрашиваем подтверждение
	cfmRes, err := c.cfmService.Request(ctx, cfmResCreate.CfmId())
	switch err.(type) {
	case nil:
	case std.ErrorNotFound,
		cfm.ErrorFinished,
		cfm.ErrorNoAttemptsLeft,
		cfm.ErrorRequestsFrequency,
		std.ErrorRuntime:
		// только что создали, поэтому любая из ошибок логики -- бага
		return resultError(std.WrapErrorToRuntime(err, c, "Execute", "cfmService", "Request"))
	default:
		panic(err)
	}

	// только что создали конфирмацию, поэтому нулей не будет
	return resultSuccess(
		actReqId,
		cfmRes.CanReqAttemptsLeft(),
		cfmRes.CanReqAfter(),
		cfmResCreate.ExpireAt(),
	)
}
