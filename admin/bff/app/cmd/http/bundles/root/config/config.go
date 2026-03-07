package config

import assert "github.com/selyukovn/go-wm-assert"

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Config struct {
	urlForGuest      string
	urlForAuthorized string
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func New(
	urlForGuest string,
	urlForAuthorized string,
) *Config {
	assert.Str().NotEmpty().Must(urlForGuest)
	assert.Str().NotEmpty().Must(urlForAuthorized)

	return &Config{
		urlForGuest:      urlForGuest,
		urlForAuthorized: urlForAuthorized,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (c *Config) UrlForGuest() string {
	return c.urlForGuest
}

func (c *Config) UrlForAuthorized() string {
	return c.urlForAuthorized
}

// ---------------------------------------------------------------------------------------------------------------------
