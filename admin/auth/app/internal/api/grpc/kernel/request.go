package kernel

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

// EnrichCtx
//
// Паникует при нулевых аргументах.
func EnrichCtx(ctx context.Context, requestId string, fullMethod string) context.Context {
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
// Паникует, если контекст не обогащен через `kernel.EnrichCtx()`.
func RequestId(ctx context.Context) string {
	assert.NotNilDeepMust(ctx)

	v := ctx.Value(grpcCtxUnaryServerInfoRequestIdKey)

	if v == nil {
		panic("`kernel.RequestId`: похоже, `kernel.EnrichCtx` не был вызван")
	}

	return v.(string)
}

// FullMethod
//
// Паникует при нулевых аргументах.
// Паникует, если контекст не обогащен через `kernel.EnrichCtx()`.
func FullMethod(ctx context.Context) string {
	assert.NotNilDeepMust(ctx)

	v := ctx.Value(grpcCtxUnaryServerInfoFullMethodKey)

	if v == nil {
		panic("`kernel.FullMethod`: похоже, `kernel.EnrichCtx` не был вызван")
	}

	return v.(string)
}

// MetadataKeyFirst
//
// Паникует при нулевых аргументах.
func MetadataKeyFirst(ctx context.Context, key string) (string, bool) {
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
