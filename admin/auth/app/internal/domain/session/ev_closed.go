package session

import (
	"example/admin/auth/internal/domain/account"
	"github.com/selyukovn/go-events"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type EventClosed struct {
	event.Event
	sessId Id
	accId  account.Id
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewEventClosed
//
// Паникует при нулевых аргументах.
func NewEventClosed(occurredAt time.Time, sessId Id, accId account.Id) EventClosed {
	assert.FalseMust(occurredAt.IsZero())
	assert.FalseMust(sessId.IsNil())
	assert.FalseMust(accId.IsNil())

	return EventClosed{
		Event:  event.NewEvent(occurredAt),
		sessId: sessId,
		accId:  accId,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (e EventClosed) SessId() Id {
	return e.sessId
}

func (e EventClosed) AccId() account.Id {
	return e.accId
}

// ---------------------------------------------------------------------------------------------------------------------
