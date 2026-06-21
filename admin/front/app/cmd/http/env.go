package main

import (
	assert "github.com/selyukovn/go-wm-assert"
	"os"
)

type tEnv = struct {
	AppName           string
	BaseUrl           string
	SessionCookieName string
	ApiBaseUrl        string
}

func loadEnv() tEnv {
	return tEnv{
		AppName:           assert.Str().NotEmpty().MustGet(os.Getenv("APP_NAME"), "env: APP_NAME"),
		BaseUrl:           assert.Str().NotEmpty().MustGet(os.Getenv("BASE_URL"), "env: BASE_URL"),
		SessionCookieName: assert.Str().NotEmpty().MustGet(os.Getenv("SESSION_COOKIE_NAME"), "env: SESSION_COOKIE_NAME"),
		ApiBaseUrl:        assert.Str().NotEmpty().MustGet(os.Getenv("API_BASEURL"), "env: API_BASEURL"),
	}
}
