package session

import (
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/action_request"
	"example/admin/auth/internal/domain/client"
	"github.com/selyukovn/go-events"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Session struct {
	id            Id
	accId         account.Id
	signInId      action_request.Id
	initialClient client.Client
	initiatedAt   time.Time
	expireAt      time.Time
	closedAt      time.Time
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// См. Factory.Create
var _ = Factory{}

// ---------------------------------------------------------------------------------------------------------------------

func ReflectExtract(s *Session) (
	id Id,
	accId account.Id,
	signInId action_request.Id,
	initialClient client.Client,
	initiatedAt time.Time,
	expireAt time.Time,
	closedAt time.Time,
) {
	return s.id,
		s.accId,
		s.signInId,
		s.initialClient,
		s.initiatedAt,
		s.expireAt,
		s.closedAt
}

func ReflectRestore(
	id Id,
	accId account.Id,
	signInId action_request.Id,
	initialClient client.Client,
	initiatedAt time.Time,
	expireAt time.Time,
	closedAt time.Time,
) *Session {
	return &Session{
		id:            id,
		accId:         accId,
		signInId:      signInId,
		initialClient: initialClient,
		initiatedAt:   initiatedAt,
		expireAt:      expireAt,
		closedAt:      closedAt,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (s *Session) close(now time.Time, evs *event.Collection) {
	s.closedAt = now
	evs.Add(NewEventClosed(now, s.id, s.accId))
}

// TickTime
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - ErrorClosed -- закрытые сессии не обновляются
//   - std.ErrorAlreadyDone -- когда нечего менять на данный момент
func (s *Session) TickTime(now time.Time, evs *event.Collection) error {
	assert.Time().NotZero().Must(now)
	assert.Cmp[*event.Collection]().NotEq(nil).Must(evs)

	if s.isClosedAsGoingToExpire(now) {
		s.close(now, evs)
		return nil
	} else if s.IsClosed(now) {
		return NewErrorSessionClosed(s.id, s.closedAt, s.expireAt)
	}

	return std.NewErrorAlreadyDoneFf("Нечего менять для сессии %q на момент времени %s", s.id, now)
}

// Close
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - ErrorClosed
func (s *Session) Close(now time.Time, evs *event.Collection) error {
	assert.Time().NotZero().Must(now)
	assert.Cmp[*event.Collection]().NotEq(nil).Must(evs)

	if s.IsClosed(now) {
		return NewErrorSessionClosed(s.id, s.closedAt, s.expireAt)
	}

	s.close(now, evs)

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (s *Session) Id() Id {
	return s.id
}

func (s *Session) AccId() account.Id {
	return s.accId
}

func (s *Session) SignInId() action_request.Id {
	return s.signInId
}

func (s *Session) ExpireAt() time.Time {
	return s.expireAt
}

// isClosedAsGoingToExpire
//
// Нельзя просто так взять и Close точно в момент expireAt
// (технически возможно, например, по таймеру в отдельной горутине, но это гораздо накладнее, чем "кроном"),
// поэтому между моментом expireAt и фактическим closedAt может пройти какое-то время.
// Это нужно учитывать также в поисковых запросах вроде RepositoryInterface.GetIdsOfGoingToExpire
func (s *Session) isClosedAsGoingToExpire(now time.Time) bool {
	return s.closedAt.IsZero() && now.After(s.expireAt)
}

func (s *Session) IsClosed(now time.Time) bool {
	return !s.closedAt.IsZero() || s.isClosedAsGoingToExpire(now)
}

// ---------------------------------------------------------------------------------------------------------------------
