package main

import (
	"context"
	adapt_domain_account "example/admin/auth/internal/adapt/domain/account"
	adapt_domain_action_request "example/admin/auth/internal/adapt/domain/action_request"
	adapt_domain_cfm "example/admin/auth/internal/adapt/domain/cfm"
	adapt_domain_event_storage "example/admin/auth/internal/adapt/domain/event_storage"
	adapt_domain_session "example/admin/auth/internal/adapt/domain/session"
	adapt_infra_clients_cfm_loggable "example/admin/auth/internal/adapt/infra/clients/cfm/loggable"
	api_grpc "example/admin/auth/internal/api/grpc"
	api_grpc_handlers_auth "example/admin/auth/internal/api/grpc/handlers/auth"
	api_grpc_kernel "example/admin/auth/internal/api/grpc/kernel"
	api_grpc_middlewares "example/admin/auth/internal/api/grpc/middlewares"
	domain_account "example/admin/auth/internal/domain/account"
	domain_action_request "example/admin/auth/internal/domain/action_request"
	domain_cfm "example/admin/auth/internal/domain/cfm"
	domain_event_storage "example/admin/auth/internal/domain/event_storage"
	domain_session "example/admin/auth/internal/domain/session"
	infra_clients_cfm_grpc "example/admin/auth/internal/infra/clients/cfm/grpc"
	opera_domain_facades "example/admin/auth/internal/opera/domain_facades"
	opera_use_cases_check_session "example/admin/auth/internal/opera/use_cases/check_session"
	opera_use_cases_sign_in_confirm "example/admin/auth/internal/opera/use_cases/sign_in_confirm"
	opera_use_cases_sign_in_request "example/admin/auth/internal/opera/use_cases/sign_in_request"
	opera_use_cases_sign_in_request_retry "example/admin/auth/internal/opera/use_cases/sign_in_request_retry"
	opera_use_cases_sign_out "example/admin/auth/internal/opera/use_cases/sign_out"
	"flag"
	"fmt"
	"github.com/selyukovn/example_gopkg/launcher"
	"github.com/selyukovn/example_gopkg/monitoring"
	"github.com/selyukovn/example_gopkg/processing"
	goroutiner "github.com/selyukovn/go-routiner"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"github.com/selyukovn/go-txr"
	"google.golang.org/grpc"
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

	// todo : MYSQL_HOST_REPLICA

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

	// clients - cfm
	infraCfmGrpcClient := adapt_infra_clients_cfm_loggable.NewDecorator(
		infra_clients_cfm_grpc.NewClientGrpcMust(
			env.ServiceCfmApiGrpcBaseUrl,
			env.ServiceCfmApiGrpcApiKey,
			processing.OperationId,
		),
	)

	// event storage
	evStorage := domain_event_storage.NewStorage(
		adapt_domain_event_storage.NewRepositoryImplSql(processing.OperationId),
	)

	// account
	accIdGen := adapt_domain_account.NewIdGeneratorImplUniqueRandom()
	accFactory := domain_account.NewFactory(accIdGen)
	accRepo := adapt_domain_account.NewRepositoryImplSql(mysqlMaster.FnIsDuplicateKeyError)

	// action request
	actReqIdGen := adapt_domain_action_request.NewIdGeneratorImplUniqueRandom()
	actReqFactory := domain_action_request.NewFactory(actReqIdGen)
	actReqRepo := adapt_domain_action_request.NewRepositoryImplSql(mysqlMaster.FnIsDuplicateKeyError)

	// cfm
	var cfmService domain_cfm.ServiceInterface = adapt_domain_cfm.NewServiceImplCfmService(infraCfmGrpcClient)

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
	accDomFac := opera_domain_facades.NewAccountDomFac(operaTxr, evStorage, accFactory, accRepo)
	actReqDomFac := opera_domain_facades.NewActionRequestDomFac(operaTxr, evStorage, actReqFactory, actReqRepo)
	sessDomFac := opera_domain_facades.NewSessionDomFac(operaTxr, evStorage, sessFactory, sessRepo)

	// -----------------------------------------------------------------------------------------------------------------
	// Launch
	// -----------------------------------------------------------------------------------------------------------------

	grpcServer := api_grpc.NewServer(func() *grpc.Server {
		s := grpc.NewServer(grpc.ChainUnaryInterceptor(
			api_grpc_kernel.RootMiddleware(),
			api_grpc_middlewares.PanicToError(func(ctx context.Context, pv any, ds []byte) error {
				logger.PanicFf(ctx, pv, ds, "api_grpc_middlewares"+"."+"PanicToError"+"(outer)")
				return api_grpc_kernel.ErrorInternal()
			}),
			api_grpc_middlewares.Metrics(),
			api_grpc_middlewares.TraceGet(),
			api_grpc_middlewares.LogBeginEnd(),
			api_grpc_middlewares.PanicToError(func(ctx context.Context, pv any, ds []byte) error {
				logger.PanicFf(ctx, pv, ds, "api_grpc_middlewares"+"."+"PanicToError"+"(inner)")
				return api_grpc_kernel.ErrorInternal()
			}),
			api_grpc_middlewares.AccessKey(env.ApiGrpcApiKey),
		))
		api_grpc_handlers_auth.Register(
			s,
			opera_use_cases_sign_in_request.NewCommand(accDomFac, actReqDomFac, cfmService),
			opera_use_cases_sign_in_request_retry.NewCommand(operaGrt, accDomFac, actReqDomFac, cfmService, sessDomFac),
			opera_use_cases_sign_in_confirm.NewCommand(operaGrt, accDomFac, actReqDomFac, cfmService, sessDomFac),
			opera_use_cases_sign_out.NewCommand(accDomFac, sessDomFac),
			opera_use_cases_check_session.NewCommand(accDomFac, sessDomFac),
			nil,
		)
		return s
	}())

	monServer := monitoring.NewMonitoringServer()

	launcher.LaunchServers([]launcher.Server{
		{
			"GRPC-сервер",
			grpcServer.Start,
			grpcServer.Stop,
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
