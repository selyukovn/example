package kernel

import assert "github.com/selyukovn/go-wm-assert"

var config = struct {
	isSet             bool
	baseUrl           string
	sessionCookieName string
}{}

func Configure(
	baseUrl string,
	sessionCookieName string,
) {
	assert.Str().NotEmpty().Must(baseUrl)
	assert.Str().NotEmpty().Must(sessionCookieName)

	config.baseUrl = baseUrl
	config.sessionCookieName = sessionCookieName

	config.isSet = true
}
