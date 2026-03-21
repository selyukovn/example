package container

import (
	"context"
	"database/sql"
	adapt_domain_cfm "example/admin/cfm/internal/adapt/domain/cfm"
	adapt_domain_cfm_code "example/admin/cfm/internal/adapt/domain/cfm/code"
	adapt_domain_event_storage "example/admin/cfm/internal/adapt/domain/event_storage"
	domain_cfm "example/admin/cfm/internal/domain/cfm"
	domain_event_storage "example/admin/cfm/internal/domain/event_storage"
	opera_domain_facades "example/admin/cfm/internal/opera/domain_facades"
	"example/admin/cfm/internal/opera/use_cases/confirm"
	"example/admin/cfm/internal/opera/use_cases/create_for_email"
	"example/admin/cfm/internal/opera/use_cases/request"
	"example/admin/cfm/internal/opera/use_cases/tick_time"
	"fmt"
	goroutiner "github.com/selyukovn/go-routiner"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	"github.com/selyukovn/go-txr"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

type Container struct {
	UseCases UseCases
}

type UseCases = struct {
	CreateForEmail *create_for_email.Command
	Request        *request.Command
	Confirm        *confirm.Command
	TickTime       *tick_time.Command
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
		adapt_domain_event_storage.NewRepositoryImplSql(),
	)

	// cfm
	cfmCodeGen := adapt_domain_cfm_code.NewGeneratorImplUintRand1()
	cfmCodeHasher := adapt_domain_cfm_code.NewHasherImplBcrypt10()
	cfmCodeSender := adapt_domain_cfm_code.NewSenderImplDummy()
	cfmIdGen := adapt_domain_cfm.NewIdGeneratorImplUniqueRandom()
	cfmFactory := domain_cfm.NewFactory(cfmIdGen, cfmCodeGen, cfmCodeHasher)
	cfmRepo := adapt_domain_cfm.NewRepositoryImplSql(sqlDbFnIsDuplicateKeyError)

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
	cfmDomFac := opera_domain_facades.NewCfmDomFac(
		operaTxr,
		evStorage,
		cfmFactory,
		cfmRepo,
		cfmCodeSender,
		cfmCodeHasher,
	)

	// -----------------------------------------------------------------------------------------------------------------

	return &Container{
		UseCases: UseCases{
			CreateForEmail: create_for_email.NewCommand(cfmDomFac),
			Request:        request.NewCommand(operaGrt, cfmDomFac),
			Confirm:        confirm.NewCommand(cfmDomFac),
			TickTime:       tick_time.NewCommand(operaGrt, cfmDomFac),
		},
	}
}
