package code

import (
	"context"
	"example/admin/cfm/internal/domain/cfm/code"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"math/rand"
	"strconv"
	"strings"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

// GeneratorImplUintRand1
//
// Rand 1 -- одним случайным значением.
type GeneratorImplUintRand1 struct{}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewGeneratorImplUintRand1() *GeneratorImplUintRand1 {
	return &GeneratorImplUintRand1{}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Generate
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (g *GeneratorImplUintRand1) Generate(ctx context.Context) (code.Code, error) {
	assert.NotNilDeepMust(ctx)

	var rMin uint32 = 0
	var rMax uint32 = 999_999

	u := rMin + rand.Uint32()%(rMax-rMin+1)

	s := strconv.FormatUint(uint64(u), 10)
	s = strings.Repeat("0", 6-len(s)) + s

	cc, err := code.CodeFromString(s)
	if err != nil {
		return code.CodeNil, std.WrapErrorToRuntime(err, g, "Generate")
	}

	return cc, nil
}

// ---------------------------------------------------------------------------------------------------------------------
