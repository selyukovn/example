package middlewares

import (
	"context"
	"example/admin/cfm/internal/api/grpc/kernel"
	"github.com/selyukovn/example_gopkg/processing"
	"github.com/selyukovn/go-std/logger"
	"google.golang.org/grpc"
)

func TraceGet() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		operationId, ok := kernel.MetadataKeyFirst(ctx, "x-operation-id")
		if !ok || operationId == "" {
			return nil, kernel.ErrorInvalidArgument("x-operation-id header not found")
		}

		ctx = processing.EnrichCtx(ctx, operationId)
		ctx = logger.AddAttrToCtx(ctx, "processing.OperationId", processing.OperationId(ctx))

		// todo : trace... -- OpenTelemetry?

		return handler(ctx, req)
	}
}
