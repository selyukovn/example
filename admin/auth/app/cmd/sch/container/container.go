package container

import (
	"context"
	"database/sql"
	"example/admin/auth/cmd/common/components/processing"
	adapt_domain_event_storage "example/admin/auth/internal/adapt/domain/event_storage"
	adapt_domain_session "example/admin/auth/internal/adapt/domain/session"
	domain_event_storage "example/admin/auth/internal/domain/event_storage"
	domain_session "example/admin/auth/internal/domain/session"
	opera_domain_facades "example/admin/auth/internal/opera/domain_facades"
	"example/admin/auth/internal/opera/use_cases/session_tick_time"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	goroutiner "github.com/selyukovn/go-routiner"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"github.com/selyukovn/go-txr"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

type Container = struct {
	UseCases UseCases
}

type UseCases = struct {
	SessionTickTime session_tick_time.Command
}

func New(
	sqlDb *sql.DB,
	sqlDbFnIsDeadlockError func(error) bool,
	sqlDbFnIsDuplicateKeyError func(error) bool,
) *Container {
	assert.NotNilDeepMust(sqlDb)
	assert.NotNilDeepMust(sqlDbFnIsDeadlockError)
	assert.NotNilDeepMust(sqlDbFnIsDuplicateKeyError)

	// -----------------------------------------------------------------------------------------------------------------
	// Domain
	// -----------------------------------------------------------------------------------------------------------------

	// event storage
	evStorage := domain_event_storage.NewStorage(
		adapt_domain_event_storage.NewRepositoryImplSql(processing.OperationId),
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

	// goroutiner
	operaGrt := goroutiner.New(
		goroutiner.MwPanicToError(func(panicValue any, debugStack []byte, ctx context.Context) error {
			logger.PanicFf(ctx, panicValue, debugStack, "container.operaGrt.MwPanicToError")
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
	sessDomFac := opera_domain_facades.NewSessionDomFac(operaTxr, evStorage, sessFactory, sessRepo)

	// -----------------------------------------------------------------------------------------------------------------

	// Контейнер -- структура потенциально "растущая" (будут добавляться новые сервисы и т.д.).
	// Поэтому лучше сразу использовать контейнер через указатель.
	return &Container{
		UseCases: UseCases{
			SessionTickTime: session_tick_time.NewCommand(operaGrt, sessDomFac),
		},
	}
}
