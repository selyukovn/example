package main

import (
	"context"
	"example/admin/auth/cmd/common/launcher"
	"example/admin/auth/cmd/common/monitoring"
	"example/admin/auth/cmd/common/resources"
	"example/admin/auth/cmd/grpc"
	"example/admin/auth/cmd/grpc/container"
	"flag"
	"fmt"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
	"io"
	"log/slog"
	"os"
)

func main() {
	// -----------------------------------------------------------------------------------------------------------------
	// Args
	// -----------------------------------------------------------------------------------------------------------------

	argDebug := *flag.Bool("debug", false, "")
	argLogFile := *flag.String("log-file", "/state/app.log", "путь к log-файлу")

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

	// mysql - master
	mysqlMaster := resources.OpenMysql(
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_HOST_MASTER"), "env: MYSQL_HOST_MASTER"),
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_USER"), "env: MYSQL_USER"),
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_PASSWORD"), "env: MYSQL_PASSWORD"),
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_DB"), "env: MYSQL_DB"),
	)
	defer fnClose("mysqlMaster", mysqlMaster.Db)

	// todo : MYSQL_HOST_REPLICA

	// -----------------------------------------------------------------------------------------------------------------
	// Globals
	// -----------------------------------------------------------------------------------------------------------------

	xLogger := logger.NewSlogLogger(slog.NewJSONHandler(logIo, &slog.HandlerOptions{
		Level: std.Ternary(argDebug, slog.LevelDebug, slog.LevelInfo),
	}))
	logger.SetDefault(xLogger)
	slog.SetDefault(xLogger.SlogLogger())

	// -----------------------------------------------------------------------------------------------------------------
	// Container
	// -----------------------------------------------------------------------------------------------------------------

	ctr := container.New(
		mysqlMaster.Db,
		mysqlMaster.FnIsDeadlockError,
		mysqlMaster.FnIsDuplicateKeyError,
		assert.Str().NotEmpty().MustGet(os.Getenv("SERVICE_CFM_API_GRPC_BASEURL"), "env: SERVICE_CFM_API_GRPC_BASEURL"),
		assert.Str().NotEmpty().MustGet(os.Getenv("SERVICE_CFM_API_GRPC_APIKEY"), "env: SERVICE_CFM_API_GRPC_APIKEY"),
	)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	grpcServer := grpc.NewServer(
		ctr,
		assert.Str().NotEmpty().MustGet(os.Getenv("API_GRPC_APIKEY"), "env: API_GRPC_APIKEY"),
	)
	monServer := monitoring.NewMonitoringServer()

	launcher.LaunchServers([]launcher.Server{
		{
			"GRPC-сервер",
			func(context.Context) error { return grpcServer.Start() },
			func(context.Context) error { return grpcServer.Stop() },
		},
		{
			"Monitoring-сервер",
			func(context.Context) error { return monServer.Start() },
			func(context.Context) error { return monServer.Stop() },
		},
	})

	// -----------------------------------------------------------------------------------------------------------------
}
