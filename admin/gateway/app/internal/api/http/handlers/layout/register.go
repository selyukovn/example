package layout

import (
	"example/admin/gateway/internal/api/http/components/security"
	"example/admin/gateway/internal/api/http/handlers/layout/actions"
	"example/admin/gateway/internal/api/http/handlers/layout/openapi"
	"net/http"
)

func Register(
	mux *http.ServeMux,
	sec security.Security,
	fnDecorateService func(openapi.StrictServerInterface) openapi.StrictServerInterface,
) {
	var service openapi.StrictServerInterface = NewService(
		actions.NewUserInfo(sec),
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
