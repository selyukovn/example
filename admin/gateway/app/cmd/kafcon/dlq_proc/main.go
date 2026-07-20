package main

import (
	"context"
	"errors"
	adapt_api_components_dlq "example/admin/gateway/internal/adapt/api/kafcon/components/dlq"
	adapt_api_kafcon_handlers_admin_auth_events_loggable "example/admin/gateway/internal/adapt/api/kafcon/handlers/admin_auth_events/loggable"
	adapt_api_kafcon_handlers_admin_auth_events_trace_get "example/admin/gateway/internal/adapt/api/kafcon/handlers/admin_auth_events/trace_get"
	adapt_infra_cache_loggable "example/admin/gateway/internal/adapt/infra/cache/loggable"
	adapt_infra_clients_auth_cachable "example/admin/gateway/internal/adapt/infra/clients/auth/cachable"
	api_kafcon "example/admin/gateway/internal/api/kafcon"
	api_kafcon_handlers_admin_auth_events "example/admin/gateway/internal/api/kafcon/handlers/admin_auth_events"
	api_kafcon_handlers_admin_auth_events_kafapi "example/admin/gateway/internal/api/kafcon/handlers/admin_auth_events/kafapi"
	api_kafcon_kernel "example/admin/gateway/internal/api/kafcon/kernel"
	infra_cache "example/admin/gateway/internal/infra/cache"
	infra_cache_redis "example/admin/gateway/internal/infra/cache/redis"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"github.com/selyukovn/go-txr"
	assert "github.com/selyukovn/go-wm-assert"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

func main() {
	// -----------------------------------------------------------------------------------------------------------------
	// Args
	// -----------------------------------------------------------------------------------------------------------------

	_argDebug := flag.Bool("debug", false, "")
	_argLogFile := flag.String("log-file", "/state/cli_dlq.log", "путь к log-файлу")
	_argTopic := flag.String("topic", "", "топик")
	_argGroupId := flag.String("groupId", "", "идентификатор группы сообщений")
	flag.Parse()
	argDebug := *_argDebug
	argLogFile := *_argLogFile
	argTopic := *_argTopic
	argGroupId := *_argGroupId

	if err := errors.Join(
		assert.Str().NotEmpty("topic").Check(argTopic),
		assert.Str().NotEmpty("groupId").Check(argGroupId),
	); err != nil {
		println(err.Error())
		os.Exit(1)
	}

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

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	router := api_kafcon_kernel.NewRouter()
	api_kafcon_handlers_admin_auth_events.Register(
		router,
		sAuthCacher,
		[]func(api_kafcon_handlers_admin_auth_events_kafapi.ServiceInterface) api_kafcon_handlers_admin_auth_events_kafapi.ServiceInterface{
			func(service api_kafcon_handlers_admin_auth_events_kafapi.ServiceInterface) api_kafcon_handlers_admin_auth_events_kafapi.ServiceInterface {
				return adapt_api_kafcon_handlers_admin_auth_events_trace_get.NewDecorator(
					adapt_api_kafcon_handlers_admin_auth_events_loggable.NewDecorator(
						service,
					),
				)
			},
		},
		nil,
	)

	// Внимание!
	// Функционал может быть существенно расширен и даже представлен в виде отдельного сервиса с UI,
	// однако в рамках данного проекта в этом нет необходимости -- достаточно простого обработчика.

	dlqProcessor := api_kafcon.NewDlqProcessor(dlqStorage, router)

	ctx := context.Background()
	ctx = logger.AddAttrToCtx(ctx, "dlq_cli_id", strings.Replace(uuid.Must(uuid.NewRandom()).String(), "-", "", -1))
	ctx = logger.AddAttrToCtx(ctx, "dlq_topic", argTopic)
	ctx = logger.AddAttrToCtx(ctx, "dlq_group", argGroupId)

	err := dlqProcessor.Process(ctx, argTopic, argGroupId)
	if err != nil {
		logger.ErrorFf(ctx, err.Error())
		os.Exit(1)
	}
}

func closeResource(name string, resource io.Closer) {
	if err := resource.Close(); err != nil {
		fmt.Printf("Ошибка закрытия ресурса %s: %s - %#v\n", name, err, err)
	} else {
		fmt.Printf("Ресурс %s закрыт!\n", name)
	}
}
