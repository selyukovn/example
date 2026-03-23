package event_storage

import (
	"context"
	"database/sql"
	"example/admin/cfm/internal/domain/cfm"
	"example/admin/cfm/internal/domain/event_storage"
	"example/admin/cfm/internal/infra/sql/db"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-txr"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ event_storage.RepositoryInterface = RepositoryImplSql{}

type RepositoryImplSql struct{}

const (
	EventTypeCfmFinished = "cfm_finished"
)

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewRepositoryImplSql() RepositoryImplSql {
	return RepositoryImplSql{}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// AddCfmFinished
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (r RepositoryImplSql) AddCfmFinished(ctx context.Context, e cfm.EventFinished) error {
	assert.NotNilDeepMust(ctx)
	assert.Cmp[cfm.EventFinished]().NotEq(cfm.EventFinished{}).Must(e)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	dbRow := &db.EventRow{}
	dbRow.CreatedAt = time.Now()
	dbRow.Type = EventTypeCfmFinished
	dbRow.Version = 1
	dbRow.OccurredAt = e.OccurredAt()
	// --
	dbRow.ExtraCfmId = sql.NullString{String: e.CfmId().String(), Valid: true}
	dbRow.ExtraFinishType = sql.NullByte{Byte: byte(e.FinishType()), Valid: true}

	err := db.EventTable.Insert(ctx, tx, dbRow)
	if err != nil {
		return std.WrapErrorToRuntime(err, r, "AddCfmFinished")
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
