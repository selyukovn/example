package handlers

import (
	"example/admin/front/cmd/http/kernel"
	"example/admin/front/internal/infra/clients/gateway"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
	"net/http"
	"time"
)

func NewConfirm(
	apiClient gateway.ApiClient,
	redirectUrlForAuthorized string,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if kernel.CookieHasSessId(r) {
			kernel.Error403(w)
			return
		}

		rData := kernel.ParseRequestJson(r, struct {
			SignInId string `json:"sign_in_id"`
			Code     string `json:"code"`
		}{})
		if rData == nil {
			kernel.Error400(w)
			return
		}

		ctx := r.Context()

		fromIp := kernel.ClientIp(r)
		fromUag := kernel.ClientUag(r)

		// --

		resp, err := apiClient.AuthSignInConfirm(fromIp, fromUag, rData.SignInId, rData.Code)
		if err != nil {
			kernel.Error500(w)
			return
		} else if resp.JSON422 != nil && resp.JSON422.Code == 400 {
			kernel.Error400(w, resp.JSON422.Message)
			return
		} else if resp.JSON422 != nil && resp.JSON422.Code == 401 {
			kernel.CookieUnsetSessId(w)
			kernel.Error401(w, resp.JSON422.Message)
			return
		} else if resp.JSON422 != nil && resp.JSON422.Code == 403 {
			kernel.Error403(w, resp.JSON422.Message)
			return
		} else if resp.JSON422 != nil && resp.JSON422.Code == 404 {
			kernel.Error404(w, resp.JSON422.Message)
			return
		} else if resp.JSON422 != nil && resp.JSON422.Code == 422 {
			kernel.Error422(w, resp.JSON422.Message)
			return
		} else {
			assert.NotNilDeepMust(resp.JSON200)
		}

		// --

		if *resp.JSON200.IsPassed {
			sessExpAt, err := time.Parse(time.RFC3339, *resp.JSON200.SessionExpireAt)
			if err != nil {
				logger.ErrorFf(ctx, err.Error())
				kernel.Error500(w)
				return
			}
			kernel.CookieSetSessId(w, *resp.JSON200.SessionId, sessExpAt)
		}

		if err := kernel.RenderJson(w, struct {
			IsPassed     bool   `json:"is_passed"`
			AttemptsLeft uint   `json:"attempts_left"`
			RedirectUrl  string `json:"redirect_url"`
		}{
			IsPassed:     *resp.JSON200.IsPassed,
			AttemptsLeft: uint(*resp.JSON200.AttemptsLeft),
			RedirectUrl:  redirectUrlForAuthorized,
		}); err != nil {
			logger.ErrorFf(ctx, err.Error())
			kernel.Error500(w)
		}
	})
}
