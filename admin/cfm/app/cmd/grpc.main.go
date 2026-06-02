package main

import (
	"context"
	"example/admin/cfm/cmd/common/container"
	"example/admin/cfm/cmd/common/launcher"
	"example/admin/cfm/cmd/common/monitoring"
	"example/admin/cfm/cmd/common/resources"
	"example/admin/cfm/cmd/grpc"
	"flag"
	"fmt"
	"io"
)

func main() {
	// -----------------------------------------------------------------------------------------------------------------
	// Params
	// -----------------------------------------------------------------------------------------------------------------

	_argDebug := flag.Bool("debug", false, "")
	_argLogFile := flag.String("log-file", "/state/app.log", "путь к log-файлу")
	flag.Parse()
	argDebug := *_argDebug
	argLogFile := *_argLogFile

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

	// mysql
	mysql := resources.OpenMysql(
		env.MysqlHost,
		env.MysqlUser,
		env.MysqlPassword,
		env.MysqlDb,
	)
	defer fnClose("mysql", mysql.Db)

	// -----------------------------------------------------------------------------------------------------------------
	// Container
	// -----------------------------------------------------------------------------------------------------------------

	ctr := container.New(
		logIo,
		argDebug,
		mysql.Db,
		mysql.FnIsDeadlockError,
		mysql.FnIsDuplicateKeyError,
	)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	grpcServer := grpc.NewServer(ctr, env.ApiGrpcApiKey)
	monServer := monitoring.NewMonitoringServer()

	launcher.LaunchServers(ctr.Logger, []launcher.Server{
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
