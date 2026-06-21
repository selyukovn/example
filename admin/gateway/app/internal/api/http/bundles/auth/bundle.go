package auth

import (
	"context"
	"example/admin/gateway/internal/api/http/bundles/auth/handlers"
	"example/admin/gateway/internal/api/http/bundles/auth/openapi"
	"example/admin/gateway/internal/api/http/components/security"
	"example/admin/gateway/internal/infra/clients/auth"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"net/http"
)

// ---------------------------------------------------------------------------------------------------------------------
// Register
// ---------------------------------------------------------------------------------------------------------------------

func Register(
	mux *http.ServeMux,
	sec security.Security,
	sAuth auth.ClientInterface,
) {
	server := sDecoratorLoggable{
		sDefault{
			signInRequest:      handlers.NewSignInRequest(sec, sAuth),
			signInRequestRetry: handlers.NewSignInRequestRetry(sec, sAuth),
			signInConfirm:      handlers.NewSignInConfirm(sec, sAuth),
			signOut:            handlers.NewSignOut(sec),
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
	signInRequest      func(context.Context, openapi.PostAuthSignInRequestRequestObject) (openapi.PostAuthSignInRequestResponseObject, error)
	signInRequestRetry func(context.Context, openapi.PutAuthSignInRequestRetryRequestObject) (openapi.PutAuthSignInRequestRetryResponseObject, error)
	signInConfirm      func(context.Context, openapi.PutAuthSignInConfirmRequestObject) (openapi.PutAuthSignInConfirmResponseObject, error)
	signOut            func(context.Context, openapi.DeleteAuthSignOutRequestObject) (openapi.DeleteAuthSignOutResponseObject, error)
}

func (r sDefault) PostAuthSignInRequest(ctx context.Context, request openapi.PostAuthSignInRequestRequestObject) (openapi.PostAuthSignInRequestResponseObject, error) {
	return r.signInRequest(ctx, request)
}

func (r sDefault) PutAuthSignInRequestRetry(ctx context.Context, request openapi.PutAuthSignInRequestRetryRequestObject) (openapi.PutAuthSignInRequestRetryResponseObject, error) {
	return r.signInRequestRetry(ctx, request)
}

func (r sDefault) PutAuthSignInConfirm(ctx context.Context, request openapi.PutAuthSignInConfirmRequestObject) (openapi.PutAuthSignInConfirmResponseObject, error) {
	return r.signInConfirm(ctx, request)
}

func (r sDefault) DeleteAuthSignOut(ctx context.Context, request openapi.DeleteAuthSignOutRequestObject) (openapi.DeleteAuthSignOutResponseObject, error) {
	return r.signOut(ctx, request)
}

// ---------------------------------------------------------------------------------------------------------------------
// LOGGABLE
// ---------------------------------------------------------------------------------------------------------------------

var _ openapi.StrictServerInterface = sDecoratorLoggable{}

type sDecoratorLoggable struct {
	origin openapi.StrictServerInterface
}

func (s sDecoratorLoggable) PostAuthSignInRequest(ctx context.Context, request openapi.PostAuthSignInRequestRequestObject) (openapi.PostAuthSignInRequestResponseObject, error) {
	logger.InfoFf(ctx, "%T: %+v", request, struct {
		Email string
	}{
		Email: *request.Body.Email,
	})

	resp, err := s.origin.PostAuthSignInRequest(ctx, request)

	switch vResp := resp.(type) {
	case openapi.PostAuthSignInRequest200JSONResponse:
		logger.InfoFf(ctx, "%T: %+v", resp, struct {
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
		logger.InfoFf(ctx, "%T: %+v", resp, resp)
	}

	return resp, err
}

func (s sDecoratorLoggable) PutAuthSignInRequestRetry(ctx context.Context, request openapi.PutAuthSignInRequestRetryRequestObject) (openapi.PutAuthSignInRequestRetryResponseObject, error) {
	logger.InfoFf(ctx, "%T: %+v", request, struct {
		SignInId string
	}{
		SignInId: *request.Body.SignInId,
	})

	resp, err := s.origin.PutAuthSignInRequestRetry(ctx, request)

	switch vResp := resp.(type) {
	case openapi.PutAuthSignInRequestRetry200JSONResponse:
		logger.InfoFf(ctx, "%T: %+v", resp, struct {
			CanRetryAt  string
			RetriesLeft int
			SignInId    string
		}{
			CanRetryAt:  *vResp.CanRetryAt,
			RetriesLeft: *vResp.RetriesLeft,
			SignInId:    *vResp.SignInId,
		})
	default:
		logger.InfoFf(ctx, "%T: %+v", resp, resp)
	}

	return resp, err
}

func (s sDecoratorLoggable) PutAuthSignInConfirm(ctx context.Context, request openapi.PutAuthSignInConfirmRequestObject) (openapi.PutAuthSignInConfirmResponseObject, error) {
	logger.InfoFf(ctx, "%T: %+v", request, struct {
		SignInId string
		Code     string
	}{
		SignInId: *request.Body.SignInId,
		Code:     std.MaskStrNotFirstLast(*request.Body.Code),
	})

	resp, err := s.origin.PutAuthSignInConfirm(ctx, request)

	switch vResp := resp.(type) {
	case openapi.PutAuthSignInConfirm200JSONResponse:
		logger.InfoFf(ctx, "%T: %+v", resp, struct {
			AttemptsLeft    int
			IsPassed        bool
			SessionId       string
			SessionExpireAt string
		}{
			AttemptsLeft:    *vResp.AttemptsLeft,
			IsPassed:        *vResp.IsPassed,
			SessionId:       std.MaskStrNotFirstLast(*vResp.SessionId),
			SessionExpireAt: *vResp.SessionExpireAt,
		})
	default:
		logger.InfoFf(ctx, "%T: %+v", resp, resp)
	}

	return resp, err
}

func (s sDecoratorLoggable) DeleteAuthSignOut(ctx context.Context, request openapi.DeleteAuthSignOutRequestObject) (openapi.DeleteAuthSignOutResponseObject, error) {
	logger.InfoFf(ctx, "%T: %+v", request, request)

	resp, err := s.origin.DeleteAuthSignOut(ctx, request)

	switch vResp := resp.(type) {
	case openapi.DeleteAuthSignOut200JSONResponse:
		logger.InfoFf(ctx, "%T: %+v", resp, struct {
			Success bool
		}{
			Success: *vResp.Success,
		})
	default:
		logger.InfoFf(ctx, "%T: %+v", resp, resp)
	}

	return resp, err
}

// ---------------------------------------------------------------------------------------------------------------------
