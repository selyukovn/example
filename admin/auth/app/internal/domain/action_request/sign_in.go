package action_request

import (
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/cfm"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type SignIn struct {
	id          Id
	accId       account.Id
	cfmId       cfm.Id
	requestedAt time.Time
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// См. Factory.CreateSignIn
var _ = Factory{}

// ---------------------------------------------------------------------------------------------------------------------

func ReflectExtract(s *SignIn) (
	id Id,
	accId account.Id,
	cfmId cfm.Id,
	requestedAt time.Time,
) {
	return s.id,
		s.accId,
		s.cfmId,
		s.requestedAt
}

func ReflectRestore(
	id Id,
	accId account.Id,
	cfmId cfm.Id,
	requestedAt time.Time,
) *SignIn {
	return &SignIn{
		id:          id,
		accId:       accId,
		cfmId:       cfmId,
		requestedAt: requestedAt,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (s *SignIn) Id() Id {
	return s.id
}

func (s *SignIn) AccId() account.Id {
	return s.accId
}

func (s *SignIn) CfmId() cfm.Id {
	return s.cfmId
}

func (s *SignIn) RequestedAt() time.Time {
	return s.requestedAt
}

// ---------------------------------------------------------------------------------------------------------------------
