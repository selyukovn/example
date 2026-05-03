package main

import (
	"example/admin/cfm/cmd/common/container"
	"example/admin/cfm/cmd/common/launcher"
	"example/admin/cfm/cmd/common/monitoring"
	"example/admin/cfm/cmd/common/resources"
	"example/admin/cfm/cmd/grpc"
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

	_argDebug := flag.Bool("debug", false, "")
	_argLogFile := flag.String("log-file", "/state/app.log", "путь к log-файлу")
	flag.Parse()
	argDebug := *_argDebug
	argLogFile := *_argLogFile

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

	grpcServer := grpc.NewServer(
		ctr,
		assert.Str().NotEmpty().MustGet(os.Getenv("API_GRPC_APIKEY"), "env: API_GRPC_APIKEY"),
	)
	monServer := monitoring.NewMonitoringServer()

	launcher.LaunchServers([]launcher.Server{
		{
			"GRPC-сервер",
			grpcServer.Start,
			grpcServer.Stop,
		},
		{
			"Monitoring-сервер",
			monServer.Start,
			monServer.Stop,
		},
	})

	// -----------------------------------------------------------------------------------------------------------------
}
