package kernel

import (
	assert "github.com/selyukovn/go-wm-assert"
	"net/http"
	"net/netip"
	"strings"
)

// ---------------------------------------------------------------------------------------------------------------------

const requestIdKey = "kernel.requestId"

func enrichRequest(r *http.Request, requestId string) {
	assert.NotNilDeepMust(r)
	assert.Str().NotEmpty().Must(requestId)

	// Добавление заголовков к полученному запросу кажется костылем, но даже если и так, то это меньшее зло.
	// Интерфейсы методов `enrichRequest(*http.Request)` и `RequestId(*http.Request)` согласованы -- это важнее.
	// Альтернативы:
	// - Enrich(ctx) и RequestId(ctx) -- слишком широко, и с пакетом вяжется только термин RequestId.
	// - r.Context() вместо хедеров использовать нельзя, поскольку контекст может быть извлечен до вызова этого метода
	//   для независимого обогащения с последующим r = r.WithContext() -- добавленный здесь requestId будет потерян.
	r.Header.Add(requestIdKey, requestId)
}

// RequestId -- см. `kernel.RequestIdFromCtx()`
//
// Паникует при нулевых аргументах.
// Паникует, если запрос не был обогащен через `kernel.enrichRequest()`.
func RequestId(r *http.Request) string {
	assert.NotNilDeepMust(r)

	v := r.Header.Get(requestIdKey)

	if v == "" {
		panic("`kernel.RequestId`: похоже, `kernel.enrichRequest` не был вызван")
	}

	return v
}

// ---------------------------------------------------------------------------------------------------------------------

func UserIp(r *http.Request) netip.Addr {
	// При открытом доступе к `gateway` IP следовало бы получать из `r.RemoteAddr`,
	// дополнительно проверяя X-Forwarded-For для случаев использования доверенных прокси.
	// Однако, в текущей схеме все запросы поступают из `front`, который передает клиентский IP как `X-Client-Ip`.
	ip := r.Header.Get("X-Client-Ip")
	ip = strings.Split(ip, ":")[0]
	// `r.RemoteAddr` не может быть некорректным, поэтому Must.
	return netip.MustParseAddr(ip)
}

// UserAgent
//
// Может возвращать пустое значение!
func UserAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}

// ---------------------------------------------------------------------------------------------------------------------
