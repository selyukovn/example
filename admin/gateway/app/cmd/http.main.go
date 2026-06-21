package main

import (
	"context"
	"example/admin/gateway/cmd/http"
	"example/admin/gateway/cmd/http/container"
	"flag"
	"fmt"
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
		assert.Str().NotEmpty().MustGet(os.Getenv("SERVICE_AUTH_API_GRPC_BASEURL"), "env: SERVICE_AUTH_API_GRPC_BASEURL"),
		assert.Str().NotEmpty().MustGet(os.Getenv("SERVICE_AUTH_API_GRPC_APIKEY"), "env: SERVICE_AUTH_API_GRPC_APIKEY"),
	)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	httpServer := http.NewServer(ctr)
	monServer := monitoring.NewMonitoringServer()

	launcher.LaunchServers([]launcher.Server{
		{
			"HTTP-сервер",
			httpServer.Start,
			httpServer.Stop,
		},
		{
			"Monitoring-сервер",
			monServer.Start,
			monServer.Stop,
		},
	})
}
