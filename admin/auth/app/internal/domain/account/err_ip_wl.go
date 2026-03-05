package account

import (
	"fmt"
	assert "github.com/selyukovn/go-wm-assert"
	"net/netip"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type ErrorIpWhitelist struct {
	accId    Id
	missedIp netip.Addr
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewErrorIpWhitelist
//
// Паникует при нулевых аргументах.
func NewErrorIpWhitelist(accId Id, missedIp netip.Addr) ErrorIpWhitelist {
	assert.Cmp[Id]().NotEq(IdNil).Must(accId)
	assert.Cmp[netip.Addr]().NotEq(netip.Addr{}).Must(missedIp)

	return ErrorIpWhitelist{
		accId:    accId,
		missedIp: missedIp,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (e ErrorIpWhitelist) Error() string {
	return fmt.Sprintf(
		"IP-адрес %q не входит в список разрешенных IP-адресов аккаунта %q",
		e.missedIp.String(),
		e.accId.String(),
	)
}

func (e ErrorIpWhitelist) AccId() Id {
	return e.accId
}

func (e ErrorIpWhitelist) MissedIp() netip.Addr {
	return e.missedIp
}

// ---------------------------------------------------------------------------------------------------------------------
