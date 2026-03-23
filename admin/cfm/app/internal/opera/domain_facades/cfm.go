package domain_facades

import (
	"context"
	"example/admin/cfm/internal/domain/cfm"
	"example/admin/cfm/internal/domain/cfm/code"
	"example/admin/cfm/internal/domain/event_storage"
	"github.com/selyukovn/go-events"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-txr"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Const
// ---------------------------------------------------------------------------------------------------------------------

var CfmDomFacNil = CfmDomFac{}

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type CfmDomFac struct {
	txr        txr.TxrInterface
	es         event_storage.Storage
	cfmFactory cfm.Factory
	cfmRepo    cfm.RepositoryInterface
	codeSender code.SenderInterface
	codeHasher code.HasherInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewCfmDomFac(
	txr txr.TxrInterface,
	es event_storage.Storage,
	cfmFactory cfm.Factory,
	cfmRepo cfm.RepositoryInterface,
	codeSender code.SenderInterface,
	codeHasher code.HasherInterface,
) CfmDomFac {
	assert.NotNilDeepMust(txr)
	assert.Cmp[event_storage.Storage]().NotEq(event_storage.StorageNil).Must(es)
	assert.Cmp[cfm.Factory]().NotEq(cfm.FactoryNil).Must(cfmFactory)
	assert.NotNilDeepMust(cfmRepo)
	assert.NotNilDeepMust(codeSender)
	assert.NotNilDeepMust(codeHasher)

	return CfmDomFac{
		txr:        txr,
		es:         es,
		cfmFactory: cfmFactory,
		cfmRepo:    cfmRepo,
		codeSender: codeSender,
		codeHasher: codeHasher,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// CreateForEmail
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
//
// Результат:
//   - Id
//   - ExpireAt
func (f CfmDomFac) CreateForEmail(ctx context.Context, email std.Email) (cfm.Id, time.Time, error) {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(email.IsNil())

	var cId cfm.Id
	var cExpireAt time.Time
	if err := f.txr.Tx(ctx, func(ctx context.Context) error {
		now := time.Now()

		c, err := f.cfmFactory.CreateEmailCfm(ctx, email, now)
		if err != nil {
			return err
		}

		err = f.cfmRepo.Add(ctx, c)
		switch err.(type) {
		case nil:
		case std.ErrorAlreadyDone:
			// AlreadyDone -- тоже бага: только что создали же
			return std.WrapErrorToRuntime(err, f, "CreateForEmail", "Add")
		default:
			return err
		}

		cId = c.Id()
		cExpireAt = c.ExpireAt()
		return nil
	}); err != nil {
		return cfm.IdNil, time.Time{}, err
	}

	return cId, cExpireAt, nil
}

// Request
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - cfm.ErrorFinished
//   - cfm.ErrorNoAttemptsLeft
//   - cfm.ErrorRequestsFrequency
//   - std.ErrorRuntime
//
// Результат:
//   - новый код
//   - email для отправки кода
//   - можно ли еще запросить
//   - можно ли еще запросить - если да, сколько раз
//   - можно ли еще запросить - если да, после какого момента времени
func (f CfmDomFac) Request(
	ctx context.Context,
	id cfm.Id,
) (
	rCode code.Code,
	rEmail std.Email,
	rCanReqAgain bool,
	rCanReqAttemptsLeft uint,
	rCanReqAfter time.Time,
	rErr error,
) {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(id.IsNil())

	rCode = code.CodeNil
	rEmail = std.EmailNil
	rCanReqAgain = false
	rCanReqAttemptsLeft = 0
	rCanReqAfter = time.Time{}
	rErr = f.txr.Tx(ctx, func(ctx context.Context) error {
		now := time.Now()

		c, err := f.cfmRepo.GetById(ctx, id)
		if err != nil {
			return err
		}

		cCode, cEmail, cCanReqAgain, cCanReqAttemptsLeft, cCanReqAfter, err := c.Request(now, f.cfmFactory, ctx)
		if err != nil {
			return err
		}

		err = f.cfmRepo.Update(ctx, c)
		if err != nil {
			return err
		}

		rCode = cCode
		rEmail = cEmail
		rCanReqAgain = cCanReqAgain
		rCanReqAttemptsLeft = cCanReqAttemptsLeft
		rCanReqAfter = cCanReqAfter
		return nil
	})

	return
}

// SendToEmail
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (f CfmDomFac) SendToEmail(ctx context.Context, cCode code.Code, email std.Email) error {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(cCode.IsNil())
	assert.FalseMust(email.IsNil())

	return f.codeSender.Send(ctx, cCode, email)
}

// Confirm
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - cfm.ErrorFinished
//   - std.ErrorUnprocessable -- если еще не запрашивалась (Request)
//   - std.ErrorRuntime
//
// Результат:
//   - если вызов завершил конфирмацию, время завершения
//   - если вызов завершил конфирмацию, успешно ли
//   - если вызов не завершил конфирмацию, сколько осталось попыток
func (f CfmDomFac) Confirm(
	ctx context.Context,
	id cfm.Id,
	cCode code.Code,
) (
	rFinishedAt time.Time,
	rIsFinishedAsPassed bool,
	rFailsLeft uint,
	rErr error,
) {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(id.IsNil())

	rFinishedAt = time.Time{}
	rIsFinishedAsPassed = false
	rFailsLeft = 0
	rErr = f.txr.Tx(ctx, func(ctx context.Context) error {
		now := time.Now()
		evs := event.NewCollection()

		c, err := f.cfmRepo.GetById(ctx, id)
		if err != nil {
			return err
		}

		finishedAt, isAsPassed, failsLeft, err := c.Confirm(now, evs, f.codeHasher, ctx, cCode)
		if err != nil {
			return err
		}

		err = f.cfmRepo.Update(ctx, c)
		if err != nil {
			return err
		}

		err = f.es.Store(ctx, evs)
		if err != nil {
			return err
		}

		rFinishedAt = finishedAt
		rIsFinishedAsPassed = isAsPassed
		rFailsLeft = failsLeft
		return nil
	})

	return
}

// TickTime
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - cfm.ErrorFinished -- завершенные не обновляются
//   - std.ErrorAlreadyDone -- когда нечего менять на данный момент
//   - std.ErrorRuntime
func (f CfmDomFac) TickTime(ctx context.Context, id cfm.Id) error {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(id.IsNil())

	return f.txr.Tx(ctx, func(ctx context.Context) error {
		now := time.Now()
		evs := event.NewCollection()

		c, err := f.cfmRepo.GetById(ctx, id)
		if err != nil {
			return err
		}

		err = c.TickTime(now, evs)
		if err != nil {
			return err
		}

		err = f.cfmRepo.Update(ctx, c)
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
// State
// ---------------------------------------------------------------------------------------------------------------------

// GetIdsGoingToExpire
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (f CfmDomFac) GetIdsGoingToExpire(ctx context.Context, limit uint) ([]cfm.Id, error) {
	assert.NotNilDeepMust(ctx)
	assert.Num[uint]().Positive().Must(limit)

	var cfmIds []cfm.Id
	err := f.txr.Tx(ctx, func(ctx context.Context) error {
		now := time.Now()

		ids, err := f.cfmRepo.GetIdsOfGoingToExpire(ctx, now, limit)
		if err != nil {
			return err
		}

		cfmIds = ids
		return nil
	})

	return cfmIds, err
}

// ---------------------------------------------------------------------------------------------------------------------
