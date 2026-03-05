package sign_in_request_retry

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
	signInId           action_request.Id
	canReqAgain        bool
	canReqAttemptsLeft uint
	canReqAfter        time.Time
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// resultSuccessCanAgain
//
// Паникует при нулевых аргументах.
func resultSuccessCanAgain(signInId action_request.Id, canReqAttemptsLeft uint, canReqAfter time.Time) (Result, error) {
	assert.FalseMust(signInId.IsNil())
	assert.FalseMust(canReqAttemptsLeft == 0)
	assert.FalseMust(canReqAfter.IsZero())

	return Result{
		signInId:           signInId,
		canReqAgain:        true,
		canReqAttemptsLeft: canReqAttemptsLeft,
		canReqAfter:        canReqAfter,
	}, nil
}

// resultSuccessLast
//
// Паникует при нулевых аргументах.
func resultSuccessLast(signInId action_request.Id) (Result, error) {
	assert.FalseMust(signInId.IsNil())

	return Result{
		signInId:           signInId,
		canReqAgain:        false,
		canReqAttemptsLeft: 0,
		canReqAfter:        time.Time{},
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

func (r Result) CanReqAgain() bool {
	return r.canReqAgain
}

func (r Result) RetriesLeft() uint {
	return r.canReqAttemptsLeft
}

func (r Result) CanRetryAt() time.Time {
	return r.canReqAfter
}

// ---------------------------------------------------------------------------------------------------------------------
