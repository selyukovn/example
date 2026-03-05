package handlers

import (
	"example/admin/bff/cmd/http/bundles/auth/config"
	"example/admin/bff/cmd/http/components/security"
	"example/admin/bff/cmd/http/container"
	"example/admin/bff/cmd/http/kernel"
	"example/admin/bff/cmd/http/kernel_ext"
	"example/admin/bff/internal/infra/clients/auth"
	"fmt"
	"github.com/selyukovn/go-std"
	"net/http"
)

func NewSignInConfirm(ctr *container.Container, cfg *config.Config, sec *security.Security) http.Handler {
	type Response = struct {
		IsPassed     bool   `json:"is_passed"`
		AttemptsLeft uint   `json:"attempts_left"`
		RedirectUrl  string `json:"redirect_url"`
	}

	return sec.AllowOnlyGuests(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		traceId := kernel_ext.TraceId(r)
		fromIp := kernel_ext.UserIp(r)
		fromUag := kernel_ext.UserAgent(r)

		rData := kernel.ParseRequestJson(r, struct {
			SignInId string `json:"sign_in_id"`
			Code     string `json:"code"`
		}{})
		if rData == nil {
			kernel.Error400(w)
			return
		}

		signInId := rData.SignInId
		code := rData.Code

		// --

		res, err := ctr.Services.Auth.SignInConfirm(ctx, traceId, fromIp, fromUag, signInId, code)
		switch vErr := err.(type) {
		case nil:
		case auth.ErrorValidation:
			kernel.Error400(w, fmt.Sprintf("%s: %s", vErr.Field, vErr.Message))
			return
		case std.ErrorNotFound:
			kernel.Error404(w)
			return
		case auth.ErrorAccountAccessDenied:
			kernel.Error403(w)
			return
		case auth.ErrorSignInFinished:
			ctr.Logger.CtxWarnFf(ctx, "Обращение к завершенному SignIn %q: %#v", signInId, vErr)
			if vErr.IsPassed {
				kernel.Error422(w, "Уже подтверждено")
			} else if vErr.IsFailed {
				kernel.Error422(w, "Уже провалено")
			} else if vErr.IsExpired {
				kernel.Error422(w, "Время вышло")
			} else {
				panic(vErr)
			}
			return
		case std.ErrorUnprocessable:
			// todo : по логике это дубликат IsAsPassed случая cfm.ErrorFinished, но...
			ctr.Logger.CtxWarnFf(ctx, "Обращение к завершенному SignIn %q с сессией: %#v", signInId, vErr)
			kernel.Error422(w, "Уже есть сессия")
			return
		case std.ErrorRuntime:
			ctr.Logger.CtxErrorFf(ctx, vErr.Error())
			kernel.Error500(w)
			return
		default:
			panic(err)
		}

		// --

		if res.IsPassed {
			sec.AuthorizeClient(w, res.SessionId, res.SessionExpireAt)
		}

		if err := kernel.RenderJson(w, Response{
			IsPassed:     res.IsPassed,
			AttemptsLeft: uint(res.AttemptsLeft),
			RedirectUrl:  cfg.UrlSignInWelcome(),
		}); err != nil {
			ctr.Logger.CtxErrorFf(ctx, err.Error())
			kernel.Error500(w)
		}
	}))
}
