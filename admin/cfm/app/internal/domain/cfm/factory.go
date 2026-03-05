package cfm

import (
	"context"
	"errors"
	"example/admin/cfm/internal/domain/cfm/code"
	"fmt"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Factory struct {
	idGenerator   IdGeneratorInterface
	codeGenerator code.GeneratorInterface
	codeHasher    code.HasherInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewFactory
//
// Паникует при нулевых аргументах.
func NewFactory(
	idGenerator IdGeneratorInterface,
	codeGenerator code.GeneratorInterface,
	codeHasher code.HasherInterface,
) *Factory {
	assert.NotNilDeepMust(idGenerator)
	assert.NotNilDeepMust(codeGenerator)
	assert.NotNilDeepMust(codeHasher)

	return &Factory{
		idGenerator:   idGenerator,
		codeGenerator: codeGenerator,
		codeHasher:    codeHasher,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// CreateEmailCfm
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (f *Factory) CreateEmailCfm(ctx context.Context, email std.Email, now time.Time) (*Cfm, error) {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(email.IsNil())
	assert.FalseMust(now.IsZero())

	var cfm *Cfm

	id, err := f.idGenerator.Generate(ctx)
	if err != nil {
		err = fmt.Errorf("не удалось создать %T : %w", cfm, err)
		return nil, std.WrapErrorToRuntime(err, f, "CreateEmailCfm")
	}

	// Логично было бы устанавливать срок протухания при первом запросе конфирмации.
	// Однако, протухать нужно начать при создании, чтобы избежать накопления конфирмаций в начальном состоянии.
	// Это проще, чем дополнительные таймеры отдельно для неотправленных конфирмаций или т.п.
	expireAt := now.Add(expirePeriod)
	requests := make([]*CfmRequest, 0)

	cfm = &Cfm{
		id:         id,
		email:      email,
		expireAt:   expireAt,
		finishedAt: time.Time{},
		finishType: FinishNotYet,
		requests:   requests,
		failsMade:  0,
	}

	return cfm, nil
}

// GenerateCodeAndHash
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (f *Factory) GenerateCodeAndHash(ctx context.Context) (code.Code, code.Hash, error) {
	assert.NotNilDeepMust(ctx)

	// todo : возможно, есть смысл объединить генерацию и хеширование в каком-то другом объекте,
	// но пока лучшего места, чем фабрика конфирмации, не найдено.
	cc, err1 := f.codeGenerator.Generate(ctx)
	cch, err2 := f.codeHasher.Hash(ctx, cc)

	if err := errors.Join(err1, err2); err != nil {
		err = fmt.Errorf("не удалось сгенерировать код и/или хеш : %w", err)
		return code.CodeNil, code.HashNil, std.WrapErrorToRuntime(err, f, "GenerateCodeAndHash")
	}

	return cc, cch, nil
}

// ---------------------------------------------------------------------------------------------------------------------
