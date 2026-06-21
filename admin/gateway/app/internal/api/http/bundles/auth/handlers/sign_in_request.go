package handlers

import (
	"context"
	"example/admin/gateway/internal/api/http/bundles/auth/openapi"
	"example/admin/gateway/internal/api/http/components/security"
	"example/admin/gateway/internal/infra/clients/auth"
	"github.com/selyukovn/go-std"
	"net/http"
	"time"
)

func NewSignInRequest(
	sec security.Security,
	sAuth auth.ClientInterface,
) func(context.Context, openapi.PostAuthSignInRequestRequestObject) (openapi.PostAuthSignInRequestResponseObject, error) {
	return func(ctx context.Context, r openapi.PostAuthSignInRequestRequestObject) (openapi.PostAuthSignInRequestResponseObject, error) {
		user := sec.AssociatedUser(ctx)
		if user.IsAuthenticated() {
			return openapi.PostAuthSignInRequest422JSONResponse{
				Code:    http.StatusForbidden,
				Message: http.StatusText(http.StatusForbidden),
			}, nil
		}

		// --

		email, err := std.EmailFromString(*r.Body.Email)
		if err != nil {
			return openapi.PostAuthSignInRequest422JSONResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		}

		// --

		res, err := sAuth.SignInRequest(ctx, user.Ip(), user.UserAgent(), email)
		switch vErr := err.(type) {
		case nil:
		case auth.ErrorValidation:
			return openapi.PostAuthSignInRequest422JSONResponse{
				Code:    http.StatusBadRequest,
				Message: vErr.Message,
			}, nil
		case std.ErrorNotFound:
			return openapi.PostAuthSignInRequest422JSONResponse{
				Code:    http.StatusNotFound,
				Message: http.StatusText(http.StatusNotFound),
			}, nil
		case auth.ErrorAccountAccessDenied:
			return openapi.PostAuthSignInRequest422JSONResponse{
				Code:    http.StatusForbidden,
				Message: http.StatusText(http.StatusForbidden),
			}, nil
		case std.ErrorRuntime:
			return nil, err
		default:
			panic(err)
		}

		// --

		canRetryAtStr := res.CanRetryAt.Format(time.RFC3339)
		expireAtStr := res.ExpireAt.Format(time.RFC3339)
		return openapi.PostAuthSignInRequest200JSONResponse{
			SignInId:    &res.SignInId,
			RetriesLeft: &res.RetriesLeft,
			CanRetryAt:  &canRetryAtStr,
			ExpireAt:    &expireAtStr,
		}, nil
	}
}
