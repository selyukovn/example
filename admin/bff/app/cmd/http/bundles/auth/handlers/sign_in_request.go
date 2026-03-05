package handlers

import (
	"example/admin/bff/cmd/http/components/security"
	"example/admin/bff/cmd/http/container"
	"example/admin/bff/cmd/http/kernel"
	"example/admin/bff/cmd/http/kernel_ext"
	"example/admin/bff/internal/infra/clients/auth"
	"fmt"
	"github.com/selyukovn/go-std"
	"net/http"
	"time"
)

func NewSignInRequest(ctr *container.Container, sec *security.Security) http.Handler {
	type Response = struct {
		SignInId    string `json:"sign_in_id"`
		RetriesLeft uint   `json:"retries_left"`
		CanRetryAt  string `json:"can_retry_at"`
		ExpireAt    string `json:"expire_at"`
	}

	return sec.AllowOnlyGuests(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		traceId := kernel_ext.TraceId(r)
		fromIp := kernel_ext.UserIp(r)
		fromUag := kernel_ext.UserAgent(r)

		rData := kernel.ParseRequestJson(r, struct {
			Email string `json:"email"`
		}{})
		if rData == nil {
			kernel.Error400(w, "no data")
			return
		}

		email, err := std.EmailFromString(rData.Email)
		if err != nil {
			kernel.Error400(w, "incorrect email")
			return
		}

		// --

		res, err := ctr.Services.Auth.SignInRequest(ctx, traceId, fromIp, fromUag, email)
		switch vErr := err.(type) {
		case nil:
		case auth.ErrorValidation:
			kernel.Error400(w, fmt.Sprintf("%s: %s", vErr.Field, vErr.Message))
			return
		case std.ErrorNotFound:
			kernel.Error404(w, "Аккаунт не найден")
			return
		case auth.ErrorAccountAccessDenied:
			kernel.Error403(w, "Доступ запрещен")
			return
		case std.ErrorRuntime:
			ctr.Logger.CtxErrorFf(ctx, vErr.Error())
			kernel.Error500(w)
			return
		default:
			panic(err)
		}

		// --

		if err := kernel.RenderJson(w, Response{
			SignInId:    res.SignInId,
			RetriesLeft: uint(res.RetriesLeft),
			// это всегда первый запрос -- тут не бывает финальных попыток, поэтому CanReqAfter всегда не 0
			CanRetryAt: res.CanRetryAt.Format(time.RFC3339),
			ExpireAt:   res.ExpireAt.Format(time.RFC3339),
		}); err != nil {
			ctr.Logger.CtxErrorFf(ctx, err.Error())
			kernel.Error500(w)
		}
	}))
}
