package request

import (
	"example/admin/cfm/internal/domain/cfm"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Const
// ---------------------------------------------------------------------------------------------------------------------

var ResultNil = Result{}

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Result struct {
	cfmId          cfm.Id
	canReqAgain    bool
	canReqAttempts uint
	canReqAfter    time.Time
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// resultSuccessCanAgain
//
// Паникует при нулевых аргументах.
func resultSuccessCanAgain(
	cfmId cfm.Id,
	canReqAttemptsLeft uint,
	canReqAfter time.Time,
) (Result, error) {
	assert.FalseMust(cfmId.IsNil())
	assert.FalseMust(canReqAttemptsLeft == 0)
	assert.FalseMust(canReqAfter.IsZero())

	return Result{
		cfmId:          cfmId,
		canReqAgain:    true,
		canReqAttempts: canReqAttemptsLeft,
		canReqAfter:    canReqAfter,
	}, nil
}

// resultSuccessLast
//
// Паникует при нулевых аргументах.
func resultSuccessLast(cfmId cfm.Id) (Result, error) {
	assert.FalseMust(cfmId.IsNil())

	return Result{
		cfmId:          cfmId,
		canReqAgain:    false,
		canReqAttempts: 0,
		canReqAfter:    time.Time{},
	}, nil
}

// resultError
//
// Паникует при нулевых аргументах.
func resultError(err error) (Result, error) {
	assert.NotNilDeepMust(err)

	return ResultNil, err
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (r Result) IsNil() bool {
	return r == ResultNil
}

func (r Result) CfmId() cfm.Id {
	return r.cfmId
}

func (r Result) CanReqAgain() bool {
	return r.canReqAgain
}

// CanReqAttemptsLeft
//
// Возвращает нулевое значение, если не CanReqAgain()
func (r Result) CanReqAttemptsLeft() uint {
	return r.canReqAttempts
}

// CanReqAfter
//
// Возвращает нулевое значение, если не CanReqAgain()
func (r Result) CanReqAfter() time.Time {
	return r.canReqAfter
}

// ---------------------------------------------------------------------------------------------------------------------
