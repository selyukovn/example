package layout

import (
	"context"
	"example/admin/gateway/cmd/http/bundles/layout/handlers"
	"example/admin/gateway/cmd/http/bundles/layout/openapi"
	"example/admin/gateway/cmd/http/components/security"
	"example/admin/gateway/cmd/http/container"
	"github.com/selyukovn/go-std/logger"
	"net/http"
)

// ---------------------------------------------------------------------------------------------------------------------
// Register
// ---------------------------------------------------------------------------------------------------------------------

func Register(
	mux *http.ServeMux,
	middlewares []func(http.Handler) http.Handler,
	ctr *container.Container,
	sec security.Security,
) {
	var xRouter openapi.StrictServerInterface = router{
		userInfo: handlers.NewUserInfo(ctr, sec),
	}
	xRouter = routerLogRequestResponseData{
		StrictServerInterface: xRouter,
	}

	// Генерируемый `github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1` код
	// использует middleware в порядке "последний -- внешний", поэтому набор нужно инвертировать.
	mwLen := len(middlewares)
	openApiMiddlewares := make([]openapi.MiddlewareFunc, mwLen)
	for i, m := range middlewares {
		openApiMiddlewares[mwLen-i-1] = m
	}

	openapi.HandlerWithOptions(
		openapi.NewStrictHandlerWithOptions(
			xRouter,
			[]openapi.StrictMiddlewareFunc{},
			openapi.StrictHTTPServerOptions{
				ResponseErrorHandlerFunc: openapi.NewStrictResponseErrorHandler(),
			},
		),
		openapi.StdHTTPServerOptions{
			BaseURL:          "",
			BaseRouter:       mux,
			Middlewares:      openApiMiddlewares,
			ErrorHandlerFunc: nil,
		},
	)
}

// ---------------------------------------------------------------------------------------------------------------------
// DEFAULT
// ---------------------------------------------------------------------------------------------------------------------

var _ openapi.StrictServerInterface = router{}

type router struct {
	userInfo func(context.Context, openapi.GetLayoutUserInfoRequestObject) (openapi.GetLayoutUserInfoResponseObject, error)
}

func (r router) GetLayoutUserInfo(ctx context.Context, request openapi.GetLayoutUserInfoRequestObject) (openapi.GetLayoutUserInfoResponseObject, error) {
	return r.userInfo(ctx, request)
}

// ---------------------------------------------------------------------------------------------------------------------
// LOGGABLE
// ---------------------------------------------------------------------------------------------------------------------

var _ openapi.StrictServerInterface = routerLogRequestResponseData{}

type routerLogRequestResponseData struct {
	openapi.StrictServerInterface
}

func (r routerLogRequestResponseData) GetLayoutUserInfo(ctx context.Context, request openapi.GetLayoutUserInfoRequestObject) (openapi.GetLayoutUserInfoResponseObject, error) {
	logger.InfoFf(ctx, "%T: %T=%+v", r, request, struct{}{})

	resp, err := r.StrictServerInterface.GetLayoutUserInfo(ctx, request)

	switch vResp := resp.(type) {
	case openapi.GetLayoutUserInfo200JSONResponse:
		logger.InfoFf(ctx, "%T: %T=%+v", r, resp, struct {
			Username  string
			AvatarUrl string
		}{
			Username:  *vResp.Username,
			AvatarUrl: *vResp.AvatarUrl,
		})
	default:
		logger.InfoFf(ctx, "%T: %T=%+v", r, resp, resp)
	}

	return resp, err
}

// ---------------------------------------------------------------------------------------------------------------------
