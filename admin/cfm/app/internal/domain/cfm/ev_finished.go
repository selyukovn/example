package cfm

import (
	"github.com/selyukovn/go-events"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

type EventFinished struct {
	event.Event
	cfmId      Id
	finishType FinishType
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewEventFinished
//
// Паникует при нулевых аргументах.
func NewEventFinished(
	occurredAt time.Time,
	cfmId Id,
	finishType FinishType,
) EventFinished {
	assert.Time().NotZero().Must(occurredAt)
	assert.FalseMust(cfmId.IsNil())
	assert.Cmp[FinishType]().NotEq(FinishNotYet).Must(finishType)

	return EventFinished{
		Event:      event.NewEvent(occurredAt),
		cfmId:      cfmId,
		finishType: finishType,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (e EventFinished) CfmId() Id {
	return e.cfmId
}

func (e EventFinished) FinishType() FinishType {
	return e.finishType
}

func (e EventFinished) IsExpired() bool {
	return e.finishType == FinishExpired
}

func (e EventFinished) IsPassed() bool {
	return e.finishType == FinishPassed
}

func (e EventFinished) IsFailed() bool {
	return e.finishType == FinishFailed
}

// ---------------------------------------------------------------------------------------------------------------------
