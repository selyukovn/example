package code

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testDataProvider_Code_correctValues() []string {
	return []string{
		"000000",
		"000001",
		"000011",
		"001100",
		"100000",
		"111111",
		"123456",
		"999999",
	}
}

func testDataProvider_Code_incorrectValues() []string {
	return []string{
		"",
		" ",
		" 0",
		"0 ",
		"&",
		"&=23",
		"23-5",
		"a",
		"aaabbb",
		"aaa999",
		"999aaa",
		"999_999",
		"999 999",
		"0",
		"5",
		"12",
		"321",
		"3214",
		"32145",
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// CodeFromString
// ---------------------------------------------------------------------------------------------------------------------

func Test_CodeFromString(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		tCases := testDataProvider_Code_correctValues()
		for _, tCase := range tCases {
			cc, err := CodeFromString(tCase)
			// no error
			assert.NoError(t, err)
			// not nil
			assert.NotEqual(t, CodeNil, cc)
			assert.False(t, cc.IsNil())
			assert.Equal(t, tCase, cc.String())
		}
	})

	t.Run("incorrect", func(t *testing.T) {
		tCases := testDataProvider_Code_incorrectValues()
		for _, tCase := range tCases {
			cc, err := CodeFromString(tCase)
			// error
			assert.NotNil(t, err)
			// nil
			assert.Equal(t, CodeNil, cc)
			assert.True(t, cc.IsNil())
		}
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// IsNil
// ---------------------------------------------------------------------------------------------------------------------

func Test_Code_IsNil(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		assert.True(t, CodeNil.IsNil())
	})

	t.Run("false", func(t *testing.T) {
		tCases := testDataProvider_Code_correctValues()
		for _, tCase := range tCases {
			cc, _ := CodeFromString(tCase)
			assert.False(t, cc.IsNil())
		}
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// String
// ---------------------------------------------------------------------------------------------------------------------

func Test_Code_String(t *testing.T) {
	tCases := testDataProvider_Code_correctValues()
	for _, tCase := range tCases {
		cc, _ := CodeFromString(tCase)
		assert.Equal(t, tCase, cc.String())
	}
}

// ---------------------------------------------------------------------------------------------------------------------
