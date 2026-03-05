package account

import (
	"github.com/selyukovn/go-std"
	"net/netip"
)

// ---------------------------------------------------------------------------------------------------------------------
// Const
// ---------------------------------------------------------------------------------------------------------------------

var IpWhitelistEmpty = IpWhitelist{}

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

const IpWhitelistMaxLen = 4

type IpWhitelist struct {
	// Массив -- поскольку слайсы не являются сравнимыми, а нам нужно сравнивать с IpWhitelistEmpty
	// IpWhitelistMaxLen = 4 -- в 99% случаев более, чем достаточно, и subnets целиком помещается в VARBINARY(255).
	subnets [IpWhitelistMaxLen]netip.Prefix
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// IpWhitelistFromPrefixes
//
// Срез должен содержать не более IpWhitelistMaxLen уникальных ненулевых адресов подсетей.
//
//   - std.ErrorValidation
func IpWhitelistFromPrefixes(subnets []netip.Prefix) (IpWhitelist, error) {
	uniqNotNil := make([]netip.Prefix, IpWhitelistMaxLen)

	deduplicateMap := make(map[netip.Prefix]struct{}, len(subnets))
	for k, subnet := range subnets {
		if !subnet.IsValid() {
			continue
		}

		if _, ok := deduplicateMap[subnet]; ok {
			continue
		}

		if len(uniqNotNil) == IpWhitelistMaxLen {
			return IpWhitelistEmpty, std.NewErrorValidationFf(
				"cрез должен содержать не более %d уникальных ненулевых значений, а не %v",
				IpWhitelistMaxLen,
				subnets,
			)
		}

		deduplicateMap[subnet] = struct{}{}
		uniqNotNil[k] = subnet
	}

	return IpWhitelist{subnets: ([IpWhitelistMaxLen]netip.Prefix)(uniqNotNil)}, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (wl IpWhitelist) IsEmpty() bool {
	return wl == IpWhitelistEmpty
}

func (wl IpWhitelist) IsIpAllowed(ip netip.Addr) bool {
	if !ip.IsValid() {
		return false
	}

	if wl.IsEmpty() {
		return true
	}

	for _, subnet := range wl.subnets {
		if subnet.Contains(ip) {
			return true
		}
	}

	return false
}

func (wl IpWhitelist) Subnets() []netip.Prefix {
	// Поскольку элементы неизменяемые,
	// можно спокойно вернуть исходные данные -- копия не нужна.
	return wl.subnets[:]
}

// ---------------------------------------------------------------------------------------------------------------------
