package processing

import (
	"context"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------

const operationIdCtxKey = "processing.operationId"

// ---------------------------------------------------------------------------------------------------------------------

// EnrichCtx
//
// Паникует при нулевых аргументах.
func EnrichCtx(ctx context.Context, operationId string) context.Context {
	assert.Str().NotEmpty().Must(operationId)

	ctx = context.WithValue(ctx, operationIdCtxKey, operationId)

	return ctx
}

// ---------------------------------------------------------------------------------------------------------------------

// OperationId
//
// Идентификатор выполняемой операции.
// Связывает все действия во всех сервисах в рамках одной операции, инициируемой клиентом.
// Назначается в точке входа, но может быть назначен клиентом для соблюдения идемпотентности.
//
// Паникует, если контекст не обогащен через `processing.EnrichCtx()`.
func OperationId(ctx context.Context) string {
	v := ctx.Value(operationIdCtxKey)

	if v == nil {
		panic("`processing.OperationId`: похоже, `processing.EnrichCtx` не был вызван")
	}

	return v.(string)
}

// ---------------------------------------------------------------------------------------------------------------------
