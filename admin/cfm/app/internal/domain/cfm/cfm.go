package cfm

import (
	"context"
	"example/admin/cfm/internal/domain/cfm/code"
	"github.com/selyukovn/go-events"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

const (
	requestsMax             = 3
	requestsThresholdPeriod = 75 * time.Second
	confirmFailsMax         = 3
	expirePeriod            = (requestsMax * requestsThresholdPeriod) + (confirmFailsMax * 10 * time.Second)
)

type FinishType uint

const (
	FinishNotYet  FinishType = 0
	FinishExpired FinishType = 1
	FinishPassed  FinishType = 2
	FinishFailed  FinishType = 3
)

type Cfm struct {
	id         Id
	email      std.Email // емайл / телефон / ... тут, потому что при повторных запросах должна отправляться туда же.
	expireAt   time.Time
	finishedAt time.Time
	finishType FinishType
	requests   []*CfmRequest
	failsMade  uint
}

type CfmRequest = struct {
	Number      uint
	CodeHash    code.Hash
	RequestedAt time.Time
}

// ---------------------------------------------------------------------------------------------------------------------
// Reflect
// ---------------------------------------------------------------------------------------------------------------------

func ReflectExtract(c *Cfm) (
	id Id,
	email std.Email,
	expireAt time.Time,
	finishedAt time.Time,
	finishType FinishType,
	requests []*CfmRequest,
	failsMade uint,
) {
	id = c.id
	email = c.email
	expireAt = c.expireAt
	finishedAt = c.finishedAt
	finishType = c.finishType
	requests = c.requests
	failsMade = c.failsMade
	return
}

func ReflectRestore(
	id Id,
	email std.Email,
	expireAt time.Time,
	finishedAt time.Time,
	finishType FinishType,
	requests []*CfmRequest,
	failsMade uint,
) *Cfm {
	return &Cfm{
		id:         id,
		email:      email,
		expireAt:   expireAt,
		finishedAt: finishedAt,
		finishType: finishType,
		requests:   requests,
		failsMade:  failsMade,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// См. Factory.CreateEmailCfm
var _ = Factory{}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (c *Cfm) _makeErrorFinished(now time.Time) ErrorFinished {
	finishedAt := c.finishedAt
	finishType := c.finishType

	if c.isFinishedAsGoingToExpire(now) {
		finishedAt = c.expireAt
		finishType = FinishExpired
	}

	return NewErrorFinished(c.id, finishedAt, finishType)
}

// Request
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - ErrorFinished
//   - ErrorNoAttemptsLeft
//   - ErrorRequestsFrequency
//   - std.ErrorRuntime
//
// Результат:
//   - новый код
//   - email для отправки кода
//   - можно ли еще запросить
//   - можно ли еще запросить - если да, сколько раз
//   - можно ли еще запросить - если да, после какого момента времени
func (c *Cfm) Request(
	now time.Time,
	factory Factory,
	ctx context.Context,
) (
	rCode code.Code,
	rEmail std.Email,
	rCanReqAgain bool,
	rCanReqAttemptsLeft uint,
	rCanReqAfter time.Time,
	rErr error,
) {
	// По поводу `factory` -- да, ссылочная непрозрачность.
	// Но выносить проверки в какой-то "AssertCanRequest(now) error" метод перед генерацией
	// или, вероятно, зря генерировать код с хешем до проверок, кмк, большее зло.

	assert.FalseMust(now.IsZero())
	assert.Cmp[Factory]().NotEq(FactoryNil).Must(factory)
	assert.NotNilDeepMust(ctx)

	rCode = code.CodeNil
	rEmail = std.EmailNil
	rCanReqAgain = false
	rCanReqAttemptsLeft = 0
	rCanReqAfter = time.Time{}

	if c.IsFinished(now) {
		rErr = c._makeErrorFinished(now)
		return
	}

	reqCount := len(c.requests)
	reqLeft := requestsMax - reqCount

	if reqLeft == 0 {
		rErr = NewErrNoAttemptsLeft(c.id)
		return
	}

	if reqCount > 0 {
		canReqAfter := c.requests[reqCount-1].RequestedAt.Add(requestsThresholdPeriod)
		if canReqAfter.After(now) {
			// тут не имеет значения, протухнет к тому времени конфирмация или нет (т.е. canReqAgain = true)
			// Попытка потом Request протухшей упадет на IsFinished -- т.е. не "наши" проблемы.
			rErr = NewErrorRequestsFrequency(c.id, canReqAfter, uint(reqLeft))
			return
		}
	}

	// --

	cc, cch, err := factory.GenerateCodeAndHash(ctx)
	if err != nil {
		rErr = std.WrapErrorToRuntime(err, c, "Request", "GenerateCodeAndHash")
		return
	}

	reqLeft = reqLeft - 1
	reqCount = reqCount + 1
	c.requests = append(c.requests, &CfmRequest{
		Number:      uint(reqCount),
		CodeHash:    cch,
		RequestedAt: now,
	})
	canReqAfter := c.requests[reqCount-1].RequestedAt.Add(requestsThresholdPeriod)

	// --

	rCode = cc
	rEmail = c.email
	if canReqAgain := reqLeft > 0 && canReqAfter.Before(c.expireAt); canReqAgain {
		rCanReqAgain = true
		rCanReqAttemptsLeft = uint(reqLeft)
		rCanReqAfter = canReqAfter
	}
	rErr = nil
	return
}

// Confirm
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - ErrorFinished
//   - std.ErrorUnprocessable -- если еще не запрашивалась (Request)
//   - std.ErrorRuntime
//
// Результат:
//   - если вызов завершил конфирмацию, время завершения
//   - если вызов завершил конфирмацию, успешно ли
//   - если вызов не завершил конфирмацию, сколько осталось попыток
func (c *Cfm) Confirm(
	now time.Time,
	evs *event.Collection,
	codeHasher code.HasherInterface,
	ctx context.Context,
	code code.Code,
) (
	rFinishedAt time.Time,
	rIsFinishedAsPassed bool,
	rFailsLeft uint,
	rErr error,
) {
	// По поводу codeHasher -- да, ссылочная непрозрачность.
	// Инкапсуляция здесь кажется более важным свойством.

	assert.FalseMust(now.IsZero())
	assert.NotNilDeepMust(evs)
	assert.NotNilDeepMust(codeHasher)
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(code.IsNil())

	rFinishedAt = time.Time{}
	rIsFinishedAsPassed = false
	rFailsLeft = 0

	if c.IsFinished(now) {
		rErr = c._makeErrorFinished(now)
		return
	}

	reqCount := len(c.requests)
	if reqCount == 0 {
		rErr = std.NewErrorUnprocessableFf("Конфирмация %q еще не запрашивалась", c.id)
		return
	}

	lastCodeHash := c.requests[reqCount-1].CodeHash
	isPassed, err := codeHasher.Compare(ctx, code, lastCodeHash)
	if err != nil {
		rErr = std.WrapErrorToRuntime(err, c, "Confirm", "Compare")
		return
	}

	if isPassed {
		c.finishType = FinishPassed
		c.finishedAt = now
		evs.Add(NewEventFinished(now, c.id, c.finishType))
	} else {
		c.failsMade += 1

		if c.failsMade == confirmFailsMax {
			c.finishedAt = now
			c.finishType = FinishFailed
			evs.Add(NewEventFinished(c.finishedAt, c.id, c.finishType))
		}
	}

	rFinishedAt = c.finishedAt
	rIsFinishedAsPassed = isPassed
	if c.finishedAt.IsZero() {
		rFailsLeft = confirmFailsMax - c.failsMade
	}
	rErr = nil
	return
}

// TickTime
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - ErrorFinished
//   - std.ErrorAlreadyDone -- когда нечего менять на данный момент
func (c *Cfm) TickTime(now time.Time, evs *event.Collection) error {
	assert.FalseMust(now.IsZero())
	assert.NotNilDeepMust(evs)

	if c.isFinishedAsGoingToExpire(now) {
		c.finishedAt = c.expireAt
		c.finishType = FinishExpired
		evs.Add(NewEventFinished(c.finishedAt, c.id, c.finishType))
		return nil
	} else if c.IsFinished(now) {
		return NewErrorFinished(c.id, c.finishedAt, c.finishType)
	}

	return std.NewErrorAlreadyDoneFf("Нечего менять для конфирмации %q на момент времени %s", c.id, now)
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (c *Cfm) Id() Id {
	return c.id
}

func (c *Cfm) ExpireAt() time.Time {
	return c.expireAt
}

func (c *Cfm) Email() std.Email {
	return c.email
}

// isFinishedAsGoingToExpire
//
// Нельзя просто так взять и завершить точно в момент expireAt
// (технически возможно, например, по таймеру в отдельной горутине, но это гораздо накладнее, чем "кроном"),
// поэтому между моментом expireAt и фактическим finishedAt может пройти какое-то время.
// Это нужно учитывать также в поисковых запросах вроде RepositoryInterface.GetIdsOfGoingToExpire
func (c *Cfm) isFinishedAsGoingToExpire(now time.Time) bool {
	return c.finishType == FinishNotYet && now.After(c.expireAt)
}

func (c *Cfm) IsFinished(now time.Time) bool {
	return c.finishType != FinishNotYet || c.isFinishedAsGoingToExpire(now)
}

// ---------------------------------------------------------------------------------------------------------------------
