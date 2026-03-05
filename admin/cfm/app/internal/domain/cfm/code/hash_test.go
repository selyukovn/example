package code

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func testDataProvider_Hash_correctValues() []string {
	return []string{
		"1",
		"$2a$10$iNWxvnOCfVEy2zRRpwzC1.D2Hkk5tMvvzB6pMuhhP.pJEvCjwcwTm",
		"f30fk9$(KF03k)#(**$$##H#LKLKLKL###",
		strings.Repeat(".", 255),
		strings.Repeat("a", 255),
	}
}

func testDataProvider_Hash_incorrectValues() []string {
	return []string{
		"",
		" ",
		" . ",
		". .",
		strings.Repeat("a", 256),
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// HashFromString
// ---------------------------------------------------------------------------------------------------------------------

func Test_HashFromString(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		tCases := testDataProvider_Hash_correctValues()
		for _, tCase := range tCases {
			cch, err := HashFromString(tCase)
			// no error
			assert.NoError(t, err)
			// not nil
			assert.NotEqual(t, HashNil, cch)
			assert.False(t, cch.IsNil())
			assert.Equal(t, tCase, cch.String())
		}
	})

	t.Run("incorrect", func(t *testing.T) {
		tCases := testDataProvider_Hash_incorrectValues()
		for _, tCase := range tCases {
			cch, err := HashFromString(tCase)
			// error
			assert.NotNil(t, err)
			// nil
			assert.Equal(t, HashNil, cch)
			assert.True(t, cch.IsNil())
		}
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// IsNil
// ---------------------------------------------------------------------------------------------------------------------

func Test_Hash_IsNil(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		assert.True(t, HashNil.IsNil())
	})

	t.Run("false", func(t *testing.T) {
		tCases := testDataProvider_Hash_correctValues()
		for _, tCase := range tCases {
			cch, _ := HashFromString(tCase)
			assert.False(t, cch.IsNil())
		}
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// String
// ---------------------------------------------------------------------------------------------------------------------

func Test_Hash_String(t *testing.T) {
	tCases := testDataProvider_Hash_correctValues()
	for _, tCase := range tCases {
		cch, _ := HashFromString(tCase)
		assert.Equal(t, tCase, cch.String())
	}
}

// ---------------------------------------------------------------------------------------------------------------------
