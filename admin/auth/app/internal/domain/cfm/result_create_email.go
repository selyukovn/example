package cfm

import (
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Const
// ---------------------------------------------------------------------------------------------------------------------

var ServiceResultCreateNil = ServiceResultCreate{}

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type ServiceResultCreate struct {
	email    std.Email
	cfmId    Id
	expireAt time.Time
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewServiceResultCreate
//
// Паникует при нулевых аргументах.
func NewServiceResultCreate(email std.Email, cfmId Id, expireAt time.Time) ServiceResultCreate {
	assert.Cmp[std.Email]().NotEq(std.EmailNil).Must(email)
	assert.Cmp[Id]().NotEq(IdNil).Must(cfmId)
	assert.Time().NotZero().Must(expireAt)

	return ServiceResultCreate{
		email:    email,
		cfmId:    cfmId,
		expireAt: expireAt,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (r ServiceResultCreate) IsNil() bool {
	return r == ServiceResultCreateNil
}

func (r ServiceResultCreate) Email() std.Email {
	return r.email
}

func (r ServiceResultCreate) CfmId() Id {
	return r.cfmId
}

func (r ServiceResultCreate) ExpireAt() time.Time {
	return r.expireAt
}

// ---------------------------------------------------------------------------------------------------------------------
