package auth

import (
	"context"
	"example/admin/bff/cmd/http/bundles/auth/config"
	"example/admin/bff/cmd/http/bundles/auth/handlers"
	"example/admin/bff/cmd/http/bundles/auth/openapi"
	"example/admin/bff/cmd/http/components/security"
	"example/admin/bff/cmd/http/components/static"
	"example/admin/bff/cmd/http/container"
	"github.com/selyukovn/go-std"
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

	var xRouter openapi.StrictServerInterface = &router{
		signInWelcome:      handlers.NewSignInWelcome(sec, cfg),
		signInRequest:      handlers.NewSignInRequest(ctr, sec),
		signInRequestRetry: handlers.NewSignInRequestRetry(ctr, sec),
		signInConfirm:      handlers.NewSignInConfirm(ctr, sec, cfg),
		signOut:            handlers.NewSignOut(sec),
	}
	xRouter = &routerLogRequestResponseData{
		StrictServerInterface: xRouter,
		ctr:                   ctr,
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

// DEFAULT
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

// LOGGABLE
// ---------------------------------------------------------------------------------------------------------------------

var _ openapi.StrictServerInterface = (*routerLogRequestResponseData)(nil)

type routerLogRequestResponseData struct {
	openapi.StrictServerInterface
	ctr *container.Container
}

func (r *routerLogRequestResponseData) PostAuthSignInRequest(ctx context.Context, request openapi.PostAuthSignInRequestRequestObject) (openapi.PostAuthSignInRequestResponseObject, error) {
	r.ctr.Logger.CtxInfoFf(ctx, "%T: %+v", request, struct {
		Email string
	}{
		Email: *request.Body.Email,
	})

	resp, err := r.StrictServerInterface.PostAuthSignInRequest(ctx, request)

	switch vResp := resp.(type) {
	case openapi.PostAuthSignInRequest200JSONResponse:
		r.ctr.Logger.CtxInfoFf(ctx, "%T: %+v", resp, struct {
			CanRetryAt  string
			ExpireAt    string
			RetriesLeft int
			SignInId    string
		}{
			CanRetryAt:  *vResp.CanRetryAt,
			ExpireAt:    *vResp.ExpireAt,
			RetriesLeft: *vResp.RetriesLeft,
			SignInId:    *vResp.SignInId,
		})
	default:
		r.ctr.Logger.CtxInfoFf(ctx, "%T: %+v", resp, resp)
	}

	return resp, err
}

func (r *routerLogRequestResponseData) PutAuthSignInRequestRetry(ctx context.Context, request openapi.PutAuthSignInRequestRetryRequestObject) (openapi.PutAuthSignInRequestRetryResponseObject, error) {
	r.ctr.Logger.CtxInfoFf(ctx, "%T: %+v", request, struct {
		SignInId string
	}{
		SignInId: *request.Body.SignInId,
	})

	resp, err := r.StrictServerInterface.PutAuthSignInRequestRetry(ctx, request)

	switch vResp := resp.(type) {
	case openapi.PutAuthSignInRequestRetry200JSONResponse:
		r.ctr.Logger.CtxInfoFf(ctx, "%T: %+v", resp, struct {
			CanRetryAt  string
			RetriesLeft int
			SignInId    string
		}{
			CanRetryAt:  *vResp.CanRetryAt,
			RetriesLeft: *vResp.RetriesLeft,
			SignInId:    *vResp.SignInId,
		})
	default:
		r.ctr.Logger.CtxInfoFf(ctx, "%T: %+v", resp, resp)
	}

	return resp, err
}

func (r *routerLogRequestResponseData) PutAuthSignInConfirm(ctx context.Context, request openapi.PutAuthSignInConfirmRequestObject) (openapi.PutAuthSignInConfirmResponseObject, error) {
	r.ctr.Logger.CtxInfoFf(ctx, "%T: %+v", request, struct {
		SignInId string
		Code     string
	}{
		SignInId: *request.Body.SignInId,
		Code:     std.MaskStrNotFirstLast(*request.Body.Code),
	})

	resp, err := r.StrictServerInterface.PutAuthSignInConfirm(ctx, request)

	switch vResp := resp.(type) {
	case openapi.PutAuthSignInConfirm200JSONResponse:
		r.ctr.Logger.CtxInfoFf(ctx, "%T: %+v", resp, struct {
			AttemptsLeft int
			IsPassed     bool
			RedirectUrl  string
		}{
			AttemptsLeft: *vResp.AttemptsLeft,
			IsPassed:     *vResp.IsPassed,
			RedirectUrl:  *vResp.RedirectUrl,
		})
	default:
		r.ctr.Logger.CtxInfoFf(ctx, "%T: %+v", resp, resp)
	}

	return resp, err
}

// #####################################################################################################################
