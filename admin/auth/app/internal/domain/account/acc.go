package account

import (
	"example/admin/auth/internal/domain/client"
	"github.com/selyukovn/go-events"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Account struct {
	id            Id
	email         std.Email
	deactivatedAt time.Time
	ipWhitelist   IpWhitelist
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// См. Factory.Create
var _ = Factory{}

// ---------------------------------------------------------------------------------------------------------------------
// Reflect
// ---------------------------------------------------------------------------------------------------------------------

func ReflectExtract(a *Account) (
	id Id,
	email std.Email,
	deactivatedAt time.Time,
	ipWhitelist IpWhitelist,
) {
	return a.id,
		a.email,
		a.deactivatedAt,
		a.ipWhitelist
}

func ReflectRestore(
	id Id,
	email std.Email,
	deactivatedAt time.Time,
	ipWhitelist IpWhitelist,
) *Account {
	return &Account{
		id:            id,
		email:         email,
		deactivatedAt: deactivatedAt,
		ipWhitelist:   ipWhitelist,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Deactivate
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorAlreadyDone
func (a *Account) Deactivate(now time.Time, evs *event.Collection) error {
	assert.Time().NotZero().Must(now)
	assert.Cmp[*event.Collection]().NotEq(nil).Must(evs)

	if a.IsDeactivated() {
		return std.NewErrorAlreadyDoneFf("Аккаунт %q уже деактивирован", a.id)
	}

	a.deactivatedAt = now

	evs.Add(NewEventDeactivated(now, a.id))

	return nil
}

// ChangeIpWhitelist
//
// Паникует при нулевых аргументах:
//   - now
//   - evs
//
// Ошибки:
//   - std.ErrorAlreadyDone
func (a *Account) ChangeIpWhitelist(newList IpWhitelist, now time.Time, evs *event.Collection) error {
	assert.Time().NotZero().Must(now)
	assert.Cmp[*event.Collection]().NotEq(nil).Must(evs)

	if a.ipWhitelist == newList {
		return std.NewErrorAlreadyDoneFf("IpWhitelist для аккаунта %q уже таков %v", a.id, newList)
	}

	a.ipWhitelist = newList

	evs.Add(NewEventIpWhitelistChanged(now, a.id, newList))

	return nil
}

// AssertSignIn
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - ErrorDeactivated
//   - ErrorIpWhitelist
func (a *Account) AssertSignIn(cl client.Client, now time.Time) error {
	assert.Cmp[client.Client]().NotEq(client.ClientNil).Must(cl)
	assert.Time().NotZero().Must(now)

	if a.IsDeactivated() {
		return NewErrorDeactivated(a.id)
	}

	if !a.ipWhitelist.IsIpAllowed(cl.IpAddress()) {
		return NewErrorIpWhitelist(a.id, cl.IpAddress())
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (a *Account) Id() Id {
	return a.id
}

func (a *Account) Email() std.Email {
	return a.email
}

func (a *Account) IsDeactivated() bool {
	return !a.deactivatedAt.IsZero()
}

// ---------------------------------------------------------------------------------------------------------------------
