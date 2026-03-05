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

func NewSignOut(ctr *container.Container, cfg *config.Config, sec *security.Security) http.Handler {
	return sec.AllowOnlyAuthorized(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		traceId := kernel_ext.TraceId(r)
		fromIp := kernel_ext.UserIp(r)
		fromUag := kernel_ext.UserAgent(r)

		rData := kernel.ParseRequestJson(r, struct {
			SessionId string `json:"session_id"`
		}{})
		if rData == nil {
			kernel.Error400(w)
			return
		}

		sessionId := rData.SessionId

		// --

		err := ctr.Services.Auth.SignOut(ctx, traceId, fromIp, fromUag, sessionId)
		switch vErr := err.(type) {
		case nil:
		case auth.ErrorValidation:
			kernel.Error400(w, fmt.Sprintf("%s: %s", vErr.Field, vErr.Message))
			return
		case std.ErrorNotFound:
			kernel.Error404(w, "Сессия не найдена")
			return
		case auth.ErrorAccountAccessDenied:
			kernel.Error403(w)
			return
		case std.ErrorAlreadyDone:
			kernel.Error422(w, "Сессия уже закрыта")
			return
		case std.ErrorRuntime:
			ctr.Logger.CtxErrorFf(ctx, vErr.Error())
			kernel.Error500(w)
			return
		default:
			panic(err)
		}

		// --

		sec.UnAuthorizeClient(w)

		if err := kernel.RenderJson(w, struct {
			RedirectUrl string `json:"redirect_url"`
		}{
			RedirectUrl: cfg.UrlSignInWelcome(),
		}); err != nil {
			ctr.Logger.CtxErrorFf(ctx, err.Error())
			kernel.Error500(w)
		}
	}))
}
