package main

import (
	"context"
	"example/admin/front/cmd/common/launcher"
	"example/admin/front/cmd/http"
	"example/admin/front/internal/infra/clients/gateway"
	infra_logger "example/admin/front/internal/infra/logger"
	"flag"
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
	// Container
	// -----------------------------------------------------------------------------------------------------------------

	logger := infra_logger.NewLogger(os.Stderr, argDebug)
	apiClient := gateway.NewApiClient(env.ApiBaseUrl)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	httpServer := http.NewServer(
		logger,
		apiClient,
		env.AppName,
		env.BaseUrl,
		env.SessionCookieName,
	)

	launcher.LaunchServers(logger, []launcher.Server{
		{
			"HTTP-сервер",
			func(context.Context) error { return httpServer.Start() },
			func(context.Context) error { return httpServer.Stop() },
		},
	})

	// -----------------------------------------------------------------------------------------------------------------
}
