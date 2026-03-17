package http

import (
	"example/admin/front/cmd/http/pages/root"
	"example/admin/front/cmd/http/pages/sign_in"
	"example/admin/front/internal/infra/clients/gateway"
	"example/admin/front/internal/infra/logger"
	"net/http"
)

func registerRoutes(
	logger *logger.Logger,
	apiClient *gateway.ApiClient,
	mux *http.ServeMux,
	appName string,
) {
	root.Register(
		mux,
		sign_in.Url,
		"/TODO", // TODO : ...
	)

	sign_in.Register(
		logger,
		apiClient,
		mux,
		appName,
		root.UrlRoot,
	)

	// 404
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// static
	mux.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	mux.Handle("GET /robots.txt", http.FileServer(http.Dir("./static")))
}
