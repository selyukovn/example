package main

import (
	"example/admin/front/internal/api/http"
	"example/admin/front/internal/infra/clients/gateway"
	"flag"
	"github.com/selyukovn/example_gopkg/launcher"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
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

	apiClient := gateway.NewApiClient(env.ApiBaseUrl)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	httpServer := http.NewServer(apiClient, env.AppName, env.BaseUrl, env.SessionCookieName)

	launcher.LaunchServers([]launcher.Server{
		{
			"HTTP-сервер",
			httpServer.Start,
			httpServer.Stop,
		},
	})
}
