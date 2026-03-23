package session

import (
	"context"
	"example/admin/auth/internal/domain/session"
	"github.com/selyukovn/go-id/like_uuid"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ session.IdGeneratorInterface = IdGeneratorImplUniqueRandom{}

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
func (g IdGeneratorImplUniqueRandom) Generate(ctx context.Context) (session.Id, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)

	id, err := like_uuid.GenerateUniqueRandom()

	if err != nil {
		return session.IdNil, std.WrapErrorToRuntime(err, g, "Generate")
	}

	return session.Id(id), nil
}

// ---------------------------------------------------------------------------------------------------------------------
