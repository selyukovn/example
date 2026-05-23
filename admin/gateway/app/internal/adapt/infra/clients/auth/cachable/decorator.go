package cachable

import (
	"context"
	"example/admin/gateway/internal/infra/cache"
	infra_clients_auth "example/admin/gateway/internal/infra/clients/auth"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
	"net/netip"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ infra_clients_auth.ClientInterface = Decorator{}

type Decorator struct {
	origin infra_clients_auth.ClientInterface
	cacher Cacher
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewDecorator
//
// Паникует при нулевых аргументах.
func NewDecorator(origin infra_clients_auth.ClientInterface, cache cache.CacheInterface) Decorator {
	assert.NotNilDeepMust(origin)
	assert.NotZeroMust(cache)

	return Decorator{
		origin: origin,
		cacher: NewCacher(cache),
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Sign In
// ---------------------------------------------------------------------------------------------------------------------

// Внимание!
// Результаты данных запросов не кешируется, поскольку сами по себе запросы довольно редкие,
// выполняют уникальные действия и/или генерируют ответы с уникальными данными.
// При нормальном поведении клиента кеширование даже с малым ttl будет бессмысленным или даже вредным.
// С патологической же клиентурой скорее следует бороться путем ограничения кол-ва запросов и/или иными способами.

func (d Decorator) SignInRequest(ctx context.Context, fromIp netip.Addr, fromUserAgent string, email std.Email) (infra_clients_auth.SignInRequestResult, error) {
	return d.origin.SignInRequest(ctx, fromIp, fromUserAgent, email)
}

func (d Decorator) SignInRequestRetry(ctx context.Context, fromIp netip.Addr, fromUserAgent string, signInId string) (infra_clients_auth.SignInRequestRetryResult, error) {
	return d.origin.SignInRequestRetry(ctx, fromIp, fromUserAgent, signInId)
}

func (d Decorator) SignInConfirm(ctx context.Context, fromIp netip.Addr, fromUserAgent string, signInId string, code string) (infra_clients_auth.SignInConfirmResult, error) {
	return d.origin.SignInConfirm(ctx, fromIp, fromUserAgent, signInId, code)
}

// Sign Out
// ---------------------------------------------------------------------------------------------------------------------

func (d Decorator) SignOut(ctx context.Context, fromIp netip.Addr, fromUserAgent string, sessionId string) error {
	rErr := d.origin.SignOut(ctx, fromIp, fromUserAgent, sessionId)
	if rErr != nil {
		return rErr
	}

	// Внимание!
	// Удаление кеша может происходить по событию `SessionClosed`,
	// но в случае опоздания/потери/ошибке-при-обработке сообщения кеш продолжит отдавать сессию как активную,
	// что может как минимум вводить в заблуждение.
	// Поэтому в первую очередь нужно удалять кеш здесь!
	//
	// Кроме того, например, in-memory-реализация кеша может быть недоступна consumer'у при раздельном запуске.
	if err := d.cacher.CheckSessionUnsetBySessionId(ctx, sessionId); err != nil {
		logger.ErrorFf(ctx, std.WrapErrorToRuntime(err, d, "SignOut", "CheckSessionUnsetBySessionId").Error())
	}

	return nil
}

// Check Session
// ---------------------------------------------------------------------------------------------------------------------

// CheckSession
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - infra_clients_auth.ErrorValidation
//   - infra_clients_auth.ErrorAccountAccessDenied
//   - infra_clients_auth.ErrorSessionClosed
//   - std.ErrorRuntime
func (d Decorator) CheckSession(
	ctx context.Context,
	fromIp netip.Addr,
	fromUserAgent string,
	sessionId string,
) (infra_clients_auth.CheckSessionResult, error) {
	rRes, rErr, cErr := d.cacher.CheckSessionGet(ctx, sessionId, fromIp, fromUserAgent)
	switch cErr.(type) {
	case nil:
		// HIT
		return rRes, rErr
	case std.ErrorNotFound:
		// --> MISS
	case std.ErrorRuntime:
		// ERROR --> LOG & FALLBACK
		logger.ErrorFf(ctx, std.WrapErrorToRuntime(cErr, d, "CheckSession", "CheckSessionGet").Error())
		return d.origin.CheckSession(ctx, fromIp, fromUserAgent, sessionId)
	default:
		panic(cErr)
	}

	// MISS --> SET
	rRes, rErr = d.origin.CheckSession(ctx, fromIp, fromUserAgent, sessionId)
	switch vrErr := rErr.(type) {
	case nil:
		if err := d.cacher.CheckSessionSetSuccess(ctx, sessionId, fromIp, fromUserAgent, rRes); err != nil {
			logger.ErrorFf(ctx, std.WrapErrorToRuntime(err, d, "CheckSession", "CheckSessionSetSuccess").Error())
		}
	case std.ErrorNotFound:
		if err := d.cacher.CheckSessionSetErrorNotFound(ctx, sessionId, fromIp, fromUserAgent, vrErr); err != nil {
			logger.ErrorFf(ctx, std.WrapErrorToRuntime(err, d, "CheckSession", "CheckSessionSetErrorNotFound").Error())
		}
	case infra_clients_auth.ErrorValidation:
		if err := d.cacher.CheckSessionSetErrorValidation(ctx, sessionId, fromIp, fromUserAgent, vrErr); err != nil {
			logger.ErrorFf(ctx, std.WrapErrorToRuntime(err, d, "CheckSession", "CheckSessionSetErrorValidation").Error())
		}
	case infra_clients_auth.ErrorAccountAccessDenied:
		if err := d.cacher.CheckSessionSetErrorAccountAccessDenied(ctx, sessionId, fromIp, fromUserAgent, vrErr); err != nil {
			logger.ErrorFf(ctx, std.WrapErrorToRuntime(err, d, "CheckSession", "CheckSessionSetErrorAccountAccessDenied").Error())
		}
	case infra_clients_auth.ErrorSessionClosed:
		if err := d.cacher.CheckSessionSetErrorSessionClosed(ctx, sessionId, fromIp, fromUserAgent, vrErr); err != nil {
			logger.ErrorFf(ctx, std.WrapErrorToRuntime(err, d, "CheckSession", "CheckSessionSetErrorSessionClosed").Error())
		}
	case std.ErrorRuntime:
		rErr = std.WrapErrorToRuntime(rErr, d, "CheckSession", "origin")
	default:
		panic(rErr)
	}

	return rRes, rErr
}

// ---------------------------------------------------------------------------------------------------------------------
