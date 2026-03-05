package auth

import (
	"fmt"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// RESULTS
// ---------------------------------------------------------------------------------------------------------------------

type SignInRequestResult struct {
	SignInId    string
	RetriesLeft int
	CanRetryAt  time.Time
	ExpireAt    time.Time
}

type SignInRequestRetryResult struct {
	SignInId    string
	RetriesLeft int
	CanRetryAt  time.Time
}

type SignInConfirmResult struct {
	IsPassed        bool
	AttemptsLeft    int
	SessionId       string
	SessionExpireAt time.Time
}

type CheckSessionResult struct {
	AccountId       string
	SessionExpireAt time.Time
}

// ---------------------------------------------------------------------------------------------------------------------
// ERRORS
// ---------------------------------------------------------------------------------------------------------------------

type ErrorValidation struct {
	Field   string
	Message string
}

func (e ErrorValidation) Error() string {
	return fmt.Sprintf("%#v", e)
}

// ---------------------------------------------------------------------------------------------------------------------

type ErrorAccountAccessDenied struct{}

func (e ErrorAccountAccessDenied) Error() string {
	return fmt.Sprintf("%#v", e)
}

// ---------------------------------------------------------------------------------------------------------------------

type ErrorSignInFinished struct {
	IsPassed  bool
	IsFailed  bool
	IsExpired bool
}

func (e ErrorSignInFinished) Error() string {
	return fmt.Sprintf("%#v", e)
}

// ---------------------------------------------------------------------------------------------------------------------

type ErrorNoAttemptsLeft struct{}

func (e ErrorNoAttemptsLeft) Error() string {
	return fmt.Sprintf("%#v", e)
}

// ---------------------------------------------------------------------------------------------------------------------

type ErrorRequestsFrequency struct {
	CanReqAfter        time.Time
	CanReqAttemptsLeft int
}

func (e ErrorRequestsFrequency) Error() string {
	return fmt.Sprintf("%#v", e)
}

// ---------------------------------------------------------------------------------------------------------------------

type ErrorSessionClosed struct {
	IsExpired bool
	ClosedAt  time.Time
}

func (e ErrorSessionClosed) Error() string {
	return fmt.Sprintf("%#v", e)
}

// ---------------------------------------------------------------------------------------------------------------------
