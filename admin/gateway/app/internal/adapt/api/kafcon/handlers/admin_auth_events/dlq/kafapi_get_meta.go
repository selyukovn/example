package dlq

import (
	"context"
	"example/admin/gateway/internal/api/kafcon/handlers/admin_auth_events/kafapi"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ kafapi.ServiceInterface = new(sDlqGetMeta)

type sDlqGetMeta struct {
	Meta *kafapi.Meta
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (s *sDlqGetMeta) handle(meta *kafapi.Meta) error {
	s.Meta = meta
	return nil
}

func (s *sDlqGetMeta) AccountCreatedV1(_ context.Context, meta *kafapi.Meta, _ *kafapi.DataAccountCreatedV1) error {
	return s.handle(meta)
}
func (s *sDlqGetMeta) AccountDeactivatedV1(_ context.Context, meta *kafapi.Meta, _ *kafapi.DataAccountDeactivatedV1) error {
	return s.handle(meta)
}
func (s *sDlqGetMeta) IpWhitelistChangedV1(_ context.Context, meta *kafapi.Meta, _ *kafapi.DataIpWhitelistChangedV1) error {
	return s.handle(meta)
}
func (s *sDlqGetMeta) SessionCreatedV1(_ context.Context, meta *kafapi.Meta, _ *kafapi.DataSessionCreatedV1) error {
	return s.handle(meta)
}
func (s *sDlqGetMeta) SessionClosedV1(_ context.Context, meta *kafapi.Meta, _ *kafapi.DataSessionClosedV1) error {
	return s.handle(meta)
}

// ---------------------------------------------------------------------------------------------------------------------
