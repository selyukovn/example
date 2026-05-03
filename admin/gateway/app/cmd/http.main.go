package main

import (
	"example/admin/gateway/cmd/common/launcher"
	"example/admin/gateway/cmd/common/monitoring"
	"example/admin/gateway/cmd/common/resources"
	"example/admin/gateway/cmd/http"
	"example/admin/gateway/cmd/http/container"
	"flag"
	"fmt"
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

	// redis
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

	// -----------------------------------------------------------------------------------------------------------------
}
