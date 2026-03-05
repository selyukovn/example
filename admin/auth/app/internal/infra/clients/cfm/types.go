package cfm

import (
	"fmt"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// CREATE FOR EMAIL
// ---------------------------------------------------------------------------------------------------------------------

type CreateForEmailResult = struct {
	CfmId    string
	ExpireAt time.Time
}

// ---------------------------------------------------------------------------------------------------------------------
// REQUEST
// ---------------------------------------------------------------------------------------------------------------------

type RequestResult = struct {
	CfmId              string
	CanReqAgain        bool
	CanReqAttemptsLeft uint
	CanReqAfter        time.Time
}

// ---------------------------------------------------------------------------------------------------------------------
// CONFIRM
// ---------------------------------------------------------------------------------------------------------------------

type ConfirmResult = struct {
	CfmId              string
	FinishedAt         time.Time
	IsFinishedAsPassed bool
	IsFinishedAsFailed bool
	FailsLeft          uint
}

// ---------------------------------------------------------------------------------------------------------------------
// ERRORS
// ---------------------------------------------------------------------------------------------------------------------

type ErrorFinished struct {
	CfmId       string
	FinishedAt  time.Time
	IsAsPassed  bool
	IsAsFailed  bool
	IsAsExpired bool
}

func (e ErrorFinished) Error() string {
	return fmt.Sprintf("%#v", e)
}

// ---------------------------------------------------------------------------------------------------------------------

type ErrorNoAttemptsLeft struct {
	CfmId string
}

func (e ErrorNoAttemptsLeft) Error() string {
	return fmt.Sprintf("%#v", e)
}

// ---------------------------------------------------------------------------------------------------------------------

type ErrorRequestsFrequency struct {
	CfmId              string
	CanReqAfter        time.Time
	CanReqAttemptsLeft uint
}

func (e ErrorRequestsFrequency) Error() string {
	return fmt.Sprintf("%#v", e)
}

// ---------------------------------------------------------------------------------------------------------------------
