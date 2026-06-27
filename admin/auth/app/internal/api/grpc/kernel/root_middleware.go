package kernel

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"strings"
)

func RootMiddleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		requestId := strings.Replace(uuid.Must(uuid.NewRandom()).String(), "-", "", -1)

		ctx = enrichCtx(ctx, requestId, info.FullMethod)

		return handler(ctx, req)
	}
}
