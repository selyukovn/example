package account

import (
	"context"
	"example/admin/auth/internal/domain/account"
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
func (g *IdGeneratorImplUniqueRandom) Generate(ctx context.Context) (account.Id, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)

	id, err := g.internal.Generate()

	if err != nil {
		return account.IdNil, std.WrapErrorToRuntime(err, g, "Generate")
	}

	return account.Id(id), nil
}

// ---------------------------------------------------------------------------------------------------------------------
