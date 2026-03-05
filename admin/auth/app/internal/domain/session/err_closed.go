package session

import (
	"fmt"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type ErrorClosed struct {
	sessId    Id
	isExpired bool
	closedAt  time.Time
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewErrorSessionClosed
//
// Паникует при нулевых аргументах.
func NewErrorSessionClosed(sessId Id, closedAt time.Time, expireAt time.Time) ErrorClosed {
	assert.FalseMust(sessId.IsNil())
	assert.FalseMust(closedAt.IsZero())
	assert.FalseMust(expireAt.IsZero())

	return ErrorClosed{
		sessId:    sessId,
		closedAt:  closedAt,
		isExpired: closedAt == expireAt,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (e ErrorClosed) Error() string {
	return fmt.Sprintf("Сессия %q закрыта", e.sessId)
}

func (e ErrorClosed) SessId() Id {
	return e.sessId
}

func (e ErrorClosed) IsExpired() bool {
	return e.isExpired
}

func (e ErrorClosed) ClosedAt() time.Time {
	return e.closedAt
}

// ---------------------------------------------------------------------------------------------------------------------
