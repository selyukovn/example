package code

import (
	"context"
	"errors"
	"example/admin/cfm/internal/domain/cfm/code"
	"fmt"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"golang.org/x/crypto/bcrypt"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ code.HasherInterface = HasherImplBcrypt{}
var _ code.HasherInterface = HasherImplBcrypt10{}

// Значения хешей будут отличаться при разных cost'ах bcrypt'а,
// а для клиента это будет выглядеть как использование другого алгоритма -- т.е. интуитивно другую реализацию.
// Поэтому важно зафиксировать cost -- в данном случае достаточно default-сложности = 10

type HasherImplBcrypt struct {
	cost int
}

type HasherImplBcrypt10 struct {
	HasherImplBcrypt
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewHasherImplBcrypt10() HasherImplBcrypt10 {
	cost := 10 // !!!
	return HasherImplBcrypt10{HasherImplBcrypt{cost: cost}}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Hash
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (hr HasherImplBcrypt) Hash(ctx context.Context, c code.Code) (code.Hash, error) {
	assert.NotNilDeepMust(ctx)

	cValBytes := []byte(c.String())

	hValBytes, err := bcrypt.GenerateFromPassword(cValBytes, hr.cost)
	if err != nil {
		err = fmt.Errorf("не удалось сгенерировать хеш для %s : %w", c, err)
		return code.HashNil, std.WrapErrorToRuntime(err, hr, "Hash")
	}

	hValStr := string(hValBytes)

	hash, err := code.HashFromString(hValStr)
	if err != nil {
		err = fmt.Errorf("некорректно сгенерированный хеш для %s : %w", c, err)
		return code.HashNil, std.WrapErrorToRuntime(err, hr, "Hash")
	}

	return hash, nil
}

// Compare
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (hr HasherImplBcrypt) Compare(ctx context.Context, c code.Code, h code.Hash) (bool, error) {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(c.IsNil())
	assert.FalseMust(h.IsNil())

	err := bcrypt.CompareHashAndPassword([]byte(h.String()), []byte(c.String()))

	if err == nil {
		return true, nil
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	return false, std.WrapErrorToRuntime(err, hr, "Compare")
}

// ---------------------------------------------------------------------------------------------------------------------
