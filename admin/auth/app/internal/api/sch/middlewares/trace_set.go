package middlewares

import (
	"context"
	"example/admin/auth/internal/api/sch/kernel"
	"github.com/google/uuid"
	"github.com/selyukovn/example_gopkg/processing"
	"github.com/selyukovn/go-std/logger"
	"strings"
)

func TraceSet() kernel.Middleware {
	return func(next kernel.Handler) kernel.Handler {
		return func(ctx context.Context) {
			operationId := strings.Replace(uuid.Must(uuid.NewRandom()).String(), "-", "", -1)

			ctx = processing.EnrichCtx(ctx, operationId)

			ctx = logger.AddAttrToCtx(ctx, "processing.OperationId", operationId)

			// todo : trace... -- OpenTelemetry?

			next(ctx)
		}
	}
}
