package cfm

import (
	"fmt"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type ErrorRequestsFrequency struct {
	cfmId              Id
	canReqAfter        time.Time
	canReqAttemptsLeft uint
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewErrorRequestsFrequency
//
// Паникует при нулевых аргументах.
func NewErrorRequestsFrequency(cfmId Id, canReqAfter time.Time, canReqAttemptsLeft uint) ErrorRequestsFrequency {
	assert.FalseMust(cfmId.IsNil())
	assert.FalseMust(canReqAfter.IsZero())
	assert.FalseMust(canReqAttemptsLeft == 0)

	return ErrorRequestsFrequency{
		cfmId:              cfmId,
		canReqAfter:        canReqAfter,
		canReqAttemptsLeft: canReqAttemptsLeft,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (e ErrorRequestsFrequency) Error() string {
	return fmt.Sprintf("конфирмация %q неперезапрашиваема до %s", e.cfmId, e.canReqAfter.Format(time.RFC3339))
}

func (e ErrorRequestsFrequency) CfmId() Id {
	return e.cfmId
}

func (e ErrorRequestsFrequency) CanReqAfter() time.Time {
	return e.canReqAfter
}

func (e ErrorRequestsFrequency) CanReqAttemptsLeft() uint {
	return e.canReqAttemptsLeft
}

// ---------------------------------------------------------------------------------------------------------------------
