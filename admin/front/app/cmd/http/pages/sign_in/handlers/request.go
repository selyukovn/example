package handlers

import (
	"example/admin/front/cmd/http/kernel"
	"example/admin/front/internal/infra/clients/gateway"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
	"net/http"
)

func NewRequest(
	apiClient *gateway.ApiClient,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if kernel.CookieHasSessId(r) {
			kernel.Error403(w)
			return
		}

		rData := kernel.ParseRequestJson(r, struct {
			Email string `json:"email"`
		}{})
		if rData == nil {
			kernel.Error400(w)
			return
		}

		ctx := r.Context()

		fromIp := kernel.ClientIp(r)
		fromUag := kernel.ClientUag(r)

		// --

		resp, err := apiClient.AuthSignInRequest(fromIp, fromUag, rData.Email)
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

		if err := kernel.RenderJson(w, struct {
			SignInId    string `json:"sign_in_id"`
			RetriesLeft uint   `json:"retries_left"`
			CanRetryAt  string `json:"can_retry_at"`
			ExpireAt    string `json:"expire_at"`
		}{
			SignInId:    *resp.JSON200.SignInId,
			RetriesLeft: uint(*resp.JSON200.RetriesLeft),
			CanRetryAt:  *resp.JSON200.CanRetryAt,
			ExpireAt:    *resp.JSON200.ExpireAt,
		}); err != nil {
			logger.ErrorFf(ctx, err.Error())
			kernel.Error500(w)
		}
	})
}
