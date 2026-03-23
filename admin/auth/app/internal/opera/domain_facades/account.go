package domain_facades

import (
	"context"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/client"
	"example/admin/auth/internal/domain/event_storage"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-txr"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Const
// ---------------------------------------------------------------------------------------------------------------------

var AccountDomFacNil = AccountDomFac{}

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type AccountDomFac struct {
	txr        txr.TxrInterface
	es         event_storage.Storage
	accFactory account.Factory
	accRepo    account.RepositoryInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewAccountDomFac
//
// Паникует при нулевых аргументах.
func NewAccountDomFac(
	txr txr.TxrInterface,
	es event_storage.Storage,
	accFactory account.Factory,
	accRepo account.RepositoryInterface,
) AccountDomFac {
	assert.NotNilDeepMust(txr)
	assert.Cmp[event_storage.Storage]().NotEq(event_storage.StorageNil).Must(es)
	assert.Cmp[account.Factory]().NotEq(account.FactoryNil).Must(accFactory)
	assert.Cmp[account.RepositoryInterface]().NotEq(nil).Must(accRepo)

	return AccountDomFac{
		txr:        txr,
		es:         es,
		accFactory: accFactory,
		accRepo:    accRepo,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// CanSignIn
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - account.ErrorDeactivated
//   - account.ErrorIpWhitelist
//   - std.ErrorRuntime
func (f AccountDomFac) CanSignIn(ctx context.Context, cl client.Client, email std.Email) (account.Id, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[client.Client]().NotEq(client.ClientNil).Must(cl)
	assert.Cmp[std.Email]().NotEq(std.EmailNil).Must(email)

	accId := account.IdNil
	err := f.txr.Tx(ctx, func(ctx context.Context) error {
		now := time.Now()

		acc, err := f.accRepo.GetByEmail(ctx, email)
		if err != nil {
			return err
		}

		err = acc.AssertSignIn(cl, now)
		if err != nil {
			return err
		}

		accId = acc.Id()
		return nil
	})

	return accId, err
}

// CheckAccess
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - account.ErrorDeactivated
//   - account.ErrorIpWhitelist
//   - std.ErrorRuntime
func (f AccountDomFac) CheckAccess(ctx context.Context, cl client.Client, id account.Id) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[client.Client]().NotEq(client.ClientNil).Must(cl)
	assert.Cmp[account.Id]().NotEq(account.IdNil).Must(id)

	return f.txr.Tx(ctx, func(ctx context.Context) error {
		now := time.Now()

		acc, err := f.accRepo.GetById(ctx, id)
		if err != nil {
			return err
		}

		err = acc.AssertSignIn(cl, now)
		if err != nil {
			return err
		}

		return nil
	})
}

// ---------------------------------------------------------------------------------------------------------------------
