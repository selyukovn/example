package client

import (
	"github.com/selyukovn/go-std"
	"net/netip"
)

// ---------------------------------------------------------------------------------------------------------------------
// Const
// ---------------------------------------------------------------------------------------------------------------------

var ClientNil = Client{}

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Client struct {
	userAgent UserAgent
	ipAddress netip.Addr
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewClient
//
//   - std.ErrorValidation
func NewClient(userAgent UserAgent, ipAddress netip.Addr) (Client, error) {
	if userAgent.IsNil() {
		return ClientNil, std.NewErrorValidationFf("User-Agent не может быть пустым")
	}

	if !ipAddress.IsValid() {
		return ClientNil, std.NewErrorValidationFf("IP не может быть пустым")
	}

	c := Client{
		userAgent: userAgent,
		ipAddress: ipAddress,
	}

	return c, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (c Client) IsNil() bool {
	return c == ClientNil
}

func (c Client) UserAgent() UserAgent {
	return c.userAgent
}

func (c Client) IpAddress() netip.Addr {
	return c.ipAddress
}

// ---------------------------------------------------------------------------------------------------------------------
