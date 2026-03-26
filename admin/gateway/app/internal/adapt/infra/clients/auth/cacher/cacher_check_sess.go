package cacher

import (
	"bytes"
	"context"
	"encoding/gob"
	infra_clients_auth "example/admin/gateway/internal/infra/clients/auth"
	"fmt"
	"github.com/selyukovn/go-std"
	"net/netip"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// CheckSession
// ---------------------------------------------------------------------------------------------------------------------

const checkSessNegativeTtl = 1 * time.Minute

// checkSessMakeSessionBoxKey
//
// Ключ основан только на сессии для простоты удаления -- см. Cacher.CheckSessionUnsetBySessionId().
func checkSessMakeSessionBoxKey(sessionId string) string {
	return fmt.Sprintf("%s/checkSess|%s", keyPrefixGeneral, sessionId)
}

type checkSessSessionCacheBox = struct {
	MaxExpireAtUnixMilli int64
	Results              map[string]struct {
		ExpireAtUnixMilli int64
		Result            cachedResult
	}
}

// Set
// ---------------------------------------------------------------------------------------------------------------------

func (c Cacher) checkSessSetResult(
	ctx context.Context,
	sessionId string,
	fromIp netip.Addr,
	fromUserAgent string,
	resultType byte,
	resultData []byte,
	ttl time.Duration,
) error {
	// Внимание!
	// `resultData` не может быть `any`, иначе возникнут проблемы с декодированием.
	// Нужно выполнять кодирование типизированного значения в каждом конкретном методе.
	// Поэтому `resultData` здесь уже ждет `[]byte`.

	_m_ := "checkSessSetResult"

	// Внимание!
	// Из-за отсутствия синхронизации при обращении к кешу возникает гонка данных.
	// Теоретически опасный сценарий:
	// - с одной вкладки выполняется запись результата CheckSession в кеш (после получения оригинальных данных).
	// - с другой вкладки выполняется Sign-Out, и весь сессионный кеш удаляется.
	// - запрос с первой вкладки выполняет запись ранее полученных данных в кеш.
	// Итого: несмотря на выполненный Sign-Out, сессия остается активной на весь ttl.
	// Вероятность такого сценария мала, но ненулевая.
	// В качестве костыля можно было бы использовать небольшой ttl вместо полного времени жизни сессии,
	// но лучше найти хорошее решение.
	// todo : требуется распределенная блокировка.
	// todo : синхронизация также позволит оптимизировать хранение, разбив cacheBox на отдельные ключи.

	cacheBoxKey := checkSessMakeSessionBoxKey(sessionId)

	var cacheBox checkSessSessionCacheBox
	cacheBoxEncoded, err := c.cache.Get(ctx, cacheBoxKey)
	switch err.(type) {
	case nil:
		if err = gob.NewDecoder(bytes.NewBuffer(cacheBoxEncoded)).Decode(&cacheBox); err != nil {
			return std.WrapErrorToRuntime(err, c, _m_, "Get", "Decode")
		}
	case std.ErrorNotFound:
		cacheBox = checkSessSessionCacheBox{
			MaxExpireAtUnixMilli: 0,
			Results: make(map[string]struct {
				ExpireAtUnixMilli int64
				Result            cachedResult
			}),
		}
	case std.ErrorRuntime:
		return std.WrapErrorToRuntime(err, c, _m_, "Get")
	default:
		panic(err)
	}

	now := time.Now()
	resultExpireAtUnixMilli := now.Add(ttl).UnixMilli()

	cacheBox.Results[fromIp.String()+fromUserAgent] = struct {
		ExpireAtUnixMilli int64
		Result            cachedResult
	}{
		ExpireAtUnixMilli: resultExpireAtUnixMilli,
		Result: cachedResult{
			Data: resultData,
			Type: resultType,
		},
	}
	cacheBox.MaxExpireAtUnixMilli = max(resultExpireAtUnixMilli, cacheBox.MaxExpireAtUnixMilli)
	cacheBoxTtl := time.UnixMilli(cacheBox.MaxExpireAtUnixMilli).Sub(now)

	bb := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(bb).Encode(cacheBox); err != nil {
		return std.WrapErrorToRuntime(err, c, _m_, "Set", "Encode")
	}
	cacheBoxEncoded = bb.Bytes()

	err = c.cache.Set(ctx, cacheBoxKey, cacheBoxEncoded, cacheBoxTtl)
	if err != nil {
		return std.WrapErrorToRuntime(err, c, _m_, "Set")
	}

	return nil
}

func (c Cacher) CheckSessionSetSuccess(
	ctx context.Context,
	sessionId string,
	fromIp netip.Addr,
	fromUserAgent string,
	rRes infra_clients_auth.CheckSessionResult,
) error {
	ttl := rRes.SessionExpireAt.Sub(time.Now())

	bb := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(bb).Encode(rRes); err != nil {
		return err
	}

	return c.checkSessSetResult(
		ctx,
		sessionId,
		fromIp,
		fromUserAgent,
		crtCheckSessResult,
		bb.Bytes(),
		ttl,
	)
}

func (c Cacher) CheckSessionSetErrorNotFound(
	ctx context.Context,
	sessionId string,
	fromIp netip.Addr,
	fromUserAgent string,
	rErr std.ErrorNotFound,
) error {
	bb := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(bb).Encode(rErr); err != nil {
		return err
	}
	return c.checkSessSetResult(
		ctx,
		sessionId,
		fromIp,
		fromUserAgent,
		crtCheckSessErrorNotFound,
		bb.Bytes(),
		checkSessNegativeTtl,
	)
}

func (c Cacher) CheckSessionSetErrorValidation(
	ctx context.Context,
	sessionId string,
	fromIp netip.Addr,
	fromUserAgent string,
	rErr infra_clients_auth.ErrorValidation,
) error {
	bb := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(bb).Encode(rErr); err != nil {
		return err
	}
	return c.checkSessSetResult(
		ctx,
		sessionId,
		fromIp,
		fromUserAgent,
		crtCheckSessErrorValidation,
		bb.Bytes(),
		checkSessNegativeTtl,
	)
}

func (c Cacher) CheckSessionSetErrorAccountAccessDenied(
	ctx context.Context,
	sessionId string,
	fromIp netip.Addr,
	fromUserAgent string,
	rErr infra_clients_auth.ErrorAccountAccessDenied,
) error {
	bb := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(bb).Encode(rErr); err != nil {
		return err
	}
	return c.checkSessSetResult(
		ctx,
		sessionId,
		fromIp,
		fromUserAgent,
		crtCheckSessErrorAccountAccessDenied,
		bb.Bytes(),
		checkSessNegativeTtl,
	)
}

func (c Cacher) CheckSessionSetErrorSessionClosed(
	ctx context.Context,
	sessionId string,
	fromIp netip.Addr,
	fromUserAgent string,
	rErr infra_clients_auth.ErrorSessionClosed,
) error {
	bb := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(bb).Encode(rErr); err != nil {
		return err
	}
	return c.checkSessSetResult(
		ctx,
		sessionId,
		fromIp,
		fromUserAgent,
		crtCheckSessErrorSessionClosed,
		bb.Bytes(),
		checkSessNegativeTtl,
	)
}

// Unset
// ---------------------------------------------------------------------------------------------------------------------

// CheckSessionUnsetBySessionId
//
// Ошибки:
//   - std.ErrorRuntime
func (c Cacher) CheckSessionUnsetBySessionId(ctx context.Context, sessionId string) error {
	_m_ := "CheckSessionUnsetBySessionId"

	key := checkSessMakeSessionBoxKey(sessionId)

	err := c.cache.Unset(ctx, key)
	if err != nil {
		return std.WrapErrorToRuntime(err, c, _m_)
	}

	return nil
}

// Get
// ---------------------------------------------------------------------------------------------------------------------

func (c Cacher) CheckSessionGet(
	ctx context.Context,
	sessionId string,
	fromIp netip.Addr,
	fromUserAgent string,
) (
	rRes infra_clients_auth.CheckSessionResult,
	rErr error,
	cErr error,
) {
	_m_ := "CheckSessionGet"

	cacheKey := checkSessMakeSessionBoxKey(sessionId)

	var cacheBox checkSessSessionCacheBox
	cacheBoxEncoded, err := c.cache.Get(ctx, cacheKey)
	switch err.(type) {
	case nil:
		if err = gob.NewDecoder(bytes.NewBuffer(cacheBoxEncoded)).Decode(&cacheBox); err != nil {
			cErr = std.WrapErrorToRuntime(err, c, _m_, "Get", "Decode")
			return
		}
	case std.ErrorNotFound:
		cErr = err
		return
	case std.ErrorRuntime:
		cErr = std.WrapErrorToRuntime(err, c, _m_, "Get")
		return
	default:
		panic(err)
	}

	cacheBoxItem, exists := cacheBox.Results[fromIp.String()+fromUserAgent]
	if !exists || cacheBoxItem.ExpireAtUnixMilli < time.Now().UnixMilli() {
		cErr = std.NewErrorNotFoundFf("кеш не найден или протух")
		return
	}

	result := cacheBoxItem.Result

	if result.Type == crtCheckSessResult {
		if err := gob.NewDecoder(bytes.NewBuffer(result.Data)).Decode(&rRes); err != nil {
			cErr = c.unsetOnDecodeErrorAndMakeErrorRuntime(ctx, cacheKey, err, _m_, "result", "rRes")
		}
	} else if result.Type == crtCheckSessErrorNotFound ||
		result.Type == crtCheckSessErrorValidation ||
		result.Type == crtCheckSessErrorAccountAccessDenied ||
		result.Type == crtCheckSessErrorSessionClosed {
		if err := gob.NewDecoder(bytes.NewBuffer(result.Data)).Decode(&rErr); err != nil {
			cErr = c.unsetOnDecodeErrorAndMakeErrorRuntime(ctx, cacheKey, err, _m_, "result", "rErr")
		}
	} else {
		panic(fmt.Sprintf("%T.%s не знает, как обработать %T с типом %#v", c, _m_, result, result.Type))
	}

	return
}

// ---------------------------------------------------------------------------------------------------------------------
