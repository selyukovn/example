package http

import (
	"context"
	"example/admin/gateway/cmd/http/bundles/auth"
	"example/admin/gateway/cmd/http/components/monitoring"
	"example/admin/gateway/cmd/http/components/processing"
	"example/admin/gateway/cmd/http/components/security"
	"example/admin/gateway/cmd/http/container"
	"example/admin/gateway/cmd/http/kernel"
	"github.com/google/uuid"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"net/http"
	"runtime/debug"
	"strings"
)

func registerRoutes(
	mux *http.ServeMux,
	ctr *container.Container,
) {
	sec := security.New(ctr)

	middlewares := []func(http.Handler) http.Handler{
		_boundaryMiddleware(ctr),
		monitoring.MetricsMiddleware(),
		_securityMiddleware(sec),
	}

	auth.Register(mux, middlewares, ctr, sec)

	mux.Handle("/", http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})))
}

// Данный перехватчик должен быть самым внешним!
func _boundaryMiddleware(ctr *container.Container) func(http.Handler) http.Handler {
	// Все описанные в данном перехватчике действия слишком связаны между собой,
	// чтобы выделить каждое в отдельный перехватчик.
	// Например, логирование статуса ответа не имеет смысла без operationId из обогащенного контекста,
	// а значит обогащение контекста обязано происходить до логирования статуса ответа.
	// Но статус ответа может быть изменен во внешних перехватчиках, что приведет к расхождению с уже записанными логами.
	// Кроме того, перехват паники для корректного ее логирования потребует ретрансляции и двух точек обработки,
	// что также может привести к нарушению согласованности кодов ответа при использовании отдельных перехватчиков.
	//
	// В net/http-сервере есть перехват паники обработчика -- результатом будет разрыв соединения.
	// Это приведет к ответу с кодом 502, несмотря на то, что уже код 500 записан в логи.
	// Поэтому нужно также, как и в случае с grpc, обязательно перехватывать панику до попадания в сервер.
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var fallbackPanicCtx *context.Context

			fnResponseOnPanic := func(w http.ResponseWriter) {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			// ---------------------------------------------------------------------------------------------------------
			// Резервный перехватчик паники
			// ---------------------------------------------------------------------------------------------------------

			// Теоретически паника может возникнуть до основного перехватчика (например, при обогащении контекста),
			// поэтому обязательно нужен резервный перехватчик, чтобы не завалить весь сервер.

			defer func() {
				if pv := recover(); pv != nil {
					logger.PanicFf(*fallbackPanicCtx, pv, debug.Stack(), "http._boundaryMiddleware (резервный recover)")
					fnResponseOnPanic(w)
				}
			}()

			// ---------------------------------------------------------------------------------------------------------
			// Обогащение
			// ---------------------------------------------------------------------------------------------------------

			ctx := r.Context()
			fallbackPanicCtx = &ctx

			w = kernel.WrapResponseWriter(w)

			requestId := strings.Replace(uuid.Must(uuid.NewRandom()).String(), "-", "", -1)
			kernel.EnrichRequest(r, requestId)
			ctx = logger.AddAttrToCtx(ctx, "kernel.RequestId", requestId)
			/* см. */ _ = kernel.RequestId

			operationId := r.Header.Get("X-Operation-Id")
			operationId = std.Ternary(
				operationId == "",
				strings.Replace(uuid.Must(uuid.NewRandom()).String(), "-", "", -1),
				operationId,
			)
			ctx = processing.EnrichCtx(ctx, operationId)
			ctx = logger.AddAttrToCtx(ctx, "processing.OperationId", operationId)
			/* см. */ _ = processing.OperationId

			// todo : trace... -- OpenTelemetry?

			r = r.WithContext(ctx)
			fallbackPanicCtx = &ctx

			// ---------------------------------------------------------------------------------------------------------
			// Логирование запроса
			// ---------------------------------------------------------------------------------------------------------

			logger.InfoFf(ctx, "Запрос: %s %s", r.Method, r.URL.Path)

			defer func() {
				status := w.(*kernel.ResponseWriter).Status()
				logger.InfoFf(ctx, "Ответ: %d", status)
			}()

			// ---------------------------------------------------------------------------------------------------------
			// Основной перехватчик паники
			// ---------------------------------------------------------------------------------------------------------

			defer func() {
				if pv := recover(); pv != nil {
					fnResponseOnPanic(w)
					logger.PanicFf(ctx, pv, debug.Stack(), "http._boundaryMiddleware")
				}
			}()

			// ---------------------------------------------------------------------------------------------------------

			next.ServeHTTP(w, r)
		})
	}
}

func _securityMiddleware(sec security.Security) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return sec.Middleware(
			// Имеет смысл искать логи по аккаунту,
			// но обогащение контекста для логгера не входит в обязанности security-middleware --
			// нужна еще одна обертка.
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()

				user := sec.AssociatedUser(ctx)
				if user.IsAuthenticated() {
					ctx = logger.AddAttrToCtx(ctx, "account_id", user.AccountId())
				} else {
					// todo : аналогично сделать для гостей с каким-то guest-id ???
				}

				r = r.WithContext(ctx)

				next.ServeHTTP(w, r)
			}),
		)
	}
}
