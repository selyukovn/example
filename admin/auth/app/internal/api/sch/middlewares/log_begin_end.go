package middlewares

import (
	"context"
	"example/admin/auth/internal/api/sch/kernel"
	"github.com/selyukovn/go-std/logger"
)

func LogBeginEnd(name string) kernel.Middleware {
	return func(next kernel.Handler) kernel.Handler {
		return func(ctx context.Context) {
			ctx = logger.AddAttrToCtx(ctx, "sch.name", name)

			logger.InfoFf(ctx, "Запуск...")
			defer logger.InfoFf(ctx, "Готово")

			next(ctx)
		}
	}
}
