package main

import (
	"example/admin/front/cmd/common/launcher"
	"example/admin/front/cmd/http"
	"example/admin/front/internal/infra/clients/gateway"
	"flag"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
	"log/slog"
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
	// Globals
	// -----------------------------------------------------------------------------------------------------------------

	xLogger := logger.NewSlogLogger(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: std.Ternary(argDebug, slog.LevelDebug, slog.LevelInfo),
	}))
	logger.SetDefault(xLogger)
	slog.SetDefault(xLogger.SlogLogger())

	// -----------------------------------------------------------------------------------------------------------------
	// Container
	// -----------------------------------------------------------------------------------------------------------------

	apiClient := gateway.NewApiClient(
		assert.Str().NotEmpty().MustGet(os.Getenv("API_BASEURL"), "env: API_BASEURL"),
	)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	httpServer := http.NewServer(
		apiClient,
		assert.Str().NotEmpty().MustGet(os.Getenv("APP_NAME"), "env: APP_NAME"),
		assert.Str().NotEmpty().MustGet(os.Getenv("BASE_URL"), "env: BASE_URL"),
		assert.Str().NotEmpty().MustGet(os.Getenv("SESSION_COOKIE_NAME"), "env: SESSION_COOKIE_NAME"),
	)

	launcher.LaunchServers([]launcher.Server{
		{
			"HTTP-сервер",
			httpServer.Start,
			httpServer.Stop,
		},
	})

	// -----------------------------------------------------------------------------------------------------------------
}
