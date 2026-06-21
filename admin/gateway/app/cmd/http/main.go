package main

import (
	adapt_infra_cache_loggable "example/admin/gateway/internal/adapt/infra/cache/loggable"
	adapt_infra_clients_auth_cachable "example/admin/gateway/internal/adapt/infra/clients/auth/cachable"
	adapt_infra_clients_auth_loggable "example/admin/gateway/internal/adapt/infra/clients/auth/loggable"
	api_http "example/admin/gateway/internal/api/http"
	api_http_bundles_auth "example/admin/gateway/internal/api/http/bundles/auth"
	api_http_bundles_layout "example/admin/gateway/internal/api/http/bundles/layout"
	api_http_bundles_root "example/admin/gateway/internal/api/http/bundles/root"
	api_http_components_security "example/admin/gateway/internal/api/http/components/security"
	api_http_interceptors "example/admin/gateway/internal/api/http/interceptors"
	infra_cache "example/admin/gateway/internal/infra/cache"
	infra_cache_redis "example/admin/gateway/internal/infra/cache/redis"
	infra_clients_auth "example/admin/gateway/internal/infra/clients/auth"
	infra_clients_auth_grpc "example/admin/gateway/internal/infra/clients/auth/grpc"
	"flag"
	"fmt"
	"github.com/selyukovn/example_gopkg/launcher"
	"github.com/selyukovn/example_gopkg/monitoring"
	"github.com/selyukovn/example_gopkg/processing"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"io"
	"log/slog"
	"net/http"
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
	// Env
	// -----------------------------------------------------------------------------------------------------------------

	env := loadEnv()

	// -----------------------------------------------------------------------------------------------------------------
	// Resources
	// -----------------------------------------------------------------------------------------------------------------

	// logIo
	logIo := std.Must(os.OpenFile(argLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666))
	defer closeResource("logIo", logIo)

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

	// ---- auth-service ----
	var sAuth infra_clients_auth.ClientInterface
	sAuth = infra_clients_auth_grpc.NewClientGrpcMust(
		env.ServiceAuthApiGrpcBaseUrl,
		env.ServiceAuthApiGrpcApiKey,
		processing.OperationId,
	)
	// cachable
	var authCache infra_cache.CacheInterface
	authCache = infra_cache_redis.New(redisCacheClient)
	authCache = adapt_infra_cache_loggable.NewDecorator(authCache, true)
	sAuth = adapt_infra_clients_auth_cachable.NewDecorator(sAuth, authCache)
	// loggable
	sAuth = adapt_infra_clients_auth_loggable.NewDecorator(sAuth)
	// ---- /auth-service ----

	sec := api_http_components_security.New(sAuth)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	httpServer := api_http.NewServer(
		api_http_interceptors.Boundary()(
			api_http_interceptors.Metrics()(
				api_http_interceptors.Security(sec)(
					func() http.Handler {
						mux := http.NewServeMux()
						api_http_bundles_auth.Register(mux, sec, sAuth)
						api_http_bundles_layout.Register(mux, sec)
						api_http_bundles_root.Register(mux)
						return mux
					}(),
				),
			),
		),
	)

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

func closeResource(name string, resource io.Closer) {
	if err := resource.Close(); err != nil {
		fmt.Printf("Ошибка закрытия ресурса %s: %s - %#v\n", name, err, err)
	} else {
		fmt.Printf("Ресурс %s закрыт!\n", name)
	}
}
