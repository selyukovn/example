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
	"io"
	"log/slog"
)

func main() {
	// -----------------------------------------------------------------------------------------------------------------
	// Params
	// -----------------------------------------------------------------------------------------------------------------

	argDebug := *flag.Bool("debug", false, "")
	argLogFile := *flag.String("log-file", "/state/app.log", "путь к log-файлу")

	env := grpc.LoadEnv()

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
		env.MysqlHostMaster,
		env.MysqlUser,
		env.MysqlPassword,
		env.MysqlDb,
	)
	defer fnClose("mysqlMaster", mysqlMaster.Db)

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
		env.ServiceCfmApiGrpcBaseUrl,
		env.ServiceCfmApiGrpcApiKey,
	)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	grpcServer := grpc.NewServer(ctr, env.ApiGrpcApiKey)
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
