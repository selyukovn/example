package main

import (
	"context"
	"example/admin/gateway/cmd/common/launcher"
	"example/admin/gateway/cmd/common/monitoring"
	"example/admin/gateway/cmd/common/resources"
	"example/admin/gateway/cmd/http"
	"example/admin/gateway/cmd/http/container"
	"flag"
	"fmt"
	"io"
)

func main() {
	// -----------------------------------------------------------------------------------------------------------------
	// Params
	// -----------------------------------------------------------------------------------------------------------------

	argDebug := *flag.Bool("debug", false, "")
	argLogFile := *flag.String("log-file", "/state/app.log", "путь к log-файлу")

	env := http.LoadEnv()

	// -----------------------------------------------------------------------------------------------------------------
	// Resources
	// -----------------------------------------------------------------------------------------------------------------

	fnClose := func(name string, resource io.Closer) {
		if err := resource.Close(); err != nil {
			fmt.Printf("Ошибка закрытия ресурса %s: %s - %#v\n", name, err, err)
		} else {
			fmt.Printf("Ресурс %s закрыт!\n", name)
		}
	}

	// logIo
	logIo := resources.NewLogIoFile(argLogFile)
	defer fnClose("logIo", logIo)

	// -----------------------------------------------------------------------------------------------------------------
	// Container
	// -----------------------------------------------------------------------------------------------------------------

	ctr := container.New(
		logIo,
		argDebug,
		env.ServiceAuthApiGrpcBaseUrl,
		env.ServiceAuthApiGrpcApiKey,
	)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	httpServer := http.NewServer(ctr)
	monServer := monitoring.NewMonitoringServer()

	launcher.LaunchServers(ctr.Logger, []launcher.Server{
		{
			"HTTP-сервер",
			func(context.Context) error { return httpServer.Start() },
			func(context.Context) error { return httpServer.Stop() },
		},
		{
			"Monitoring-сервер",
			func(context.Context) error { return monServer.Start() },
			func(context.Context) error { return monServer.Stop() },
		},
	})

	// -----------------------------------------------------------------------------------------------------------------
}
