package domain_facades

import (
	"context"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/action_request"
	"example/admin/auth/internal/domain/client"
	"example/admin/auth/internal/domain/event_storage"
	"example/admin/auth/internal/domain/session"
	"github.com/selyukovn/go-events"
	"github.com/selyukovn/go-txr"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type SessionDomFac struct {
	txr         txr.TxrInterface
	es          *event_storage.Storage
	sessFactory *session.Factory
	sessRepo    session.RepositoryInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewSessionDomFac
//
// Паникует при нулевых аргументах.
func NewSessionDomFac(
	txr txr.TxrInterface,
	es *event_storage.Storage,
	sessFactory *session.Factory,
	sessRepo session.RepositoryInterface,
) *SessionDomFac {
	assert.NotNilDeepMust(txr)
	assert.Cmp[*event_storage.Storage]().NotEq(nil).Must(es)
	assert.Cmp[*session.Factory]().NotEq(nil).Must(sessFactory)
	assert.Cmp[session.RepositoryInterface]().NotEq(nil).Must(sessRepo)

	return &SessionDomFac{
		txr:         txr,
		es:          es,
		sessFactory: sessFactory,
		sessRepo:    sessRepo,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Create
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorAlreadyDone
//   - std.ErrorRuntime
func (f *SessionDomFac) Create(
	ctx context.Context,
	cl client.Client,
	accId account.Id,
	signInId action_request.Id,
) (session.Id, time.Time, error) {
	var sessId session.Id
	var sessExpAt time.Time
	err := f.txr.Tx(ctx, func(ctx context.Context) error {
		now := time.Now()
		evs := event.NewCollection()

		sess, err := f.sessFactory.Create(ctx, accId, signInId, cl, now, evs)
		if err != nil {
			return err
		}

		err = f.sessRepo.Add(ctx, sess)
		if err != nil {
			return err
		}

		err = f.es.Store(ctx, evs)
		if err != nil {
			return err
		}

		sessId = sess.Id()
		sessExpAt = sess.ExpireAt()
		return nil
	})
	return sessId, sessExpAt, err
}

// HasBySignInRequest
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (f *SessionDomFac) HasBySignInRequest(ctx context.Context, signInId action_request.Id) (bool, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[action_request.Id]().NotEq(action_request.IdNil).Must(signInId)

	var isSessionExist bool
	var err error
	_ = f.txr.Tx(ctx, func(ctx context.Context) error {
		isSessionExist, err = f.sessRepo.HasBySignInRequest(ctx, signInId)
		return nil
	})

	return isSessionExist, err
}

// GetAccId
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - std.ErrorRuntime
func (f *SessionDomFac) GetAccId(ctx context.Context, sessId session.Id) (account.Id, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[session.Id]().NotEq(session.IdNil).Must(sessId)

	var accId account.Id
	err := f.txr.Tx(ctx, func(ctx context.Context) error {
		_accId, err := f.sessRepo.GetAccIdById(ctx, sessId)
		accId = _accId
		return err
	})
	return accId, err
}

// GetAccIdAndExpAt
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - std.ErrorRuntime
func (f *SessionDomFac) GetAccIdAndExpAt(ctx context.Context, sessId session.Id) (account.Id, time.Time, error) {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(sessId.IsNil())

	var accId account.Id
	var expAt time.Time
	err := f.txr.Tx(ctx, func(ctx context.Context) error {
		_accId, _expAt, err := f.sessRepo.GetAccIdAndExpAtById(ctx, sessId)
		accId = _accId
		expAt = _expAt
		return err
	})
	return accId, expAt, err
}

// GetIdsGoingToExpire
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (f *SessionDomFac) GetIdsGoingToExpire(ctx context.Context, limit uint) ([]session.Id, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Num[uint]().Positive().Must(limit)

	var sessIds []session.Id
	err := f.txr.Tx(ctx, func(ctx context.Context) error {
		now := time.Now()

		ids, err := f.sessRepo.GetIdsOfGoingToExpire(ctx, now, limit)
		if err != nil {
			return err
		}

		sessIds = ids
		return nil
	})

	return sessIds, err
}

// TickTime
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - session.ErrorClosed -- закрытые сессии не обновляются
//   - std.ErrorAlreadyDone -- когда нечего менять на данный момент
//   - std.ErrorRuntime
func (f *SessionDomFac) TickTime(ctx context.Context, sessId session.Id) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[session.Id]().NotEq(session.IdNil).Must(sessId)

	return f.txr.Tx(ctx, func(ctx context.Context) error {
		now := time.Now()
		evs := event.NewCollection()

		sess, err := f.sessRepo.GetById(ctx, sessId)
		if err != nil {
			return err
		}

		err = sess.TickTime(now, evs)
		if err != nil {
			return err
		}

		err = f.sessRepo.Update(ctx, sess)
		if err != nil {
			return err
		}

		err = f.es.Store(ctx, evs)
		if err != nil {
			return err
		}

		return nil
	})
}

// Close
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - session.ErrorClosed
//   - std.ErrorRuntime
func (f *SessionDomFac) Close(ctx context.Context, sessId session.Id) error {
	return f.txr.Tx(ctx, func(ctx context.Context) error {
		now := time.Now()
		evs := event.NewCollection()

		sess, err := f.sessRepo.GetById(ctx, sessId)
		if err != nil {
			return err
		}

		err = sess.Close(now, evs)
		if err != nil {
			return err
		}

		err = f.sessRepo.Update(ctx, sess)
		if err != nil {
			return err
		}

		err = f.es.Store(ctx, evs)
		if err != nil {
			return err
		}

		return nil
	})
}

// ---------------------------------------------------------------------------------------------------------------------
