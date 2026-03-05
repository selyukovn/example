package config

import assert "github.com/selyukovn/go-wm-assert"

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Config struct {
	appName                string
	urlSignInWelcome       string
	urlSignInRequest       string
	urlSignInRequestRetry  string
	urlSignInConfirm       string
	urlRedirectToOnSuccess string
	staticBasePath         string
	staticBaseUrl          string
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func New(
	appName string,
	urlSignInWelcome string,
	urlSignInRequest string,
	urlSignInRequestRetry string,
	urlSignInConfirm string,
	urlRedirectToOnSuccess string,
	staticBasePath string,
	staticBaseUrl string,
) *Config {
	assert.Str().NotEmpty().Must(appName)
	assert.Str().NotEmpty().Must(urlSignInWelcome)
	assert.Str().NotEmpty().Must(urlSignInRequest)
	assert.Str().NotEmpty().Must(urlSignInRequestRetry)
	assert.Str().NotEmpty().Must(urlSignInConfirm)
	assert.Str().NotEmpty().Must(urlRedirectToOnSuccess)
	assert.Str().NotEmpty().Must(staticBasePath)
	assert.Str().NotEmpty().Must(staticBaseUrl)

	return &Config{
		appName:                appName,
		urlSignInWelcome:       urlSignInWelcome,
		urlSignInRequest:       urlSignInRequest,
		urlSignInRequestRetry:  urlSignInRequestRetry,
		urlSignInConfirm:       urlSignInConfirm,
		urlRedirectToOnSuccess: urlRedirectToOnSuccess,
		staticBasePath:         staticBasePath,
		staticBaseUrl:          staticBaseUrl,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

// У url'ов могут быть параметры, поэтому это должны быть методы во всех случаях для однообразия.

func (c *Config) AppName() string {
	return c.appName
}

func (c *Config) UrlSignInWelcome() string {
	return c.urlSignInWelcome
}

func (c *Config) UrlSignInRequest() string {
	return c.urlSignInRequest
}

func (c *Config) UrlSignInRequestRetry() string {
	return c.urlSignInRequestRetry
}

func (c *Config) UrlSignInConfirm() string {
	return c.urlSignInConfirm
}

func (c *Config) UrlRedirectToOnSuccess() string {
	return c.urlRedirectToOnSuccess
}

func (c *Config) StaticBasePath() string {
	return c.staticBasePath
}

func (c *Config) StaticBaseUrl() string {
	return c.staticBaseUrl
}

// ---------------------------------------------------------------------------------------------------------------------
