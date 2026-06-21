package kernel

import (
	"errors"
	assert "github.com/selyukovn/go-wm-assert"
	"net/http"
	"strings"
	"time"
)

func CookieHasSessId(r *http.Request) bool {
	return CookieGetSessId(r) != ""
}

func CookieGetSessId(r *http.Request) string {
	assert.NotNilDeepMust(r)

	assert.TrueMust(config.isSet)

	c, err := r.Cookie(config.sessionCookieName)
	if errors.Is(err, http.ErrNoCookie) {
		return ""
	}
	return c.Value
}

func CookieSetSessId(w http.ResponseWriter, sessId string, sessExpAt time.Time) {
	assert.Str().NotEmpty().Must(sessId)
	assert.Time().NotZero().Must(sessExpAt)

	assert.TrueMust(config.isSet)

	isHttps, domain := func(baseUrl string) (bool, string) {
		parts := strings.Split(baseUrl, ":")
		return parts[0] == "https", strings.TrimPrefix(parts[1], "//")
	}(config.baseUrl)

	http.SetCookie(w, &http.Cookie{
		Name:     config.sessionCookieName,
		Value:    sessId,
		Expires:  sessExpAt,
		Domain:   domain,
		Path:     "/",
		HttpOnly: true,
		Secure:   isHttps,
	})
}

func CookieUnsetSessId(w http.ResponseWriter) {
	CookieSetSessId(w, "-", time.Now().Add(-time.Hour))
}
