package cfm

import (
	"fmt"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type ErrorFinished struct {
	cfmId      Id
	finishedAt time.Time
	finishType FinishType
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewErrorFinished
//
// Паникует при нулевых аргументах.
func NewErrorFinished(
	cfmId Id,
	finishedAt time.Time,
	finishType FinishType,
) ErrorFinished {
	assert.FalseMust(cfmId.IsNil())
	assert.Time().NotZero().Must(finishedAt)
	assert.Cmp[FinishType]().NotEq(FinishNotYet).Must(finishType)

	return ErrorFinished{
		cfmId:      cfmId,
		finishedAt: finishedAt,
		finishType: finishType,
	}
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
	return e.finishType == FinishExpired
}

func (e ErrorFinished) IsAsPassed() bool {
	return e.finishType == FinishPassed
}

func (e ErrorFinished) IsAsFailed() bool {
	return e.finishType == FinishFailed
}

// ---------------------------------------------------------------------------------------------------------------------
