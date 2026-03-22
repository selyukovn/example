package helpers

// todo : rename to "kernel" or so...

import (
	"context"
	assert "github.com/selyukovn/go-wm-assert"
	"google.golang.org/grpc/metadata"
)

// ---------------------------------------------------------------------------------------------------------------------
// CTX
// ---------------------------------------------------------------------------------------------------------------------

const grpcCtxUnaryServerInfoRequestIdKey = "grpc.UnaryServerInfo.requestId"
const grpcCtxUnaryServerInfoFullMethodKey = "grpc.UnaryServerInfo.fullMethod"

// AddGrpcInfoToCtx
//
// Паникует при нулевых аргументах.
func AddGrpcInfoToCtx(ctx context.Context, requestId string, fullMethod string) context.Context {
	assert.NotNilDeepMust(ctx)
	assert.Str().NotEmpty().Must(requestId)
	assert.Str().NotEmpty().Must(fullMethod)

	ctx = context.WithValue(ctx, grpcCtxUnaryServerInfoRequestIdKey, requestId)
	ctx = context.WithValue(ctx, grpcCtxUnaryServerInfoFullMethodKey, fullMethod)

	return ctx
}

// ---------------------------------------------------------------------------------------------------------------------

// RequestId
//
// Паникует при нулевых аргументах.
// Паникует, если контекст не обогащен через `helpers.AddGrpcInfoToCtx()`.
func RequestId(ctx context.Context) string {
	assert.NotNilDeepMust(ctx)

	v := ctx.Value(grpcCtxUnaryServerInfoRequestIdKey)

	if v == nil {
		panic("`helpers.RequestId`: похоже, `helpers.AddGrpcInfoToCtx` не был вызван")
	}

	return v.(string)
}

// GrpcFullMethod
//
// Паникует при нулевых аргументах.
// Паникует, если контекст не обогащен через `helpers.AddGrpcInfoToCtx()`.
func GrpcFullMethod(ctx context.Context) string {
	assert.NotNilDeepMust(ctx)

	v := ctx.Value(grpcCtxUnaryServerInfoFullMethodKey)

	if v == nil {
		panic("`helpers.GrpcFullMethod`: похоже, `helpers.AddGrpcInfoToCtx` не был вызван")
	}

	return v.(string)
}

// GrpcMetadataKeyFirst
//
// Возвращает первое значение метаданных с ключом `key`, если таковые имеются.
//
// Паникует при нулевых аргументах.
func GrpcMetadataKeyFirst(ctx context.Context, key string) (string, bool) {
	assert.NotNilDeepMust(ctx)
	assert.Str().NotEmpty().Must(key)

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values, ok := md[key]; ok && len(values) > 0 {
			return values[0], true
		}
	}
	return "", false
}

// ---------------------------------------------------------------------------------------------------------------------
