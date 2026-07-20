package admin_auth_events

import (
	"context"
	"example/admin/gateway/internal/api/kafcon/handlers/admin_auth_events/actions"
	"example/admin/gateway/internal/api/kafcon/handlers/admin_auth_events/kafapi"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ kafapi.ServiceInterface = sDefault{}

type sDefault struct {
	kafapi.EmptyService
	sessionClosedV1 actions.SessionClosedV1
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func newServiceDefault(
	sessionClosedV1 actions.SessionClosedV1,
) sDefault {
	assert.NotNilDeepMust(sessionClosedV1)

	return sDefault{
		sessionClosedV1: sessionClosedV1,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (s sDefault) SessionClosedV1(ctx context.Context, meta *kafapi.Meta, msg *kafapi.DataSessionClosedV1) error {
	return s.sessionClosedV1(ctx, meta, msg)
}

// ---------------------------------------------------------------------------------------------------------------------
