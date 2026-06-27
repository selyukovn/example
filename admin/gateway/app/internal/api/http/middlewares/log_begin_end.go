package middlewares

import (
	"example/admin/gateway/internal/api/http/kernel"
	"github.com/selyukovn/go-std/logger"
	"net/http"
)

func LogBeginEnd() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ctx = logger.AddAttrToCtx(ctx, "kernel.RequestId", kernel.RequestId(r))
			r = r.WithContext(ctx)

			logger.InfoFf(ctx, "Запрос: %s %s", r.Method, r.URL.Path)
			defer func() {
				/* см. */ _ = kernel.RootMiddleware
				status := w.(*kernel.ResponseWriter).Status()
				logger.InfoFf(ctx, "Ответ: %d", status)
			}()

			next.ServeHTTP(w, r)
		})
	}
}
