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

// GrpcFullMethod
//
// Паникует при нулевых аргументах.
func GrpcFullMethod(ctx context.Context) string {
	assert.NotNilDeepMust(ctx)

	return ctx.Value(grpcCtxUnaryServerInfoFullMethodKey).(string)
}

// GrpcMetadataKeyFirst
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

// TODO : TRACER, trace+path, ...

const traceIdCtxKey = "traceId"

func TraceIdSet(ctx context.Context, traceId string) context.Context {
	assert.NotNilDeepMust(ctx)

	return context.WithValue(ctx, traceIdCtxKey, traceId)
}

func TraceIdGet(ctx context.Context) string {
	assert.NotNilDeepMust(ctx)
	return ctx.Value(traceIdCtxKey).(string)
}

// ---------------------------------------------------------------------------------------------------------------------
