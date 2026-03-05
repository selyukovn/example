package sign_in_request

import (
	"example/admin/auth/internal/domain/action_request"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Const
// ---------------------------------------------------------------------------------------------------------------------

var ResultNil = Result{}

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Result struct {
	signInId    action_request.Id
	retriesLeft uint
	canRetryAt  time.Time
	expireAt    time.Time
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// resultSuccess
//
// Паникует при нулевых аргументах:
//   - signInId
func resultSuccess(
	signInId action_request.Id,
	retriesLeft uint,
	canRetryAt time.Time,
	expireAt time.Time,
) (Result, error) {
	assert.Cmp[action_request.Id]().NotEq(action_request.IdNil).Must(signInId)

	return Result{
		signInId:    signInId,
		retriesLeft: retriesLeft,
		canRetryAt:  canRetryAt,
		expireAt:    expireAt,
	}, nil
}

func resultError(err error) (Result, error) {
	return ResultNil, err
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (r Result) IsNil() bool {
	return r == ResultNil
}

func (r Result) SignInId() action_request.Id {
	return r.signInId
}

// RetriesLeft
//
// Возвращает нулевое значение, если это была последняя попытка.
func (r Result) RetriesLeft() uint {
	return r.retriesLeft
}

// CanRetryAt
//
// Возвращает нулевое значение, если это была последняя попытка.
func (r Result) CanRetryAt() time.Time {
	return r.canRetryAt
}

func (r Result) ExpireAt() time.Time {
	return r.expireAt
}

// ---------------------------------------------------------------------------------------------------------------------
