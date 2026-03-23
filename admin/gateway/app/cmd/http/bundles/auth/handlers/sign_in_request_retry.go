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

func NewSignInRequestRetry(
	ctr *container.Container,
	sec security.Security,
) func(context.Context, openapi.PutAuthSignInRequestRetryRequestObject) (openapi.PutAuthSignInRequestRetryResponseObject, error) {
	return func(ctx context.Context, r openapi.PutAuthSignInRequestRetryRequestObject) (openapi.PutAuthSignInRequestRetryResponseObject, error) {
		user := sec.AssociatedUser(ctx)
		if user.IsAuthenticated() {
			return openapi.PutAuthSignInRequestRetry422JSONResponse{
				Code:    http.StatusForbidden,
				Message: http.StatusText(http.StatusForbidden),
			}, nil
		}

		// --

		signInId := *r.Body.SignInId

		// --

		res, err := ctr.Services.Auth.SignInRequestRetry(ctx, user.Ip(), user.UserAgent(), signInId)
		switch vErr := err.(type) {
		case nil:
		case auth.ErrorValidation:
			return openapi.PutAuthSignInRequestRetry422JSONResponse{
				Code:    http.StatusBadRequest,
				Message: vErr.Message,
			}, nil
		case std.ErrorNotFound:
			return openapi.PutAuthSignInRequestRetry422JSONResponse{
				Code:    http.StatusNotFound,
				Message: http.StatusText(http.StatusNotFound),
			}, nil
		case auth.ErrorAccountAccessDenied:
			return openapi.PutAuthSignInRequestRetry422JSONResponse{
				Code:    http.StatusForbidden,
				Message: http.StatusText(http.StatusForbidden),
			}, nil
		case auth.ErrorSignInFinished:
			logger.WarnFf(ctx, "Обращение к завершенному SignIn %q: %#v", signInId, vErr)
			if vErr.IsPassed {
				return openapi.PutAuthSignInRequestRetry422JSONResponse{
					Code:    http.StatusUnprocessableEntity,
					Message: "Уже подтверждено",
				}, nil
			} else if vErr.IsFailed {
				return openapi.PutAuthSignInRequestRetry422JSONResponse{
					Code:    http.StatusUnprocessableEntity,
					Message: "Уже провалено",
				}, nil
			} else if vErr.IsExpired {
				return openapi.PutAuthSignInRequestRetry422JSONResponse{
					Code:    http.StatusUnprocessableEntity,
					Message: "Время вышло",
				}, nil
			} else {
				panic(vErr)
			}
		case auth.ErrorNoAttemptsLeft:
			return openapi.PutAuthSignInRequestRetry422JSONResponse{
				Code:    http.StatusUnprocessableEntity,
				Message: "Попытки кончились",
			}, nil
		case auth.ErrorRequestsFrequency:
			// фронт обновляет данные, а не реагирует на статус, поэтому отвечаем как при успехе
			canRetryAtStr := vErr.CanReqAfter.Format(time.RFC3339)
			return openapi.PutAuthSignInRequestRetry200JSONResponse{
				SignInId:    &signInId,
				RetriesLeft: &vErr.CanReqAttemptsLeft,
				CanRetryAt:  &canRetryAtStr,
			}, nil
		case std.ErrorUnprocessable:
			// todo : по логике это дубликат IsAsPassed случая cfm.ErrorFinished, но...
			logger.WarnFf(ctx, "Обращение к завершенному SignIn %q с сессией: %#v", signInId, vErr)
			return openapi.PutAuthSignInRequestRetry422JSONResponse{
				Code:    http.StatusUnprocessableEntity,
				Message: "Уже есть сессия",
			}, nil
		case std.ErrorRuntime:
			return nil, err
		default:
			panic(err)
		}

		// --

		canRetryAtStr := std.Ternary(res.RetriesLeft > 0, res.CanRetryAt.Format(time.RFC3339), "")
		return openapi.PutAuthSignInRequestRetry200JSONResponse{
			SignInId:    &res.SignInId,
			RetriesLeft: &res.RetriesLeft,
			CanRetryAt:  &canRetryAtStr,
		}, nil
	}
}
