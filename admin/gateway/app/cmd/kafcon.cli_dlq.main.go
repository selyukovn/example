package main

import (
	"context"
	"errors"
	"example/admin/gateway/cmd/common/resources"
	"example/admin/gateway/cmd/kafcon"
	"example/admin/gateway/cmd/kafcon/container"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
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
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_HOST")),
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_USER")),
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_PASSWORD")),
		assert.Str().NotEmpty().MustGet(os.Getenv("MYSQL_DB")),
	)
	defer fnClose("mysql", mysql.Db)

	// redis-cache
	redisCacheClient := resources.OpenRedis(
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
		mysql.Db,
		mysql.FnIsDeadlockError,
		mysql.FnIsDuplicateKeyError,
	)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	ctx := context.Background()
	ctx = logger.AddAttrToCtx(ctx, "dlq_cli_id", strings.Replace(uuid.Must(uuid.NewRandom()).String(), "-", "", -1))
	ctx = logger.AddAttrToCtx(ctx, "dlq_topic", argTopic)
	ctx = logger.AddAttrToCtx(ctx, "dlq_group", argGroupId)

	// Внимание!
	// Функционал может быть существенно расширен и даже представлен в виде отдельного сервиса с UI,
	// однако в рамках данного проекта в этом нет необходимости -- достаточно простого обработчика.

	manager := kafcon.NewDlqProcessor(ctr)

	err := manager.Process(ctx, argTopic, argGroupId)
	if err != nil {
		logger.ErrorFf(ctx, err.Error())
		os.Exit(1)
	}

	// -----------------------------------------------------------------------------------------------------------------
}
