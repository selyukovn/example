package account

import (
	"github.com/selyukovn/go-events"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type EventCreated struct {
	event.Event
	accId Id
	email std.Email
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewEventCreated
//
// Паникует при нулевых аргументах.
func NewEventCreated(occurredAt time.Time, accId Id, email std.Email) EventCreated {
	assert.Time().NotZero().Must(occurredAt)
	assert.Cmp[Id]().NotEq(IdNil).Must(accId)
	assert.Cmp[std.Email]().NotEq(std.EmailNil).Must(email)

	return EventCreated{
		Event: event.NewEvent(occurredAt),
		accId: accId,
		email: email,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (e EventCreated) AccId() Id {
	return e.accId
}

func (e EventCreated) Email() std.Email {
	return e.email
}

// ---------------------------------------------------------------------------------------------------------------------
