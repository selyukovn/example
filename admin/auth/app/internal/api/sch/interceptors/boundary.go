package interceptors

import (
	"context"
	"example/admin/auth/internal/api/sch/kernel"
	"github.com/google/uuid"
	"github.com/selyukovn/example_gopkg/processing"
	"github.com/selyukovn/go-std/logger"
	"runtime/debug"
	"strings"
	"time"
)

// NewBoundary
//
// Данный перехватчик должен быть самым внешним!
func NewBoundary() func(string, time.Duration) func(kernel.Handler) kernel.Handler {
	return func(name string, _ time.Duration) func(kernel.Handler) kernel.Handler {
		// Все описанные в данном перехватчике действия слишком связаны между собой,
		// чтобы выделить каждое в отдельный перехватчик.
		return func(next kernel.Handler) kernel.Handler {
			return func(ctx context.Context) {
				// ---------------------------------------------------------------------------------------------------------
				// Резервный перехватчик паники
				// ---------------------------------------------------------------------------------------------------------

				// Теоретически паника может возникнуть до основного перехватчика (например, при обогащении контекста),
				// поэтому обязательно нужен резервный перехватчик, чтобы не завалить весь сервер.

				defer func() {
					if pv := recover(); pv != nil {
						logger.PanicFf(ctx, pv, debug.Stack(), "shc.interceptors.NewBoundary")
					}
				}()

				// ---------------------------------------------------------------------------------------------------------
				// Обогащение
				// ---------------------------------------------------------------------------------------------------------

				operationId := strings.Replace(uuid.Must(uuid.NewRandom()).String(), "-", "", -1)
				ctx = processing.EnrichCtx(ctx, operationId)
				ctx = logger.AddAttrToCtx(ctx, "processing.OperationId", operationId)
				/* см. */ _ = processing.OperationId

				// todo : возможно, нужен аналог RequestId...

				// todo : trace... -- OpenTelemetry?

				ctx = logger.AddAttrToCtx(ctx, "job", name)

				// ---------------------------------------------------------------------------------------------------------
				// Логирование
				// ---------------------------------------------------------------------------------------------------------

				logger.InfoFf(ctx, "%s - запуск...", name)
				defer logger.InfoFf(ctx, "%s - готово!", name)

				// ---------------------------------------------------------------------------------------------------------
				// Основной перехватчик паники
				// ---------------------------------------------------------------------------------------------------------

				// Grpc-сервер, в отличие от net/http-сервера, не умеет сам перехватывать панику.
				// Если панику не перехватить, все приложение рухнет, т.к. будет достигнута вершина стека горутины обработчика.
				// Поскольку не все обработчики могут быть сломаны, имеет смысл панику перехватывать, чтоб не валить весь сервер.

				defer func() {
					if pv := recover(); pv != nil {
						logger.PanicFf(ctx, pv, debug.Stack(), "grpc.interceptors.NewBoundary")
					}
				}()

				// ---------------------------------------------------------------------------------------------------------

				next(ctx)
			}
		}
	}
}
