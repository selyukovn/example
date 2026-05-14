package handlers

import (
	"context"
	"example/admin/gateway/cmd/kafcon/bundles/admin_auth_events/kafapi"
	"example/admin/gateway/cmd/kafcon/container"
)

func NewSessionClosedV1(ctr *container.Container) func(context.Context, *kafapi.Meta, *kafapi.DataSessionClosedV1) error {
	return func(ctx context.Context, meta *kafapi.Meta, data *kafapi.DataSessionClosedV1) error {
		return ctr.AuthServiceCacher.CheckSessionUnsetBySessionId(ctx, data.SessionId)
	}
}
