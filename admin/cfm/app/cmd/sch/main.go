package main

import (
	"context"
	adapt_domain_cfm "example/admin/cfm/internal/adapt/domain/cfm"
	adapt_domain_cfm_code "example/admin/cfm/internal/adapt/domain/cfm/code"
	adapt_domain_event_storage "example/admin/cfm/internal/adapt/domain/event_storage"
	api_sch "example/admin/cfm/internal/api/sch"
	api_sch_handlers_tick_time "example/admin/cfm/internal/api/sch/handlers/tick_time"
	api_sch_kernel "example/admin/cfm/internal/api/sch/kernel"
	api_sch_middlewares "example/admin/cfm/internal/api/sch/middlewares"
	domain_cfm "example/admin/cfm/internal/domain/cfm"
	domain_event_storage "example/admin/cfm/internal/domain/event_storage"
	opera_domain_facades "example/admin/cfm/internal/opera/domain_facades"
	opera_use_cases_tick_time "example/admin/cfm/internal/opera/use_cases/tick_time"
	"flag"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/selyukovn/example_gopkg/launcher"
	"github.com/selyukovn/example_gopkg/monitoring"
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

	// mysql
	xMysql := openMysql(env.MysqlHost, env.MysqlUser, env.MysqlPassword, env.MysqlDb)
	defer closeResource("xMysql", xMysql.Db)

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
		adapt_domain_event_storage.NewRepositoryImplSql(),
	)

	// cfm
	cfmCodeGen := adapt_domain_cfm_code.NewGeneratorImplUintRand1()
	cfmCodeHasher := adapt_domain_cfm_code.NewHasherImplBcrypt10()
	cfmCodeSender := adapt_domain_cfm_code.NewSenderImplDummy()
	cfmIdGen := adapt_domain_cfm.NewIdGeneratorImplUniqueRandom()
	cfmFactory := domain_cfm.NewFactory(cfmIdGen, cfmCodeGen, cfmCodeHasher)
	cfmRepo := adapt_domain_cfm.NewRepositoryImplSql(xMysql.FnIsDuplicateKeyError)

	// txr
	operaTxr := txr.NewTxrImplSql(xMysql.Db, 2, 50*time.Millisecond, xMysql.FnIsDeadlockError)

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
	cfmDomFac := opera_domain_facades.NewCfmDomFac(
		operaTxr,
		evStorage,
		cfmFactory,
		cfmRepo,
		cfmCodeSender,
		cfmCodeHasher,
	)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	scheduler := api_sch.NewScheduler(func() *cron.Cron {
		xCron := cron.New()
		api_sch_handlers_tick_time.Register(
			xCron,
			10*time.Minute,
			opera_use_cases_tick_time.NewCommand(operaGrt, cfmDomFac),
			[]api_sch_kernel.Middleware{
				api_sch_middlewares.TraceSet(),
				api_sch_middlewares.LogBeginEnd("TickTime"),
				api_sch_middlewares.OnPanic(func(ctx context.Context, pv any, ds []byte) {
					logger.PanicFf(ctx, pv, ds, "api_sch_middlewares"+"."+"OnPanic")
				}),
			},
		)
		return xCron
	}())

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
