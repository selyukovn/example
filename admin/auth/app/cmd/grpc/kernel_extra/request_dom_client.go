package kernel_extra

import (
	"example/admin/auth/internal/domain/client"
	"github.com/selyukovn/go-std"
	"net/netip"
)

// ParseClient
//
//   - std.ErrorValidation
func ParseClient(fromIp string, fromUserAgent string) (client.Client, error) {
	ua, err := client.UserAgentFromString(fromUserAgent)
	if err != nil {
		return client.ClientNil, err
	}

	ip, err := netip.ParseAddr(fromIp)
	if err != nil {
		return client.ClientNil, std.NewErrorValidationFf(err.Error())
	}

	cl, err := client.NewClient(ua, ip)
	if err != nil {
		return client.ClientNil, err
	}

	return cl, nil
}
