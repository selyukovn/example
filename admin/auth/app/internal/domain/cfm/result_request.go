package cfm

import (
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Const
// ---------------------------------------------------------------------------------------------------------------------

var ServiceResultRequestNil = ServiceResultRequest{}

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type ServiceResultRequest struct {
	cfmId              Id
	canReqAgain        bool
	canReqAttemptsLeft uint
	canReqAfter        time.Time
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewServiceResultRequestCanAgain
//
// Паникует при нулевых аргументах.
func NewServiceResultRequestCanAgain(cfmId Id, canReqAttemptsLeft uint, canReqAfter time.Time) ServiceResultRequest {
	assert.FalseMust(cfmId.IsNil())
	assert.FalseMust(canReqAttemptsLeft == 0)
	assert.FalseMust(canReqAfter.IsZero())

	return ServiceResultRequest{
		cfmId:              cfmId,
		canReqAgain:        true,
		canReqAttemptsLeft: canReqAttemptsLeft,
		canReqAfter:        canReqAfter,
	}
}

// NewServiceResultRequestLast
//
// Паникует при нулевых аргументах.
func NewServiceResultRequestLast(cfmId Id) ServiceResultRequest {
	assert.FalseMust(cfmId.IsNil())

	return ServiceResultRequest{
		cfmId:              cfmId,
		canReqAgain:        false,
		canReqAttemptsLeft: 0,
		canReqAfter:        time.Time{},
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (r ServiceResultRequest) IsNil() bool {
	return r == ServiceResultRequestNil
}

func (r ServiceResultRequest) CfmId() Id {
	return r.cfmId
}

func (r ServiceResultRequest) CanReqAgain() bool {
	return r.canReqAgain
}

// CanReqAttemptsLeft
//
// Возвращает нулевое значение, если не CanReqAgain()
func (r ServiceResultRequest) CanReqAttemptsLeft() uint {
	return r.canReqAttemptsLeft
}

// CanReqAfter
//
// Возвращает нулевое значение, если не CanReqAgain()
func (r ServiceResultRequest) CanReqAfter() time.Time {
	return r.canReqAfter
}

// ---------------------------------------------------------------------------------------------------------------------
