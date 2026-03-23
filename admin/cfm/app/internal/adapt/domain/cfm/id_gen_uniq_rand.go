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

var _ cfm.IdGeneratorInterface = IdGeneratorImplUniqueRandom{}

type IdGeneratorImplUniqueRandom struct{}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewIdGeneratorImplUniqueRandom() IdGeneratorImplUniqueRandom {
	return IdGeneratorImplUniqueRandom{}
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
func (g IdGeneratorImplUniqueRandom) Generate(ctx context.Context) (cfm.Id, error) {
	assert.NotNilDeepMust(ctx)

	id, err := like_uuid.GenerateUniqueRandom()

	if err != nil {
		return cfm.IdNil, std.WrapErrorToRuntime(err, g, "Generate")
	}

	return cfm.Id(id), err
}

// ---------------------------------------------------------------------------------------------------------------------
