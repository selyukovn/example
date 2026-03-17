package http

import (
	"fmt"
	assert "github.com/selyukovn/go-wm-assert"
	"os"
	"strconv"
)

type Env = struct {
	AppName           string
	BaseUrl           string
	SessionCookieName string
	ApiBaseUrl        string
}

func LoadEnv() Env {
	fnMustNGet := func(eVarName string) string {
		v := os.Getenv(eVarName)
		assert.Str().NotEmpty().Must(v, fmt.Sprintf("env: %q не должен быть пустым", eVarName))
		return v
	}

	fnEnvBool := func(eVarName string) bool {
		str := fnMustNGet(eVarName)
		assert.Str().In([]string{"0", "1"}).Must(str, fmt.Sprintf("env: %q ожидает 0 или 1", eVarName))
		return str == "1"
	}

	fnEnvPort := func(eVarName string) uint {
		str := fnMustNGet(eVarName)
		i, err := strconv.Atoi(str)
		if err != nil {
			panic(fmt.Errorf("env: %q ждет номер порта, а не %q: %w", eVarName, str, err))
		}
		return uint(i)
	}

	_ = fnEnvBool
	_ = fnEnvPort

	return Env{
		AppName:           fnMustNGet("APP_NAME"),
		BaseUrl:           fnMustNGet("BASE_URL"),
		SessionCookieName: fnMustNGet("SESSION_COOKIE_NAME"),
		ApiBaseUrl:        fnMustNGet("API_BASEURL"),
	}
}
