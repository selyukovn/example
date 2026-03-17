package kernel

import (
	"net/http"
	"net/netip"
	"strings"
)

func TraceId(r *http.Request) string {
	return r.Header.Get("X-Trace-Id")
}

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
