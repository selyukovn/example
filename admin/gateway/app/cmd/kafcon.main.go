package main

import (
	"context"
	"example/admin/gateway/cmd/common/launcher"
	"example/admin/gateway/cmd/common/monitoring"
	"example/admin/gateway/cmd/common/resources"
	"example/admin/gateway/cmd/kafcon"
	"example/admin/gateway/cmd/kafcon/container"
	"flag"
	"fmt"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

const ownerServiceName = "gateway"

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

	ctr := container.New(redisCacheClient)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	var launcherServers []launcher.Server

	// Внимание!
	// В идеальной среде в каждом сервисе-приемнике использовалось бы несколько консьюмер-контейнеров на каждый топик.
	// В данном проекте для экономии ресурсов каждый сервис-приемник использует один консьюмер-контейнер на все топики,
	// однако, разделение по топикам для наглядности выполнено с помощью отдельных горутин
	// (`kafka.Consumer` может обрабатывать все топики и одной горутиной-консьюмером, но этот подход менее нагляден).

	brokerHostPorts := strings.Split(assert.Str().NotEmpty().MustGet(os.Getenv("KAFKA_BROKERS_HOSTPORTS")), ",")

	fnAddConsumers := func(topicName string, consumersNumber int) {
		for i := 1; i <= consumersNumber; i++ {
			name := fmt.Sprintf("consumer-%s-%d", topicName, i)
			server := kafcon.NewServer(ownerServiceName, brokerHostPorts, topicName, strconv.Itoa(i), ctr)
			launcherServers = append(launcherServers, launcher.Server{
				Name:    name,
				FnStart: func(server *kafcon.Server) func(ctx context.Context) error { return server.Start }(server),
				FnStop:  func(server *kafcon.Server) func(ctx context.Context) error { return server.Stop }(server),
			})
		}
	}

	fnAddConsumers(kafcon.AdminAuthEventsTopic, 3)

	// --

	monServer := monitoring.NewMonitoringServer()
	launcherServers = append(launcherServers, launcher.Server{
		Name:    "Monitoring-сервер",
		FnStart: monServer.Start,
		FnStop:  monServer.Stop,
	})

	launcher.LaunchServers(launcherServers)

	// -----------------------------------------------------------------------------------------------------------------
}
