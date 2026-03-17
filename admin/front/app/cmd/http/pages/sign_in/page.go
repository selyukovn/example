package sign_in

import (
	"example/admin/front/cmd/http/kernel"
	"example/admin/front/cmd/http/pages/sign_in/handlers"
	"example/admin/front/internal/infra/clients/gateway"
	"example/admin/front/internal/infra/logger"
	"html/template"
	"net/http"
)

// ---------------------------------------------------------------------------------------------------------------------

const Url = "/sign-in/"

// ---------------------------------------------------------------------------------------------------------------------

func Register(
	logger *logger.Logger,
	apiClient *gateway.ApiClient,
	mux *http.ServeMux,
	appName string,
	redirectUrlForAuthorized string,
) {
	const HandlerUrlRequest = "/sign-in/request/"
	const HandlerUrlRequestRetry = "/sign-in/request-retry/"
	const HandlerUrlConfirm = "/sign-in/confirm/"

	mux.Handle("GET "+Url+"{$}", newRenderer(
		logger,
		appName,
		redirectUrlForAuthorized,
		HandlerUrlRequest,
		HandlerUrlRequestRetry,
		HandlerUrlConfirm,
	))
	mux.Handle("POST "+HandlerUrlRequest+"{$}", handlers.NewRequest(logger, apiClient))
	mux.Handle("PUT "+HandlerUrlRequestRetry+"{$}", handlers.NewRequestRetry(logger, apiClient))
	mux.Handle("PUT "+HandlerUrlConfirm+"{$}", handlers.NewConfirm(logger, apiClient, redirectUrlForAuthorized))
}

// ---------------------------------------------------------------------------------------------------------------------

func newRenderer(
	logger *logger.Logger,
	appName string,
	redirectUrlForAuthorized string,
	handlerUrlRequest string,
	handlerUrlRequestRetry string,
	handlerUrlConfirm string,
) http.Handler {
	tpl := template.Must(template.ParseFiles("static/pages/sign_in/page.html"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if kernel.CookieHasSessId(r) {
			kernel.Redirect307(w, r, redirectUrlForAuthorized)
			return
		}

		if err := tpl.Execute(w, struct {
			AppName         string
			UrlRequest      string
			UrlRequestRetry string
			UrlConfirm      string
		}{
			AppName:         appName,
			UrlRequest:      handlerUrlRequest,
			UrlRequestRetry: handlerUrlRequestRetry,
			UrlConfirm:      handlerUrlConfirm,
		}); err != nil {
			logger.GeneralErrorFf(err.Error())
			kernel.Error500(w)
		}
	})
}

// ---------------------------------------------------------------------------------------------------------------------
