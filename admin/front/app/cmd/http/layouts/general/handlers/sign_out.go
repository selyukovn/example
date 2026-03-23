package handlers

import (
	"example/admin/front/cmd/http/kernel"
	"example/admin/front/internal/infra/clients/gateway"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
	"net/http"
)

func NewSignOut(
	apiClient gateway.ApiClient,
	redirectUrlForGuests string,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessId := kernel.CookieGetSessId(r)
		if sessId == "" {
			kernel.Error401(w)
			return
		}

		ctx := r.Context()

		fromIp := kernel.ClientIp(r)
		fromUag := kernel.ClientUag(r)

		// --

		resp, err := apiClient.AuthSignOut(fromIp, fromUag, sessId)
		if err != nil {
			kernel.Error500(w)
			return
		} else if resp.JSON422 != nil && resp.JSON422.Code == 400 {
			kernel.Error400(w, resp.JSON422.Message)
			return
		} else if resp.JSON422 != nil && resp.JSON422.Code == 401 {
			// для текущей реализации фронта это равносильно success'у
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

		kernel.CookieUnsetSessId(w)

		if err := kernel.RenderJson(w, struct {
			RedirectUrl string `json:"redirect_url"`
		}{
			RedirectUrl: redirectUrlForGuests,
		}); err != nil {
			logger.ErrorFf(ctx, err.Error())
			kernel.Error500(w)
		}
	})
}
