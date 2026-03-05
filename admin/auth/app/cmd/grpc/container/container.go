package container

import (
	"context"
	"database/sql"
	adapt_domain_account "example/admin/auth/internal/adapt/domain/account"
	adapt_domain_action_request "example/admin/auth/internal/adapt/domain/action_request"
	adapt_domain_cfm "example/admin/auth/internal/adapt/domain/cfm"
	adapt_domain_event_storage "example/admin/auth/internal/adapt/domain/event_storage"
	adapt_domain_session "example/admin/auth/internal/adapt/domain/session"
	adapt_infra_clients_cfm "example/admin/auth/internal/adapt/infra/clients/cfm"
	adapt_opera_components "example/admin/auth/internal/adapt/opera/components"
	domain_account "example/admin/auth/internal/domain/account"
	domain_action_request "example/admin/auth/internal/domain/action_request"
	domain_cfm "example/admin/auth/internal/domain/cfm"
	domain_event_storage "example/admin/auth/internal/domain/event_storage"
	domain_session "example/admin/auth/internal/domain/session"
	infra_clients_cfm_grpc "example/admin/auth/internal/infra/clients/cfm/grpc"
	infra_logger "example/admin/auth/internal/infra/logger"
	opera_domain_facades "example/admin/auth/internal/opera/domain_facades"
	"example/admin/auth/internal/opera/use_cases/check_session"
	"example/admin/auth/internal/opera/use_cases/sign_in_confirm"
	"example/admin/auth/internal/opera/use_cases/sign_in_request"
	"example/admin/auth/internal/opera/use_cases/sign_in_request_retry"
	"example/admin/auth/internal/opera/use_cases/sign_out"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	goroutiner "github.com/selyukovn/go-routiner"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-txr"
	assert "github.com/selyukovn/go-wm-assert"
	"io"
	"time"
)

type Container struct {
	Logger   *infra_logger.Logger
	UseCases UseCases
}

type UseCases = struct {
	SignInRequest      *sign_in_request.Command
	SignInRequestRetry *sign_in_request_retry.Command
	SignInConfirm      *sign_in_confirm.Command
	SignOut            *sign_out.Command
	CheckSession       *check_session.Command
}

func New(
	logIo io.Writer,
	isDebug bool,
	sqlDb *sql.DB,
	sqlDbFnIsDeadlockError func(error) bool,
	sqlDbFnIsDuplicateKeyError func(error) bool,
	appCfmApiGrpcBaseUrl string,
	appCfmApiGrpcApiKey string,
) *Container {
	assert.NotNilDeepMust(logIo)
	assert.NotNilDeepMust(logIo)
	assert.NotNilDeepMust(sqlDb)
	assert.NotNilDeepMust(sqlDbFnIsDeadlockError)
	assert.NotNilDeepMust(sqlDbFnIsDuplicateKeyError)
	assert.NotNilDeepMust(appCfmApiGrpcBaseUrl)
	assert.NotNilDeepMust(appCfmApiGrpcApiKey)

	// -----------------------------------------------------------------------------------------------------------------
	// Infra
	// -----------------------------------------------------------------------------------------------------------------

	// logger
	infraLogger := infra_logger.NewLogger(logIo, isDebug)

	// clients - cfm
	infraCfmGrpcClient := adapt_infra_clients_cfm.NewDecoratorLoggable(
		infra_clients_cfm_grpc.NewClientGrpcMust(appCfmApiGrpcBaseUrl, appCfmApiGrpcApiKey),
		infraLogger,
	)

	// -----------------------------------------------------------------------------------------------------------------
	// Domain
	// -----------------------------------------------------------------------------------------------------------------

	// event storage
	evStorage := domain_event_storage.NewStorage(
		adapt_domain_event_storage.NewRepositoryImplSql(),
	)

	// account
	accIdGen := adapt_domain_account.NewIdGeneratorImplUniqueRandom()
	accFactory := domain_account.NewFactory(accIdGen)
	accRepo := adapt_domain_account.NewRepositoryImplSql(sqlDbFnIsDuplicateKeyError)

	// action request
	actReqIdGen := adapt_domain_action_request.NewIdGeneratorImplUniqueRandom()
	actReqFactory := domain_action_request.NewFactory(actReqIdGen)
	actReqRepo := adapt_domain_action_request.NewRepositoryImplSql(sqlDbFnIsDuplicateKeyError)

	// cfm
	var cfmService domain_cfm.ServiceInterface = adapt_domain_cfm.NewServiceImplCfmService(
		infraCfmGrpcClient,
		func(ctx context.Context) string {
			traceId, _ := infraLogger.GetTraceIdFromCtx(ctx)
			return traceId
		},
	)

	// session
	sessIdGen := adapt_domain_session.NewIdGeneratorImplUniqueRandom()
	sessFactory := domain_session.NewFactory(sessIdGen)
	sessRepo := adapt_domain_session.NewRepositoryImplSql(sqlDbFnIsDuplicateKeyError)

	// -----------------------------------------------------------------------------------------------------------------
	// Opera
	// -----------------------------------------------------------------------------------------------------------------

	// txr
	operaTxr := txr.NewTxrImplSql(sqlDb, 2, 50*time.Millisecond, sqlDbFnIsDeadlockError)

	// logger
	operaLogger := adapt_opera_components.NewLoggerImplInfraLogger(infraLogger)

	// goroutiner
	operaGrt := goroutiner.New(
		goroutiner.MwPanicToError(func(panicValue any, debugStack []byte, ctx context.Context) error {
			infraLogger.CtxPanicFf(ctx, panicValue, debugStack, "container.operaGrt.MwPanicToError")
			// --
			var err error
			switch pv := panicValue.(type) {
			case error:
				err = fmt.Errorf("panic: %w; stack: %s", pv, string(debugStack))
			case string, fmt.Stringer:
				err = fmt.Errorf("panic: %q; stack: %s", pv, string(debugStack))
			default:
				err = fmt.Errorf("panic: %#v; stack: %s", pv, string(debugStack))
			}
			return std.WrapErrorToRuntime(err, "container.operaGrt", "MwPanicToError")
		}),
	)

	// domain facades
	accDomFac := opera_domain_facades.NewAccountDomFac(operaTxr, evStorage, accFactory, accRepo)
	actReqDomFac := opera_domain_facades.NewActionRequestDomFac(operaTxr, evStorage, actReqFactory, actReqRepo)
	sessDomFac := opera_domain_facades.NewSessionDomFac(operaTxr, evStorage, sessFactory, sessRepo)

	// -----------------------------------------------------------------------------------------------------------------

	return &Container{
		Logger: infraLogger,
		UseCases: UseCases{
			SignInRequest: sign_in_request.NewCommand(operaLogger, accDomFac, actReqDomFac, cfmService),
			SignInRequestRetry: sign_in_request_retry.NewCommand(
				operaLogger,
				operaGrt,
				accDomFac,
				actReqDomFac,
				cfmService,
				sessDomFac,
			),
			SignInConfirm: sign_in_confirm.NewCommand(
				operaLogger,
				operaGrt,
				accDomFac,
				actReqDomFac,
				cfmService,
				sessDomFac,
			),
			SignOut:      sign_out.NewCommand(operaLogger, accDomFac, sessDomFac),
			CheckSession: check_session.NewCommand(operaLogger, accDomFac, sessDomFac),
		},
	}
}
