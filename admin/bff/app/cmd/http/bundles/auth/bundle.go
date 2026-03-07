package auth

import (
	"context"
	"example/admin/bff/cmd/http/bundles/auth/config"
	"example/admin/bff/cmd/http/bundles/auth/handlers"
	"example/admin/bff/cmd/http/bundles/auth/openapi"
	"example/admin/bff/cmd/http/components/security"
	"example/admin/bff/cmd/http/components/static"
	"example/admin/bff/cmd/http/container"
	"net/http"
)

// #####################################################################################################################
// CONSTANTS
// #####################################################################################################################

const (
	UrlSignInWelcome = openapi.UrlSignInWelcome
	UrlSignOut       = openapi.UrlSignOut
)

// #####################################################################################################################
// ROUTER
// #####################################################################################################################

func Register(
	mux *http.ServeMux,
	middlewares []func(http.Handler) http.Handler,
	ctr *container.Container,
	sec *security.Security,
	appName string,
	urlRedirectOnSuccess string,
) {
	staticPath, staticUrl := static.RegisterFileHandler(mux, "auth")

	cfg := config.New(
		appName,
		urlRedirectOnSuccess,
		staticPath,
		staticUrl,
	)

	xRouter := &router{
		signInWelcome:      handlers.NewSignInWelcome(sec, cfg),
		signInRequest:      handlers.NewSignInRequest(ctr, sec),
		signInRequestRetry: handlers.NewSignInRequestRetry(ctr, sec),
		signInConfirm:      handlers.NewSignInConfirm(ctr, sec, cfg),
		signOut:            handlers.NewSignOut(sec),
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
				ResponseErrorHandlerFunc: openapi.NewStrictResponseErrorHandler(ctr),
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

var _ openapi.StrictServerInterface = (*router)(nil)

type router struct {
	signInWelcome      func(context.Context, openapi.GetAuthSignInWelcomeRequestObject) (openapi.GetAuthSignInWelcomeResponseObject, error)
	signInRequest      func(context.Context, openapi.PostAuthSignInRequestRequestObject) (openapi.PostAuthSignInRequestResponseObject, error)
	signInRequestRetry func(context.Context, openapi.PutAuthSignInRequestRetryRequestObject) (openapi.PutAuthSignInRequestRetryResponseObject, error)
	signInConfirm      func(context.Context, openapi.PutAuthSignInConfirmRequestObject) (openapi.PutAuthSignInConfirmResponseObject, error)
	signOut            func(context.Context, openapi.DeleteAuthSignOutRequestObject) (openapi.DeleteAuthSignOutResponseObject, error)
}

func (r *router) GetAuthSignInWelcome(ctx context.Context, request openapi.GetAuthSignInWelcomeRequestObject) (openapi.GetAuthSignInWelcomeResponseObject, error) {
	return r.signInWelcome(ctx, request)
}

func (r *router) PostAuthSignInRequest(ctx context.Context, request openapi.PostAuthSignInRequestRequestObject) (openapi.PostAuthSignInRequestResponseObject, error) {
	return r.signInRequest(ctx, request)
}

func (r *router) PutAuthSignInRequestRetry(ctx context.Context, request openapi.PutAuthSignInRequestRetryRequestObject) (openapi.PutAuthSignInRequestRetryResponseObject, error) {
	return r.signInRequestRetry(ctx, request)
}

func (r *router) PutAuthSignInConfirm(ctx context.Context, request openapi.PutAuthSignInConfirmRequestObject) (openapi.PutAuthSignInConfirmResponseObject, error) {
	return r.signInConfirm(ctx, request)
}

func (r *router) DeleteAuthSignOut(ctx context.Context, request openapi.DeleteAuthSignOutRequestObject) (openapi.DeleteAuthSignOutResponseObject, error) {
	return r.signOut(ctx, request)
}

// #####################################################################################################################
