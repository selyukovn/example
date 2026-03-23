package handlers

import (
	"context"
	"example/admin/gateway/cmd/http/bundles/auth/openapi"
	"example/admin/gateway/cmd/http/components/security"
	"example/admin/gateway/cmd/http/container"
	"example/admin/gateway/internal/infra/clients/auth"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"net/http"
	"time"
)

func NewSignInConfirm(
	ctr *container.Container,
	sec security.Security,
) func(context.Context, openapi.PutAuthSignInConfirmRequestObject) (openapi.PutAuthSignInConfirmResponseObject, error) {
	return func(ctx context.Context, r openapi.PutAuthSignInConfirmRequestObject) (openapi.PutAuthSignInConfirmResponseObject, error) {
		user := sec.AssociatedUser(ctx)
		if user.IsAuthenticated() {
			return openapi.PutAuthSignInConfirm422JSONResponse{
				Code:    http.StatusForbidden,
				Message: http.StatusText(http.StatusForbidden),
			}, nil
		}

		// --

		signInId := *r.Body.SignInId
		code := *r.Body.Code

		// --

		res, err := ctr.Services.Auth.SignInConfirm(ctx, user.Ip(), user.UserAgent(), signInId, code)
		switch vErr := err.(type) {
		case nil:
		case auth.ErrorValidation:
			return openapi.PutAuthSignInConfirm422JSONResponse{
				Code:    http.StatusBadRequest,
				Message: vErr.Message,
			}, nil
		case std.ErrorNotFound:
			return openapi.PutAuthSignInConfirm422JSONResponse{
				Code:    http.StatusNotFound,
				Message: http.StatusText(http.StatusNotFound),
			}, nil
		case auth.ErrorAccountAccessDenied:
			return openapi.PutAuthSignInConfirm422JSONResponse{
				Code:    http.StatusForbidden,
				Message: http.StatusText(http.StatusForbidden),
			}, nil
		case auth.ErrorSignInFinished:
			logger.WarnFf(ctx, "Обращение к завершенному SignIn %q: %#v", signInId, vErr)
			if vErr.IsPassed {
				return openapi.PutAuthSignInConfirm422JSONResponse{
					Code:    http.StatusUnprocessableEntity,
					Message: "Уже подтверждено",
				}, nil
			} else if vErr.IsFailed {
				return openapi.PutAuthSignInConfirm422JSONResponse{
					Code:    http.StatusUnprocessableEntity,
					Message: "Уже провалено",
				}, nil
			} else if vErr.IsExpired {
				return openapi.PutAuthSignInConfirm422JSONResponse{
					Code:    http.StatusUnprocessableEntity,
					Message: "Время вышло",
				}, nil
			} else {
				panic(vErr)
			}
		case std.ErrorUnprocessable:
			// todo : по логике это дубликат IsAsPassed случая cfm.ErrorFinished, но...
			logger.WarnFf(ctx, "Обращение к завершенному SignIn %q с сессией: %#v", signInId, vErr)
			return openapi.PutAuthSignInConfirm422JSONResponse{
				Code:    http.StatusUnprocessableEntity,
				Message: "Уже есть сессия",
			}, nil
		case std.ErrorRuntime:
			return nil, err
		default:
			panic(err)
		}

		// --

		sessExpAtStr := res.SessionExpireAt.Format(time.RFC3339)
		return openapi.PutAuthSignInConfirm200JSONResponse{
			IsPassed:        &res.IsPassed,
			AttemptsLeft:    &res.AttemptsLeft,
			SessionId:       &res.SessionId,
			SessionExpireAt: &sessExpAtStr,
		}, nil
	}
}
