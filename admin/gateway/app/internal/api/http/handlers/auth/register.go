package auth

import (
	"example/admin/gateway/internal/api/http/components/security"
	"example/admin/gateway/internal/api/http/handlers/auth/actions"
	"example/admin/gateway/internal/api/http/handlers/auth/openapi"
	"example/admin/gateway/internal/infra/clients/auth"
	"net/http"
)

func Register(
	mux *http.ServeMux,
	sec security.Security,
	sAuth auth.ClientInterface,
	fnDecorateService func(openapi.StrictServerInterface) openapi.StrictServerInterface,
) {
	var service openapi.StrictServerInterface = newServiceDefault(
		actions.NewSignInRequest(sec, sAuth),
		actions.NewSignInRequestRetry(sec, sAuth),
		actions.NewSignInConfirm(sec, sAuth),
		actions.NewSignOut(sec),
	)

	if fnDecorateService != nil {
		service = fnDecorateService(service)
	}

	openapi.HandlerWithOptions(
		openapi.NewStrictHandlerWithOptions(
			service,
			[]openapi.StrictMiddlewareFunc{},
			openapi.StrictHTTPServerOptions{
				ResponseErrorHandlerFunc: openapi.NewStrictResponseErrorHandler(),
			},
		),
		openapi.StdHTTPServerOptions{
			BaseURL:          "",
			BaseRouter:       mux,
			ErrorHandlerFunc: nil,
		},
	)
}
