package cfm

import (
	"context"
	"example/admin/cfm/internal/domain/cfm"
	"github.com/selyukovn/go-id/like_uuid"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type IdGeneratorImplUniqueRandom struct {
	internal *like_uuid.IdGeneratorUniqueRandom
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewIdGeneratorImplUniqueRandom() *IdGeneratorImplUniqueRandom {
	return &IdGeneratorImplUniqueRandom{
		internal: like_uuid.NewIdGeneratorUniqueRandom(),
	}
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
func (g *IdGeneratorImplUniqueRandom) Generate(ctx context.Context) (cfm.Id, error) {
	assert.NotNilDeepMust(ctx)

	id, err := g.internal.Generate()

	if err != nil {
		return cfm.IdNil, std.WrapErrorToRuntime(err, g, "Generate")
	}

	return cfm.Id(id), err
}

// ---------------------------------------------------------------------------------------------------------------------
