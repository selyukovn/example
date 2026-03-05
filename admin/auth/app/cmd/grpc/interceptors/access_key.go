package interceptors

import (
	"context"
	"example/admin/auth/cmd/grpc/helpers"
	assert "github.com/selyukovn/go-wm-assert"
	"google.golang.org/grpc"
)

func NewAccessKey(apiKey string) grpc.UnaryServerInterceptor {
	assert.Str().NotEmpty().Must(apiKey)

	expectedHeaderValue := "Bearer " + apiKey

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		actualHeaderValue, ok := helpers.GrpcMetadataKeyFirst(ctx, "authorization")

		if !ok || actualHeaderValue == "" {
			return nil, helpers.ErrorUnauthenticated()
		} else if actualHeaderValue != expectedHeaderValue {
			return nil, helpers.ErrorPermissionDenied()
		}

		return handler(ctx, req)
	}
}
