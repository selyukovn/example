package http

import (
	auth_cfg "example/admin/bff/cmd/http/bundles/auth/config"
	auth_handlers "example/admin/bff/cmd/http/bundles/auth/handlers"
	"example/admin/bff/cmd/http/components/monitoring"
	"example/admin/bff/cmd/http/components/security"
	"example/admin/bff/cmd/http/container"
	"example/admin/bff/cmd/http/kernel"
	"example/admin/bff/cmd/http/kernel_ext"
	assert "github.com/selyukovn/go-wm-assert"
	"net/http"
	"runtime/debug"
)

func registerRoutes(
	mux *http.ServeMux,
	ctr *container.Container,
	appName string,
	baseUrl string,
	sessionCookieName string,
) {
	sec := security.New(ctr, baseUrl, sessionCookieName)

	globalInterceptors := []func(http.Handler) http.Handler{
		_newBoundaryInterceptor(ctr),
		monitoring.NewMetricsInterceptor(),
	}

	// -----------------------------------------------------------------------------------------------------------------
	// Root
	// -----------------------------------------------------------------------------------------------------------------

	// Большого смысла нет заводить маленький bundle для редиректа,
	// поэтому, несмотря на нарушенное однообразие, пусть будет тут.
	_registerDynamicHandlers(mux, globalInterceptors, map[string]http.Handler{
		"GET /{$}": sec.AllowAny(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := sec.AssociatedUser(r)
			if user.IsGuest() {
				kernel.Redirect307(w, r, "/auth/sign-in/welcome/")
			} else {
				// todo : kernel.Redirect307(w, r, ...)
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = w.Write([]byte("Account: " + user.AccountId()))
			}
		})),
	})

	// -----------------------------------------------------------------------------------------------------------------
	// Auth
	// -----------------------------------------------------------------------------------------------------------------

	authStaticPath, authStaticUrl := _registerStaticHandler(mux, "auth")
	authCfg := auth_cfg.New(
		appName,
		"/auth/sign-in/welcome/",
		"/auth/sign-in/request/",
		"/auth/sign-in/request-retry/",
		"/auth/sign-in/confirm/",
		"/",
		authStaticPath,
		authStaticUrl,
	)
	_registerDynamicHandlers(mux, globalInterceptors, map[string]http.Handler{
		"GET /auth/sign-in/welcome/{$}":       auth_handlers.NewSignInWelcome(ctr, authCfg, sec),
		"POST /auth/sign-in/request/{$}":      auth_handlers.NewSignInRequest(ctr, sec),
		"PUT /auth/sign-in/request-retry/{$}": auth_handlers.NewSignInRequestRetry(ctr, sec),
		"PUT /auth/sign-in/confirm/{$}":       auth_handlers.NewSignInConfirm(ctr, authCfg, sec),
		"DELETE /auth/sign-out/{$}":           auth_handlers.NewSignOut(ctr, authCfg, sec),
	})

	// -----------------------------------------------------------------------------------------------------------------

	return
}

func _registerDynamicHandlers(
	mux *http.ServeMux,
	interceptors []func(http.Handler) http.Handler,
	routes map[string]http.Handler,
) {
	for route, handler := range routes {
		for i := len(interceptors) - 1; i >= 0; i-- {
			handler = interceptors[i](handler)
		}
		mux.Handle(route, handler)
	}
}

func _registerStaticHandler(mux *http.ServeMux, bundleName string) (string, string) {
	// см. bff/build/http/Dockerfile
	dirPath := "./static/" + bundleName
	urlPrefix := "/static/" + bundleName + "/" + kernel_ext.CalcFilesVersion(dirPath)

	// `rproxy` кеширует файлы, полученные из `bff` -- cм. README.md.
	mux.Handle(
		// GET разрешает и HEAD
		"GET "+urlPrefix+"/",
		http.StripPrefix(urlPrefix, http.FileServer(http.Dir(dirPath))),
	)

	return dirPath, urlPrefix
}

// Данный перехватчик должен быть самым внешним!
func _newBoundaryInterceptor(ctr *container.Container) func(http.Handler) http.Handler {
	// Все описанные в данном перехватчике действия слишком связаны между собой,
	// чтобы выделить каждое в отдельный перехватчик.
	// Например, логирование статуса ответа не имеет смысла без trace-id из обогащенного контекста,
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
			fnResponseOnPanic := func(w http.ResponseWriter) {
				kernel.Error500(w)
			}

			// ---------------------------------------------------------------------------------------------------------
			// Резервный перехватчик паники
			// ---------------------------------------------------------------------------------------------------------

			// Теоретически паника может возникнуть до основного перехватчика (например, при обогащении контекста),
			// поэтому обязательно нужен резервный перехватчик, чтобы не завалить весь сервер.

			defer func() {
				if pv := recover(); pv != nil {
					ctr.Logger.GeneralPanicFf(pv, debug.Stack(), "http._newBoundaryInterceptor (резервный recover)")
					fnResponseOnPanic(w)
				}
			}()

			// ---------------------------------------------------------------------------------------------------------
			// Обогащение writer'а
			// ---------------------------------------------------------------------------------------------------------

			w = kernel.WrapResponseWriter(w)

			// ---------------------------------------------------------------------------------------------------------
			// Обогащение контекста
			// ---------------------------------------------------------------------------------------------------------

			// Trace Id
			// ----------------

			// Внимание!
			// X-Trace-Id прописывается rproxy -- его отсутствие есть проблема сервера, а не пользователя.
			// Поэтому вместо 400 нужно отдавать 500.
			// см. rproxy/build/server/nginx.conf
			traceId := kernel_ext.TraceId(r)
			if err := assert.Str().NotEmpty().Check(traceId); err != nil {
				ctr.Logger.GeneralErrorFf("Похоже, `rproxy` не передал TraceId!")
				kernel.Error500(w)
				return
			}

			ctx := r.Context()
			ctx = ctr.Logger.AddTraceIdToCtx(ctx, traceId)

			r = r.WithContext(ctx)

			// ---------------------------------------------------------------------------------------------------------
			// Логирование запроса
			// ---------------------------------------------------------------------------------------------------------

			ctr.Logger.CtxInfoFf(ctx, "Запрос: %s %s", r.Method, r.URL.Path)

			defer func() {
				status := w.(*kernel.ResponseWriter).Status()
				ctr.Logger.CtxInfoFf(ctx, "Ответ: %d", status)
			}()

			// ---------------------------------------------------------------------------------------------------------
			// Основной перехватчик паники
			// ---------------------------------------------------------------------------------------------------------

			defer func() {
				if pv := recover(); pv != nil {
					fnResponseOnPanic(w)
					ctr.Logger.CtxPanicFf(ctx, pv, debug.Stack(), "http._newBoundaryInterceptor")
				}
			}()

			// ---------------------------------------------------------------------------------------------------------

			next.ServeHTTP(w, r)
		})
	}
}
