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

type EventCreated struct {
	event.Event
	sessId Id
	accId  account.Id
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewEventCreated
//
// Паникует при нулевых аргументах.
func NewEventCreated(occurredAt time.Time, sessId Id, accId account.Id) EventCreated {
	assert.FalseMust(occurredAt.IsZero())
	assert.FalseMust(sessId.IsNil())
	assert.FalseMust(accId.IsNil())

	return EventCreated{
		Event:  event.NewEvent(occurredAt),
		sessId: sessId,
		accId:  accId,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (e EventCreated) SessId() Id {
	return e.sessId
}

func (e EventCreated) AccId() account.Id {
	return e.accId
}

// ---------------------------------------------------------------------------------------------------------------------
