package domain_facades

import (
	"context"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/action_request"
	"example/admin/auth/internal/domain/cfm"
	"example/admin/auth/internal/domain/event_storage"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-txr"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Const
// ---------------------------------------------------------------------------------------------------------------------

var ActionRequestDomFacNil = ActionRequestDomFac{}

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type ActionRequestDomFac struct {
	txr           txr.TxrInterface
	es            event_storage.Storage
	actReqFactory action_request.Factory
	actReqRepo    action_request.RepositoryInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewActionRequestDomFac(
	txr txr.TxrInterface,
	es event_storage.Storage,
	actReqFactory action_request.Factory,
	actReqRepo action_request.RepositoryInterface,
) ActionRequestDomFac {
	assert.NotNilDeepMust(txr)
	assert.Cmp[event_storage.Storage]().NotEq(event_storage.StorageNil).Must(es)
	assert.Cmp[action_request.Factory]().NotEq(action_request.FactoryNil).Must(actReqFactory)
	assert.Cmp[action_request.RepositoryInterface]().NotEq(nil).Must(actReqRepo)

	return ActionRequestDomFac{
		txr:           txr,
		es:            es,
		actReqFactory: actReqFactory,
		actReqRepo:    actReqRepo,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// CreateSignIn
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (f ActionRequestDomFac) CreateSignIn(
	ctx context.Context,
	accId account.Id,
	cfmId cfm.Id,
) (action_request.Id, error) {
	var actReqId action_request.Id
	err := f.txr.Tx(ctx, func(ctx context.Context) error {
		now := time.Now()
		signIn, err := f.actReqFactory.CreateSignIn(ctx, accId, cfmId, now)
		if err != nil {
			return err
		}

		err = f.actReqRepo.Add(ctx, signIn)
		switch err.(type) {
		case nil:
		case std.ErrorAlreadyDone:
			// AlreadyDone -- тоже бага: только что создали же
			return std.WrapErrorToRuntime(err, f, "CreateSignIn", "Add")
		default:
			return err
		}

		actReqId = signIn.Id()
		return nil
	})
	return actReqId, err
}

// CheckSignIn
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - std.ErrorRuntime
func (f ActionRequestDomFac) CheckSignIn(ctx context.Context, signInId action_request.Id) (account.Id, cfm.Id, error) {
	var accId account.Id
	var cfmId cfm.Id
	err := f.txr.Tx(ctx, func(ctx context.Context) error {
		signIn, err := f.actReqRepo.GetSignIn(ctx, signInId)
		if err != nil {
			return err
		}

		accId = signIn.AccId()
		cfmId = signIn.CfmId()
		return nil
	})
	return accId, cfmId, err
}

// ---------------------------------------------------------------------------------------------------------------------
