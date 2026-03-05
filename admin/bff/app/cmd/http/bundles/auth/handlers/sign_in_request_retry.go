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

func NewSignInRequestRetry(ctr *container.Container, sec *security.Security) http.Handler {
	type Response = struct {
		SignInId    string `json:"sign_in_id"`
		RetriesLeft uint   `json:"retries_left"`
		CanRetryAt  string `json:"can_retry_at"`
	}

	return sec.AllowOnlyGuests(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		traceId := kernel_ext.TraceId(r)
		fromIp := kernel_ext.UserIp(r)
		fromUag := kernel_ext.UserAgent(r)

		rData := kernel.ParseRequestJson(r, struct {
			SignInId string `json:"sign_in_id"`
		}{})
		if rData == nil {
			kernel.Error400(w)
			return
		}

		signInId := rData.SignInId

		// --

		res, err := ctr.Services.Auth.SignInRequestRetry(ctx, traceId, fromIp, fromUag, signInId)
		switch vErr := err.(type) {
		case nil:
		case auth.ErrorValidation:
			kernel.Error400(w, fmt.Sprintf("%s: %s", vErr.Field, vErr.Message))
			return
		case std.ErrorNotFound:
			kernel.Error404(w)
			return
		case auth.ErrorAccountAccessDenied:
			kernel.Error403(w, "Доступ запрещен")
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
		case auth.ErrorNoAttemptsLeft:
			kernel.Error422(w, "Попытки кончились")
			return
		case auth.ErrorRequestsFrequency:
			// фронт обновляет данные, а не реагирует на статус, поэтому отвечаем как при успехе
			if err := kernel.RenderJson(w, Response{
				SignInId:    signInId,
				RetriesLeft: uint(vErr.CanReqAttemptsLeft),
				CanRetryAt:  vErr.CanReqAfter.Format(time.RFC3339),
			}); err != nil {
				ctr.Logger.CtxErrorFf(ctx, err.Error())
				kernel.Error500(w)
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

		// can again
		if res.RetriesLeft > 0 {
			if err := kernel.RenderJson(w, Response{
				SignInId:    signInId,
				RetriesLeft: uint(res.RetriesLeft),
				CanRetryAt:  res.CanRetryAt.Format(time.RFC3339),
			}); err != nil {
				ctr.Logger.CtxErrorFf(ctx, err.Error())
				kernel.Error500(w)
			}
			return // !!!
		}

		// last
		if err := kernel.RenderJson(w, Response{
			SignInId:    signInId,
			RetriesLeft: 0,
			CanRetryAt:  "",
		}); err != nil {
			ctr.Logger.CtxErrorFf(ctx, err.Error())
			kernel.Error500(w)
		}
	}))
}
