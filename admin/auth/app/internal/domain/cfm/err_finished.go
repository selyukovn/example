package cfm

import (
	"fmt"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type FinishType uint

const (
	FinishTypeNotYet  FinishType = 0
	FinishTypeExpired FinishType = 1
	FinishTypePassed  FinishType = 2
	FinishTypeFailed  FinishType = 3
)

var finishTypes = []FinishType{
	FinishTypeNotYet,
	FinishTypeExpired,
	FinishTypePassed,
	FinishTypeFailed,
}

type ErrorFinished struct {
	cfmId      Id
	finishedAt time.Time
	finishType FinishType
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func newErrorFinished(cfmId Id, finishedAt time.Time, finishType FinishType) ErrorFinished {
	assert.FalseMust(cfmId.IsNil())
	assert.Time().NotZero().Must(finishedAt)
	assert.Cmp[FinishType]().In(finishTypes).NotEq(FinishTypeNotYet).Must(finishType)

	return ErrorFinished{
		cfmId:      cfmId,
		finishedAt: finishedAt,
		finishType: finishType,
	}
}

// NewErrorFinishedAsExpired
//
// Паникует при нулевых аргументах.
func NewErrorFinishedAsExpired(cfmId Id, finishedAt time.Time) ErrorFinished {
	return newErrorFinished(cfmId, finishedAt, FinishTypeExpired)
}

// NewErrorFinishedAsFailed
//
// Паникует при нулевых аргументах.
func NewErrorFinishedAsFailed(cfmId Id, finishedAt time.Time) ErrorFinished {
	return newErrorFinished(cfmId, finishedAt, FinishTypeFailed)
}

// NewErrorFinishedAsPassed
//
// Паникует при нулевых аргументах.
func NewErrorFinishedAsPassed(cfmId Id, finishedAt time.Time) ErrorFinished {
	return newErrorFinished(cfmId, finishedAt, FinishTypePassed)
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (e ErrorFinished) Error() string {
	return fmt.Sprintf("конфирмация %q завершена", e.cfmId)
}

func (e ErrorFinished) CfmId() Id {
	return e.cfmId
}

func (e ErrorFinished) FinishedAt() time.Time {
	return e.finishedAt
}

func (e ErrorFinished) IsAsExpired() bool {
	return e.finishType == FinishTypeExpired
}

func (e ErrorFinished) IsAsFailed() bool {
	return e.finishType == FinishTypeFailed
}

func (e ErrorFinished) IsAsPassed() bool {
	return e.finishType == FinishTypePassed
}

// ---------------------------------------------------------------------------------------------------------------------
