package middlewares

import (
	"github.com/google/uuid"
	"github.com/selyukovn/example_gopkg/processing"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"net/http"
	"strings"
)

func TraceSet() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			operationId := r.Header.Get("X-Operation-Id")
			operationId = std.Ternary(
				operationId == "",
				strings.Replace(uuid.Must(uuid.NewRandom()).String(), "-", "", -1),
				operationId,
			)

			ctx := r.Context()
			ctx = processing.EnrichCtx(ctx, operationId)
			ctx = logger.AddAttrToCtx(ctx, "processing.OperationId", operationId)

			// todo : trace... -- OpenTelemetry?

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
