package layout

import (
	"context"
	"example/admin/gateway/internal/api/http/bundles/layout/handlers"
	"example/admin/gateway/internal/api/http/bundles/layout/openapi"
	"example/admin/gateway/internal/api/http/components/security"
	"github.com/selyukovn/go-std/logger"
	"net/http"
)

// ---------------------------------------------------------------------------------------------------------------------
// Register
// ---------------------------------------------------------------------------------------------------------------------

func Register(
	mux *http.ServeMux,
	sec security.Security,
) {
	server := sDecoratorLoggable{
		sDefault{
			userInfo: handlers.NewUserInfo(sec),
		},
	}

	openapi.HandlerWithOptions(
		openapi.NewStrictHandlerWithOptions(
			server,
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

// ---------------------------------------------------------------------------------------------------------------------
// DEFAULT
// ---------------------------------------------------------------------------------------------------------------------

var _ openapi.StrictServerInterface = sDefault{}

type sDefault struct {
	userInfo func(context.Context, openapi.GetLayoutUserInfoRequestObject) (openapi.GetLayoutUserInfoResponseObject, error)
}

func (s sDefault) GetLayoutUserInfo(ctx context.Context, request openapi.GetLayoutUserInfoRequestObject) (openapi.GetLayoutUserInfoResponseObject, error) {
	return s.userInfo(ctx, request)
}

// ---------------------------------------------------------------------------------------------------------------------
// LOGGABLE
// ---------------------------------------------------------------------------------------------------------------------

var _ openapi.StrictServerInterface = sDecoratorLoggable{}

type sDecoratorLoggable struct {
	origin openapi.StrictServerInterface
}

func (s sDecoratorLoggable) GetLayoutUserInfo(ctx context.Context, request openapi.GetLayoutUserInfoRequestObject) (openapi.GetLayoutUserInfoResponseObject, error) {
	logger.InfoFf(ctx, "%T: %T=%+v", s, request, struct{}{})

	resp, err := s.origin.GetLayoutUserInfo(ctx, request)

	switch vResp := resp.(type) {
	case openapi.GetLayoutUserInfo200JSONResponse:
		logger.InfoFf(ctx, "%T: %T=%+v", s, resp, struct {
			Username  string
			AvatarUrl string
		}{
			Username:  *vResp.Username,
			AvatarUrl: *vResp.AvatarUrl,
		})
	default:
		logger.InfoFf(ctx, "%T: %T=%+v", s, resp, resp)
	}

	return resp, err
}

// ---------------------------------------------------------------------------------------------------------------------
