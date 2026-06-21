package main

import (
	"context"
	adapt_domain_event_storage "example/admin/auth/internal/adapt/domain/event_storage"
	adapt_domain_session "example/admin/auth/internal/adapt/domain/session"
	api_sch "example/admin/auth/internal/api/sch"
	api_sch_handlers "example/admin/auth/internal/api/sch/handlers"
	api_sch_interceptors "example/admin/auth/internal/api/sch/interceptors"
	domain_event_storage "example/admin/auth/internal/domain/event_storage"
	domain_session "example/admin/auth/internal/domain/session"
	opera_domain_facades "example/admin/auth/internal/opera/domain_facades"
	opera_use_cases_session_tick_time "example/admin/auth/internal/opera/use_cases/session_tick_time"
	"flag"
	"fmt"
	"github.com/selyukovn/example_gopkg/launcher"
	"github.com/selyukovn/example_gopkg/monitoring"
	"github.com/selyukovn/example_gopkg/processing"
	goroutiner "github.com/selyukovn/go-routiner"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"github.com/selyukovn/go-txr"
	"io"
	"log/slog"
	"os"
	"time"
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

	// mysql - master
	mysqlMaster := openMysql(env.MysqlHostMaster, env.MysqlUser, env.MysqlPassword, env.MysqlDb)
	defer closeResource("mysqlMaster", mysqlMaster.Db)

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

	// event storage
	evStorage := domain_event_storage.NewStorage(
		adapt_domain_event_storage.NewRepositoryImplSql(processing.OperationId),
	)

	// session
	sessIdGen := adapt_domain_session.NewIdGeneratorImplUniqueRandom()
	sessFactory := domain_session.NewFactory(sessIdGen)
	sessRepo := adapt_domain_session.NewRepositoryImplSql(mysqlMaster.FnIsDuplicateKeyError)

	// txr
	operaTxr := txr.NewTxrImplSql(mysqlMaster.Db, 2, 50*time.Millisecond, mysqlMaster.FnIsDeadlockError)

	// goroutiner
	operaGrt := goroutiner.New(
		goroutiner.MwPanicToError(func(panicValue any, debugStack []byte, ctx context.Context) error {
			_o_ := "main.operaGrt"
			_m_ := "MwPanicToError"
			logger.PanicFf(ctx, panicValue, debugStack, _o_+"."+_m_)
			var err error
			switch pv := panicValue.(type) {
			case error:
				err = fmt.Errorf("panic: %w; stack: %s", pv, string(debugStack))
			case string, fmt.Stringer:
				err = fmt.Errorf("panic: %q; stack: %s", pv, string(debugStack))
			default:
				err = fmt.Errorf("panic: %#v; stack: %s", pv, string(debugStack))
			}
			return std.WrapErrorToRuntime(err, _o_, _m_)
		}),
	)

	// domain facades
	sessDomFac := opera_domain_facades.NewSessionDomFac(operaTxr, evStorage, sessFactory, sessRepo)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	router := api_sch.NewRouter()
	router.Register(
		"sessionTickTime",
		10*time.Minute,
		api_sch_handlers.NewSessionTickTime(
			opera_use_cases_session_tick_time.NewCommand(operaGrt, sessDomFac),
		),
	)
	scheduler := api_sch.NewScheduler(
		router,
		api_sch_interceptors.NewBoundary(),
	)

	monServer := monitoring.NewMonitoringServer()

	launcher.LaunchServers([]launcher.Server{
		{
			"Scheduler",
			scheduler.Start,
			scheduler.Stop,
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
