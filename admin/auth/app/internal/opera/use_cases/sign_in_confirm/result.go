package sign_in_confirm

import (
	"example/admin/auth/internal/domain/session"
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
	sessId       session.Id
	sessExpAt    time.Time
	attemptsLeft uint
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// resultSuccess
//
// Паникует при нулевых аргументах.
func resultSuccess(sessId session.Id, sessExpAt time.Time) (Result, error) {
	assert.Cmp[session.Id]().NotEq(session.IdNil).Must(sessId)
	assert.Time().NotZero().Must(sessExpAt)

	return Result{
		sessId:       sessId,
		sessExpAt:    sessExpAt,
		attemptsLeft: 0,
	}, nil
}

func resultFail(attemptsLeft uint) (Result, error) {
	return Result{
		sessId:       session.IdNil,
		sessExpAt:    time.Time{},
		attemptsLeft: attemptsLeft,
	}, nil
}

func resultError(err error) (Result, error) {
	return ResultNil, err
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (r Result) SessId() session.Id {
	return r.sessId
}

func (r Result) SessExpAt() time.Time {
	return r.sessExpAt
}

func (r Result) IsPassed() bool {
	return !r.sessId.IsNil()
}

func (r Result) AttemptsLeft() uint {
	return r.attemptsLeft
}

// ---------------------------------------------------------------------------------------------------------------------
