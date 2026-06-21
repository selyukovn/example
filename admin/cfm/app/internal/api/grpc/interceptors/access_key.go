package interceptors

import (
	"context"
	"example/admin/cfm/internal/api/grpc/kernel"
	assert "github.com/selyukovn/go-wm-assert"
	"google.golang.org/grpc"
)

func NewAccessKey(apiKey string) grpc.UnaryServerInterceptor {
	assert.Str().NotEmpty().Must(apiKey)

	expectedHeaderValue := "Bearer " + apiKey

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		actualHeaderValue, ok := kernel.MetadataKeyFirst(ctx, "authorization")

		if !ok || actualHeaderValue == "" {
			return nil, kernel.ErrorUnauthenticated()
		} else if actualHeaderValue != expectedHeaderValue {
			return nil, kernel.ErrorPermissionDenied()
		}

		return handler(ctx, req)
	}
}
