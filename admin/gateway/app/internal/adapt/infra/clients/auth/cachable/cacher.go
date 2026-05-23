package cachable

import (
	"context"
	"errors"
	infra_cache "example/admin/gateway/internal/infra/cache"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Cacher struct {
	cache infra_cache.CacheInterface
}

type cachedResult = struct {
	Data []byte
	Type byte
}

// Константы для cachedResult.Type.
//
// iota потребует соблюдения порядка определения констант, что чревато ошибками.
const (
	crtCheckSessResult                   = byte(1)
	crtCheckSessErrorNotFound            = byte(2)
	crtCheckSessErrorValidation          = byte(3)
	crtCheckSessErrorAccountAccessDenied = byte(4)
	crtCheckSessErrorSessionClosed       = byte(5)
)

const keyPrefixGeneral = "clients.auth.cacheable"

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewCacher
//
// Паникует при нулевых аргументах.
func NewCacher(cache infra_cache.CacheInterface) Cacher {
	assert.NotNilDeepMust(cache)

	return Cacher{
		cache: cache,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (c Cacher) unsetOnDecodeErrorAndMakeErrorRuntime(
	ctx context.Context,
	cacheKey string,
	decodeErr error,
	method string,
	extraInfo1 string,
	extraInfo2 string,
) std.ErrorRuntime {
	// Внимание!
	// При хранении значений с достаточно большим ttl во внешнем хранилище (не локальной памяти) возможна ситуация,
	// когда структура данных (код) будет изменена до истечения ttl записи в кеше.
	// В таком случае при декодировании, вероятно, возникнет ошибка несоответствия данных структуре.
	// Удаление записи в кеше приведет в дальнейшем к перезаписи кеша данными с обновленной структурой.
	cUnsetErr := c.cache.Unset(ctx, cacheKey)

	err := errors.Join(decodeErr, cUnsetErr)
	return std.WrapErrorToRuntime(err, c, method, "unsetOnDecodeErrorAndMakeErrorRuntime", extraInfo1, extraInfo2)
}

// ---------------------------------------------------------------------------------------------------------------------

// ... разделены на файлы для удобства ...

// ---------------------------------------------------------------------------------------------------------------------
