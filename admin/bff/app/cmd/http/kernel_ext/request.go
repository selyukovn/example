package kernel_ext

import (
	assert "github.com/selyukovn/go-wm-assert"
	"net/http"
	"net/netip"
)

// ---------------------------------------------------------------------------------------------------------------------
// см. rproxy/build/server/nginx.conf
// ---------------------------------------------------------------------------------------------------------------------

func TraceId(r *http.Request) string {
	return r.Header.Get("X-Trace-Id")
}

func UserIp(r *http.Request) netip.Addr {
	ip, err := netip.ParseAddr(r.Header.Get("X-Real-Ip"))
	assert.TrueMust(err == nil)
	return ip
}

func UserAgent(r *http.Request) string {
	uag := r.Header.Get("X-Real-UserAgent")
	assert.Str().NotEmpty().Must(uag)
	return uag
}

// ---------------------------------------------------------------------------------------------------------------------
