package main

import (
	api_http "example/admin/front/internal/api/http"
	api_http_kernel "example/admin/front/internal/api/http/kernel"
	api_http_layout_general "example/admin/front/internal/api/http/layouts/general"
	api_http_pages_root "example/admin/front/internal/api/http/pages/root"
	api_http_pages_sign_in "example/admin/front/internal/api/http/pages/sign_in"
	api_http_pages_welcome "example/admin/front/internal/api/http/pages/welcome"
	infra_clients_gateway "example/admin/front/internal/infra/clients/gateway"
	"flag"
	"github.com/selyukovn/example_gopkg/launcher"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// -----------------------------------------------------------------------------------------------------------------
	// Args
	// -----------------------------------------------------------------------------------------------------------------

	_argDebug := flag.Bool("debug", false, "")
	flag.Parse()
	argDebug := *_argDebug

	// -----------------------------------------------------------------------------------------------------------------
	// Env
	// -----------------------------------------------------------------------------------------------------------------

	env := loadEnv()

	// -----------------------------------------------------------------------------------------------------------------
	// Globals
	// -----------------------------------------------------------------------------------------------------------------

	xLogger := logger.NewSlogLogger(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: std.Ternary(argDebug, slog.LevelDebug, slog.LevelInfo),
	}))
	logger.SetDefault(xLogger)
	slog.SetDefault(xLogger.SlogLogger())

	// -----------------------------------------------------------------------------------------------------------------
	// Build
	// -----------------------------------------------------------------------------------------------------------------

	apiClient := infra_clients_gateway.NewApiClient(env.ApiBaseUrl)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	api_http_kernel.Configure(env.BaseUrl, env.SessionCookieName)

	httpServer := api_http.NewServer(func() http.Handler {
		mux := http.NewServeMux()
		api_http_pages_root.Register(
			mux,
			api_http_pages_sign_in.Url,
			api_http_pages_welcome.Url,
		)
		api_http_layout_general.Register(
			mux,
			apiClient,
			env.AppName,
			api_http_pages_root.UrlRoot,
			map[string]string{
				api_http_pages_welcome.Title: api_http_pages_welcome.Url,
			},
		)
		api_http_pages_welcome.Register(mux, apiClient, api_http_pages_root.UrlRoot)
		api_http_pages_sign_in.Register(mux, apiClient, env.AppName, api_http_pages_root.UrlRoot)
		// 404
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		})
		// static
		mux.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
		mux.Handle("GET /robots.txt", http.FileServer(http.Dir("./static")))
		return mux
	}())

	launcher.LaunchServers([]launcher.Server{
		{
			"HTTP-сервер",
			httpServer.Start,
			httpServer.Stop,
		},
	})
}
