package main

import (
	"context"
	"database/sql"
	"example/admin/gateway/cmd/kafcon"
	"example/admin/gateway/cmd/kafcon/container"
	"flag"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"github.com/selyukovn/example_gopkg/launcher"
	"github.com/selyukovn/example_gopkg/monitoring"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

const serviceName = "gateway"

func main() {
	// -----------------------------------------------------------------------------------------------------------------
	// Args
	// -----------------------------------------------------------------------------------------------------------------

	_argDebug := flag.Bool("debug", false, "")
	_argLogFile := flag.String("log-file", "/state/server.log", "путь к log-файлу")
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
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_HOST"), "env: MYSQL_HOST"),
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_USER"), "env: MYSQL_USER"),
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_PASSWORD"), "env: MYSQL_PASSWORD"),
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_DB"), "env: MYSQL_DB"),
	)
	defer fnClose("xMysql", xMysql.Db)

	// redis-cache
	// todo : рефакторинг
	redisCacheClient := func(host string, username string, password string, dbNumber uint) *redis.Client {
		assert.Str().NotEmpty().Must(host)
		assert.Str().NotEmpty().Must(username)
		assert.Str().NotEmpty().Must(password)

		opt, err := redis.ParseURL(fmt.Sprintf("redis://%s:%s@%s:6379?db=%d", username, password, host, dbNumber))
		if err != nil {
			panic(err)
		}

		r := redis.NewClient(opt)

		if err := r.Ping(context.Background()).Err(); err != nil {
			panic(err)
		}

		return r
	}(
		assert.Str().NotEmpty().MustGet(os.Getenv("REDIS_CACHE_HOST")),
		assert.Str().NotEmpty().MustGet(os.Getenv("REDIS_CACHE_USER")),
		assert.Str().NotEmpty().MustGet(os.Getenv("REDIS_CACHE_PASSWORD")),
		uint(std.Must[uint64](strconv.ParseUint(os.Getenv("REDIS_CACHE_DB"), 10, 64))),
	)
	defer fnClose("redisCacheClient", redisCacheClient)

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
		redisCacheClient,
		xMysql.Db,
		xMysql.FnIsDeadlockError,
		xMysql.FnIsDuplicateKeyError,
	)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	brokersHostPorts := strings.Split(assert.Str().NotEmpty().MustGet(os.Getenv("KAFKA_BROKERS_HOSTPORTS")), ",")

	adminAuthEventsCnsA := kafcon.NewConsumer(serviceName, kafcon.AdminAuthEventsTopic, "a", brokersHostPorts, ctr)
	adminAuthEventsCnsB := kafcon.NewConsumer(serviceName, kafcon.AdminAuthEventsTopic, "b", brokersHostPorts, ctr)
	adminAuthEventsCnsC := kafcon.NewConsumer(serviceName, kafcon.AdminAuthEventsTopic, "c", brokersHostPorts, ctr)

	monServer := monitoring.NewMonitoringServer()

	launcher.LaunchServers([]launcher.Server{
		{
			Name:    adminAuthEventsCnsA.Id(),
			FnStart: adminAuthEventsCnsA.Start,
			FnStop:  adminAuthEventsCnsA.Stop,
		},
		{
			Name:    adminAuthEventsCnsB.Id(),
			FnStart: adminAuthEventsCnsB.Start,
			FnStop:  adminAuthEventsCnsB.Stop,
		},
		{
			Name:    adminAuthEventsCnsC.Id(),
			FnStart: adminAuthEventsCnsC.Start,
			FnStop:  adminAuthEventsCnsC.Stop,
		},
		{
			Name:    "Monitoring-сервер",
			FnStart: monServer.Start,
			FnStop:  monServer.Stop,
		},
	})
}
