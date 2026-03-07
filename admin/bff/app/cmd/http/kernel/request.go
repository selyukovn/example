package kernel

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
	return netip.MustParseAddr(r.Header.Get("X-Real-Ip"))
}

func UserAgent(r *http.Request) string {
	return assert.Str().NotEmpty().MustGet(r.Header.Get("X-Real-UserAgent"))
}

// ---------------------------------------------------------------------------------------------------------------------
