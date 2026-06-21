package main

import (
	"database/sql"
	"example/admin/auth/cmd/grpc"
	"example/admin/auth/cmd/grpc/container"
	"flag"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/selyukovn/example_gopkg/launcher"
	"github.com/selyukovn/example_gopkg/monitoring"
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
	logIo := std.Must(os.OpenFile(argLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666))
	defer fnClose("logIo", logIo)

	// mysql - master
	// todo : рефакторинг
	type tSql = struct {
		Db                    *sql.DB
		FnIsDuplicateKeyError func(error) bool
		FnIsDeadlockError     func(error) bool
	}
	xMysql := func(host string, user string, password string, dbName string) tSql {
		assert.Str().NotEmpty().Must(host)
		assert.Str().NotEmpty().Must(user)
		assert.Str().NotEmpty().Must(password)
		assert.Str().NotEmpty().Must(dbName)

		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", user, password, host, dbName))
		if err != nil {
			panic(err)
		}

		fnIsDuplicateKeyError := func(err error) bool {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok {
				return mysqlErr.Number == 1062
			}
			return false
		}

		fnIsDeadlockError := func(err error) bool {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok {
				return mysqlErr.Number == 1213
			}
			return false
		}

		return tSql{
			Db:                    db,
			FnIsDuplicateKeyError: fnIsDuplicateKeyError,
			FnIsDeadlockError:     fnIsDeadlockError,
		}
	}(
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_HOST_MASTER"), "env: MYSQL_HOST_MASTER"),
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_USER"), "env: MYSQL_USER"),
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_PASSWORD"), "env: MYSQL_PASSWORD"),
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_DB"), "env: MYSQL_DB"),
	)
	defer fnClose("xMysql", xMysql.Db)

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
		xMysql.Db,
		xMysql.FnIsDeadlockError,
		xMysql.FnIsDuplicateKeyError,
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
			grpcServer.Start,
			grpcServer.Stop,
		},
		{
			"Monitoring-сервер",
			monServer.Start,
			monServer.Stop,
		},
	})
}
