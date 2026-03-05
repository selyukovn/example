package confirm

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
	cfmId              cfm.Id
	finishedAt         time.Time
	isFinishedAsPassed bool
	failsLeft          uint
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// resultSuccessCanAgain
//
// Паникует при нулевых аргументах.
func resultSuccessCanAgain(cfmId cfm.Id, failsLeft uint) (Result, error) {
	assert.FalseMust(cfmId.IsNil())
	assert.FalseMust(failsLeft == 0)

	return Result{
		cfmId:              cfmId,
		finishedAt:         time.Time{},
		isFinishedAsPassed: false,
		failsLeft:          failsLeft,
	}, nil
}

// resultSuccessLast
//
// Паникует при нулевых аргументах:
//   - cfmId
//   - finishedAt
func resultSuccessLast(cfmId cfm.Id, finishedAt time.Time, isFinishedAsPassed bool) (Result, error) {
	assert.FalseMust(cfmId.IsNil())
	assert.FalseMust(finishedAt.IsZero())

	return Result{
		cfmId:              cfmId,
		finishedAt:         finishedAt,
		isFinishedAsPassed: isFinishedAsPassed,
		failsLeft:          0,
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

func (r Result) IsFinished() bool {
	return !r.finishedAt.IsZero()
}

func (r Result) IsFinishedAsFailed() bool {
	return r.IsFinished() && !r.isFinishedAsPassed
}

func (r Result) IsFinishedAsPassed() bool {
	return r.IsFinished() && r.isFinishedAsPassed
}

// FinishedAt
//
// Возвращает нулевое значение, если не IsFinished()
func (r Result) FinishedAt() time.Time {
	return r.finishedAt
}

// FailsLeft
//
// Возвращает нулевое значение, если IsFinished()
func (r Result) FailsLeft() uint {
	return r.failsLeft
}

// ---------------------------------------------------------------------------------------------------------------------
