package helpers

import (
	"context"
	assert "github.com/selyukovn/go-wm-assert"
	"google.golang.org/grpc/metadata"
)

// ---------------------------------------------------------------------------------------------------------------------
// CTX
// ---------------------------------------------------------------------------------------------------------------------

const grpcCtxUnaryServerInfoFullMethodKey = "grpc.UnaryServerInfo.fullMethod"

// AddGrpcInfoToCtx
//
// Паникует при нулевых аргументах.
func AddGrpcInfoToCtx(ctx context.Context, fullMethod string) context.Context {
	assert.NotNilDeepMust(ctx)
	assert.Str().NotEmpty().Must(fullMethod)

	ctx = context.WithValue(ctx, grpcCtxUnaryServerInfoFullMethodKey, fullMethod)

	return ctx
}

// ---------------------------------------------------------------------------------------------------------------------

// GrpcInfoFullMethod
//
// Паникует при нулевых аргументах.
func GrpcInfoFullMethod(ctx context.Context) string {
	assert.NotNilDeepMust(ctx)

	return ctx.Value(grpcCtxUnaryServerInfoFullMethodKey).(string)
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
