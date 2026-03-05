package create_for_email

import (
	"example/admin/cfm/internal/domain/cfm"
	"github.com/selyukovn/go-std"
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
	cfmId    cfm.Id
	email    std.Email
	expireAt time.Time
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// resultSuccess
//
// Паникует при нулевых аргументах.
func resultSuccess(cfmId cfm.Id, email std.Email, expireAt time.Time) (Result, error) {
	assert.FalseMust(cfmId.IsNil())
	assert.FalseMust(email.IsNil())
	assert.FalseMust(expireAt.IsZero())

	return Result{
		cfmId:    cfmId,
		email:    email,
		expireAt: expireAt,
	}, nil
}

// resultError
//
// Паникует при нулевых аргументах.
func resultError(err error) (Result, error) {
	assert.NotNilDeepMust(err)

	return Result{cfm.IdNil, std.EmailNil, time.Time{}}, err
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (r Result) IsNil() bool {
	return r == ResultNil
}

func (r Result) CfmId() cfm.Id {
	return r.cfmId
}

func (r Result) Email() std.Email {
	return r.email
}

func (r Result) ExpireAt() time.Time {
	return r.expireAt
}

// ---------------------------------------------------------------------------------------------------------------------
