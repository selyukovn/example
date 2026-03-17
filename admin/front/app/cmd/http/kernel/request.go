package kernel

import (
	"encoding/json"
	"net/http"
)

func ClientIp(r *http.Request) string {
	// Следовало бы проверять также `X-Forwarded-For` и список доверенных прокси,
	// но в рамках данного проекта это излишне.
	return r.RemoteAddr
}

func ClientUag(r *http.Request) string {
	return r.Header.Get("User-Agent")
}

func ParseRequestJson[T any](r *http.Request, dataExample T) *T {
	p := &dataExample
	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		return nil
	}
	return p
}
