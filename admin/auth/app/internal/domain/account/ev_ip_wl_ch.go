package account

import (
	"github.com/selyukovn/go-events"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type EventIpWhitelistChanged struct {
	event.Event
	accId   Id
	newList IpWhitelist
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewEventIpWhitelistChanged
//
// Паникует при нулевых аргументах:
//   - occurredAt
//   - accId
func NewEventIpWhitelistChanged(
	occurredAt time.Time,
	accId Id,
	newList IpWhitelist,
) EventIpWhitelistChanged {
	assert.Time().NotZero().Must(occurredAt)
	assert.Cmp[Id]().NotEq(IdNil).Must(accId)

	return EventIpWhitelistChanged{
		Event:   event.NewEvent(occurredAt),
		accId:   accId,
		newList: newList,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (e EventIpWhitelistChanged) AccId() Id {
	return e.accId
}

func (e EventIpWhitelistChanged) NewList() IpWhitelist {
	return e.newList
}

// ---------------------------------------------------------------------------------------------------------------------
