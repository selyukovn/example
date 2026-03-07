package handlers

import (
	"context"
	"example/admin/bff/cmd/http/bundles/auth/config"
	"example/admin/bff/cmd/http/bundles/auth/openapi"
	"example/admin/bff/cmd/http/components/security"
	"example/admin/bff/cmd/http/container"
	"example/admin/bff/internal/infra/clients/auth"
	"github.com/selyukovn/go-std"
	"net/http"
)

func NewSignInConfirm(
	ctr *container.Container,
	sec *security.Security,
	cfg *config.Config,
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

		res, err := ctr.Services.Auth.SignInConfirm(ctx, user.TraceId(), user.Ip(), user.UserAgent(), signInId, code)
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
			ctr.Logger.CtxWarnFf(ctx, "Обращение к завершенному SignIn %q: %#v", signInId, vErr)
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
			ctr.Logger.CtxWarnFf(ctx, "Обращение к завершенному SignIn %q с сессией: %#v", signInId, vErr)
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

		redirectUrl := cfg.UrlRedirectToOnSuccess()
		return signInConfirmSuccessResponse{
			openapi.PutAuthSignInConfirm200JSONResponse{
				IsPassed:     &res.IsPassed,
				AttemptsLeft: &res.AttemptsLeft,
				RedirectUrl:  &redirectUrl,
			},
			func(w http.ResponseWriter) error {
				if res.IsPassed {
					err := sec.Authenticate(ctx, user, w, res.SessionId)
					switch err.(type) {
					case nil:
					case std.ErrorRuntime:
						return std.WrapErrorToRuntime(err, "handlers", "NewSignInConfirm", "Authenticate")
					default:
						panic(err)
					}
				}
				return nil
			},
		}, nil
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type signInConfirmSuccessResponse struct {
	openapi.PutAuthSignInConfirm200JSONResponse
	fnPreResponse func(http.ResponseWriter) error
}

func (r signInConfirmSuccessResponse) VisitPutAuthSignInConfirmResponse(w http.ResponseWriter) error {
	if err := r.fnPreResponse(w); err != nil {
		return err
	}
	return r.PutAuthSignInConfirm200JSONResponse.VisitPutAuthSignInConfirmResponse(w)
}

// ---------------------------------------------------------------------------------------------------------------------
