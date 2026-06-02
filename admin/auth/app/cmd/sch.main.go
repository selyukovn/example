package main

import (
	"context"
	"example/admin/auth/cmd/common/launcher"
	"example/admin/auth/cmd/common/monitoring"
	"example/admin/auth/cmd/common/resources"
	"example/admin/auth/cmd/sch"
	"example/admin/auth/cmd/sch/container"
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

	env := sch.LoadEnv()

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
	// Container
	// -----------------------------------------------------------------------------------------------------------------

	ctr := container.New(
		logIo,
		argDebug,
		mysqlMaster.Db,
		mysqlMaster.FnIsDeadlockError,
		mysqlMaster.FnIsDuplicateKeyError,
	)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	scheduler := sch.NewScheduler(ctr)
	monServer := monitoring.NewMonitoringServer()

	launcher.LaunchServers(ctr.Logger, []launcher.Server{
		{
			"Scheduler",
			func(context.Context) error { return scheduler.Start() },
			func(context.Context) error { return scheduler.Stop() },
		},
		{
			"Monitoring-сервер",
			func(context.Context) error { return monServer.Start() },
			func(context.Context) error { return monServer.Stop() },
		},
	})

	// -----------------------------------------------------------------------------------------------------------------
}
