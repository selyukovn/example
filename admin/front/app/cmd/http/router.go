package http

import (
	layout_general "example/admin/front/cmd/http/layouts/general"
	"example/admin/front/cmd/http/pages/root"
	"example/admin/front/cmd/http/pages/sign_in"
	"example/admin/front/cmd/http/pages/welcome"
	"example/admin/front/internal/infra/clients/gateway"
	"net/http"
)

func registerRoutes(
	apiClient gateway.ApiClient,
	mux *http.ServeMux,
	appName string,
) {
	// Layouts
	// -----------------------------------------------------------------------------------------------------------------

	layout_general.Register(
		apiClient,
		mux,
		appName,
		root.UrlRoot,
		map[string]string{
			welcome.Title: welcome.Url,
		},
	)

	// Pages
	// -----------------------------------------------------------------------------------------------------------------

	root.Register(
		mux,
		sign_in.Url,
		welcome.Url,
	)

	sign_in.Register(
		apiClient,
		mux,
		appName,
		root.UrlRoot,
	)

	welcome.Register(
		apiClient,
		mux,
		root.UrlRoot,
	)

	// Special
	// -----------------------------------------------------------------------------------------------------------------

	// 404
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// static
	mux.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	mux.Handle("GET /robots.txt", http.FileServer(http.Dir("./static")))
}
