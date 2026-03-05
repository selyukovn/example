package event_storage

import (
	"context"
	"example/admin/cfm/internal/domain/cfm"
)

type RepositoryInterface interface {
	// AddCfmFinished
	//
	// Паникует при нулевых аргументах.
	//
	// Ошибки:
	//   - std.ErrorRuntime
	AddCfmFinished(ctx context.Context, e cfm.EventFinished) error
}
