package cfm

import (
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Const
// ---------------------------------------------------------------------------------------------------------------------

var ServiceResultConfirmNil = ServiceResultConfirm{}

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type ServiceResultConfirm struct {
	cfmId              Id
	finishedAt         time.Time
	isFinishedAsPassed bool
	failsLeft          uint
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewServiceResultConfirmCanAgain
//
// Паникует при нулевых аргументах.
func NewServiceResultConfirmCanAgain(cfmId Id, failsLeft uint) ServiceResultConfirm {
	assert.FalseMust(cfmId.IsNil())
	assert.FalseMust(failsLeft == 0)

	return ServiceResultConfirm{
		cfmId:              cfmId,
		finishedAt:         time.Time{},
		isFinishedAsPassed: false,
		failsLeft:          failsLeft,
	}
}

// NewServiceResultConfirmLast
//
// Паникует при нулевых аргументах:
//   - cfmId
//   - finishedAt
func NewServiceResultConfirmLast(cfmId Id, finishedAt time.Time, isFinishedAsPassed bool) ServiceResultConfirm {
	assert.FalseMust(cfmId.IsNil())
	assert.FalseMust(finishedAt.IsZero())

	return ServiceResultConfirm{
		cfmId:              cfmId,
		finishedAt:         finishedAt,
		isFinishedAsPassed: isFinishedAsPassed,
		failsLeft:          0,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (r ServiceResultConfirm) IsNil() bool {
	return r == ServiceResultConfirmNil
}

func (r ServiceResultConfirm) CfmId() Id {
	return r.cfmId
}

func (r ServiceResultConfirm) IsFinished() bool {
	return !r.finishedAt.IsZero()
}

func (r ServiceResultConfirm) IsFinishedAsFailed() bool {
	return r.IsFinished() && !r.isFinishedAsPassed
}

func (r ServiceResultConfirm) IsFinishedAsPassed() bool {
	return r.IsFinished() && r.isFinishedAsPassed
}

// FinishedAt
//
// Возвращает нулевое значение, если не IsFinished()
func (r ServiceResultConfirm) FinishedAt() time.Time {
	return r.finishedAt
}

// FailsLeft
//
// Возвращает нулевое значение, если IsFinished()
func (r ServiceResultConfirm) FailsLeft() uint {
	return r.failsLeft
}

// ---------------------------------------------------------------------------------------------------------------------
