package actions

import (
	"context"
	adapt_infra_clients_auth_cachable "example/admin/gateway/internal/adapt/infra/clients/auth/cachable"
	"example/admin/gateway/internal/api/kafcon/handlers/admin_auth_events/kafapi"
)

// ---------------------------------------------------------------------------------------------------------------------

type SessionClosedV1 = func(context.Context, *kafapi.Meta, *kafapi.DataSessionClosedV1) error

// ---------------------------------------------------------------------------------------------------------------------

func NewSessionClosedV1(sAuthCacher adapt_infra_clients_auth_cachable.Cacher) SessionClosedV1 {
	return func(ctx context.Context, meta *kafapi.Meta, data *kafapi.DataSessionClosedV1) error {
		return sAuthCacher.CheckSessionUnsetBySessionId(ctx, data.SessionId)
	}
}

// ---------------------------------------------------------------------------------------------------------------------
