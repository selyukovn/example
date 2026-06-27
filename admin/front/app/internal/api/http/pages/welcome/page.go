package welcome

import (
	"example/admin/front/internal/api/http/kernel"
	"example/admin/front/internal/api/http/layouts/general"
	"example/admin/front/internal/infra/clients/gateway"
	"github.com/selyukovn/go-std/logger"
	"math/rand"
	"net/http"
)

// ---------------------------------------------------------------------------------------------------------------------

const Title = "Главная"
const Url = "/welcome/"

// ---------------------------------------------------------------------------------------------------------------------

func Register(
	mux *http.ServeMux,
	apiClient gateway.ApiClient,
	redirectUrlForGuests string,
) {
	mux.Handle("GET "+Url+"{$}", newRenderer(
		apiClient,
		redirectUrlForGuests,
	))
}

// ---------------------------------------------------------------------------------------------------------------------

func newRenderer(
	apiClient gateway.ApiClient,
	redirectUrlForGuests string,
) http.Handler {
	view := general.MakeView(apiClient, "static/pages/welcome/page.html")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessId := kernel.CookieGetSessId(r)
		if sessId == "" {
			kernel.Redirect307(w, r, redirectUrlForGuests)
			return
		}

		ctx := r.Context()

		// --

		quotes := []string{
			"Шаг влево, шаг вправо — два шага.",
			"Одна ошибка — и ты ошибся.",
			"Работа — не волк. Никто не волк. Только волк — волк.",
			"В жизни всегда есть две дороги: одна — первая, а другая — вторая.",
			"Если заблудился в лесу, иди домой.",
			"Делай, как надо. Как не надо, не делай.",
		}

		// --

		if err := view.Render(w, r, struct {
			Title string
			Quote string
		}{
			Title: Title,
			Quote: quotes[rand.Intn(len(quotes))],
		}); err != nil {
			logger.ErrorFf(ctx, err.Error())
			kernel.Error500(w)
		}
	})
}

// ---------------------------------------------------------------------------------------------------------------------
