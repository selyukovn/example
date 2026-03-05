package client

import (
	"github.com/selyukovn/go-std"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_UserAgentFromString(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		tCases := []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			"Chrome/89.0.4389.90",
			"curl/7.64.1",
			"A1.B-C/D_E(F)G:H",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
		}

		for _, tCase := range tCases {
			ua, err := UserAgentFromString(tCase)
			assert.NoError(t, err)
			assert.NotEqual(t, UserAgentNil, ua)
			assert.Equal(t, ua.String(), tCase)
			assert.False(t, ua.IsNil())
		}
	})

	t.Run("incorrect", func(t *testing.T) {
		tCases := []string{
			"",
			"   ",
			"\t",
			"\n",
			"a",
			"1invalid",
			"invalid@char",
			"toolong" + strings.Repeat("g", 285),
			"trailing-",
			"-leading",
			"space at end ",
		}

		for _, tCase := range tCases {
			ua, err := UserAgentFromString(tCase)
			assert.Error(t, err)
			assert.IsType(t, std.ErrorValidation{}, err)
			assert.Equal(t, UserAgentNil, ua)
			assert.True(t, ua.IsNil())
		}
	})
}
