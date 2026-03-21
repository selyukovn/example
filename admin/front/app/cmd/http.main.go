package main

import (
	"context"
	"example/admin/front/cmd/common/launcher"
	"example/admin/front/cmd/http"
	"example/admin/front/internal/infra/clients/gateway"
	"flag"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"log/slog"
	"os"
)

func main() {
	// -----------------------------------------------------------------------------------------------------------------
	// Params
	// -----------------------------------------------------------------------------------------------------------------

	_argDebug := flag.Bool("debug", false, "")
	flag.Parse()
	argDebug := *_argDebug

	env := http.LoadEnv()

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

	apiClient := gateway.NewApiClient(env.ApiBaseUrl)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	httpServer := http.NewServer(
		apiClient,
		env.AppName,
		env.BaseUrl,
		env.SessionCookieName,
	)

	launcher.LaunchServers([]launcher.Server{
		{
			"HTTP-сервер",
			func(context.Context) error { return httpServer.Start() },
			func(context.Context) error { return httpServer.Stop() },
		},
	})

	// -----------------------------------------------------------------------------------------------------------------
}
