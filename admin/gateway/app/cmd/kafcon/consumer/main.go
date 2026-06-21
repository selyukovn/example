package main

import (
	adapt_api_components_dlq "example/admin/gateway/internal/adapt/api/kafcon/components/dlq"
	adapt_infra_cache_loggable "example/admin/gateway/internal/adapt/infra/cache/loggable"
	adapt_infra_clients_auth_cachable "example/admin/gateway/internal/adapt/infra/clients/auth/cachable"
	api_kafcon "example/admin/gateway/internal/api/kafcon"
	api_kafcon_bundles "example/admin/gateway/internal/api/kafcon/bundles"
	api_kafcon_bundles_admin_auth_events "example/admin/gateway/internal/api/kafcon/bundles/admin_auth_events"
	infra_cache "example/admin/gateway/internal/infra/cache"
	infra_cache_redis "example/admin/gateway/internal/infra/cache/redis"
	"flag"
	"fmt"
	"github.com/selyukovn/example_gopkg/launcher"
	"github.com/selyukovn/example_gopkg/monitoring"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"github.com/selyukovn/go-txr"
	"io"
	"log/slog"
	"os"
	"time"
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
	// Env
	// -----------------------------------------------------------------------------------------------------------------

	env := loadEnv()

	// -----------------------------------------------------------------------------------------------------------------
	// Resources
	// -----------------------------------------------------------------------------------------------------------------

	// logIo
	logIo := std.Must(os.OpenFile(argLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666))
	defer closeResource("logIo", logIo)

	// mysql
	xMysql := openMysql(env.MysqlHost, env.MysqlUser, env.MysqlPassword, env.MysqlDb)
	defer closeResource("xMysql", xMysql.Db)

	// redis-cache
	redisCacheClient := openRedis(env.RedisCacheHost, env.RedisCacheUser, env.RedisCachePassword, env.RedisCacheDb)
	defer closeResource("redisCacheClient", redisCacheClient)

	// -----------------------------------------------------------------------------------------------------------------
	// Globals
	// -----------------------------------------------------------------------------------------------------------------

	xLogger := logger.NewSlogLogger(slog.NewJSONHandler(logIo, &slog.HandlerOptions{
		Level: std.Ternary(argDebug, slog.LevelDebug, slog.LevelInfo),
	}))
	logger.SetDefault(xLogger)
	slog.SetDefault(xLogger.SlogLogger())

	// -----------------------------------------------------------------------------------------------------------------
	// Build
	// -----------------------------------------------------------------------------------------------------------------

	var cache infra_cache.CacheInterface
	cache = infra_cache_redis.New(redisCacheClient)
	cache = adapt_infra_cache_loggable.NewDecorator(cache, true)

	sqlTxr := txr.NewTxrImplSql(xMysql.Db, 2, 50*time.Millisecond, xMysql.FnIsDeadlockError)

	sAuthCacher := adapt_infra_clients_auth_cachable.NewCacher(cache)

	dlqStorage := adapt_api_components_dlq.NewStorageSQL(sqlTxr)

	router := api_kafcon_bundles.NewRouter()
	router.Register(
		api_kafcon_bundles_admin_auth_events.TopicName,
		api_kafcon_bundles_admin_auth_events.NewHandlerDecoratorDlq(
			api_kafcon_bundles_admin_auth_events.NewHandlerDefault(sAuthCacher),
			dlqStorage,
		),
	)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	adminAuthEventsCnsA := api_kafcon.NewConsumer(
		serviceName,
		api_kafcon_bundles_admin_auth_events.TopicName,
		"a",
		env.KafkaBrokersHostPorts,
		router,
	)
	adminAuthEventsCnsB := api_kafcon.NewConsumer(
		serviceName,
		api_kafcon_bundles_admin_auth_events.TopicName,
		"b",
		env.KafkaBrokersHostPorts,
		router,
	)
	adminAuthEventsCnsC := api_kafcon.NewConsumer(
		serviceName,
		api_kafcon_bundles_admin_auth_events.TopicName,
		"c",
		env.KafkaBrokersHostPorts,
		router,
	)

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

func closeResource(name string, resource io.Closer) {
	if err := resource.Close(); err != nil {
		fmt.Printf("Ошибка закрытия ресурса %s: %s - %#v\n", name, err, err)
	} else {
		fmt.Printf("Ресурс %s закрыт!\n", name)
	}
}
