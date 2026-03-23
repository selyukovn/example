package main

import (
	"context"
	"example/admin/cfm/cmd/common/container"
	"example/admin/cfm/cmd/common/launcher"
	"example/admin/cfm/cmd/common/monitoring"
	"example/admin/cfm/cmd/common/resources"
	"example/admin/cfm/cmd/sch"
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

	// mysql
	mysql := resources.OpenMysql(
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_HOST"), "env: MYSQL_HOST"),
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_USER"), "env: MYSQL_USER"),
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_PASSWORD"), "env: MYSQL_PASSWORD"),
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_DB"), "env: MYSQL_DB"),
	)
	defer fnClose("mysql", mysql.Db)

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
		mysql.Db,
		mysql.FnIsDeadlockError,
		mysql.FnIsDuplicateKeyError,
	)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	scheduler := sch.NewScheduler(ctr)
	monServer := monitoring.NewMonitoringServer()

	launcher.LaunchServers([]launcher.Server{
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
