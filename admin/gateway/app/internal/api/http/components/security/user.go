package security

import (
	assert "github.com/selyukovn/go-wm-assert"
	"net/netip"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type User struct {
	ip        netip.Addr
	userAgent string
	sessId    string
	sessExpAt time.Time
	accId     string
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func newUserGuest(ip netip.Addr, userAgent string) *User {
	assert.NotZeroMust(ip)
	assert.Str().NotEmpty(userAgent)

	return &User{
		ip:        ip,
		userAgent: userAgent,
	}
}

func newUserAuthorized(
	ip netip.Addr,
	userAgent string,
	sessId string,
	sessExpAt time.Time,
	accId string,
) *User {
	assert.NotZeroMust(ip)
	assert.Str().NotEmpty(userAgent)
	assert.Str().NotEmpty().Must(sessId)
	assert.Time().NotZero().Must(sessExpAt)
	assert.Str().NotEmpty().Must(accId)

	return &User{
		ip:        ip,
		userAgent: userAgent,
		sessId:    sessId,
		sessExpAt: sessExpAt,
		accId:     accId,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (u *User) authenticate(sessId string, sessExpAt time.Time, accId string) {
	assert.Str().NotEmpty().Must(sessId)
	assert.Time().NotZero().Must(sessExpAt)
	assert.Str().NotEmpty().Must(accId)

	assert.TrueMust(u.IsGuest())

	u.sessId = sessId
	u.sessExpAt = sessExpAt
	u.accId = accId
}

func (u *User) unAuthenticate() {
	assert.TrueMust(u.IsAuthenticated())

	u.sessId = ""
	u.sessExpAt = time.Time{}
	u.accId = ""
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (u *User) sessionId() string {
	return u.sessId
}

// ---------------------------------------------------------------------------------------------------------------------

func (u *User) Ip() netip.Addr {
	return u.ip
}

func (u *User) UserAgent() string {
	return u.userAgent
}

func (u *User) IsGuest() bool {
	return u.accId == ""
}

func (u *User) IsAuthenticated() bool {
	return !u.IsGuest()
}

// AccountId
//
// Вернет пустую строку, если IsGuest() == true.
func (u *User) AccountId() string {
	return u.accId
}

// ---------------------------------------------------------------------------------------------------------------------
