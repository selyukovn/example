package db

import (
	"context"
	"database/sql"
	"github.com/huandu/go-sqlbuilder"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// #####################################################################################################################
// ROW
// #####################################################################################################################

type EventRow struct {
	OccurredAt      time.Time
	Type            string
	Version         uint
	ExtraCfmId      sql.NullString
	ExtraFinishType sql.NullByte
	CreatedAt       time.Time
}

// #####################################################################################################################
// TABLE
// #####################################################################################################################

var EventTable = eventTable{}

// ---------------------------------------------------------------------------------------------------------------------

type eventTable struct{}

func (t eventTable) Name() string                  { return "event" }
func (t eventTable) ColumnIdAutoincrement() string { return "id" }
func (t eventTable) ColumnOccurredAt() string      { return "occurred_at" }
func (t eventTable) ColumnType() string            { return "type" }
func (t eventTable) ColumnVersion() string         { return "version" }
func (t eventTable) ColumnExtraCfmId() string      { return "extra_confirmation_id" }
func (t eventTable) ColumnExtraFinishType() string { return "extra_finish_type" }
func (t eventTable) ColumnCreatedAt() string       { return "created_at" }

// ---------------------------------------------------------------------------------------------------------------------

// Insert
//
// Паникует при нулевых аргументах.
func (t eventTable) Insert(ctx context.Context, tx *sql.Tx, row *EventRow) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*sql.Tx]().NotEq(nil).Must(tx)
	assert.Cmp[*EventRow]().NotEq(nil).Must(row)

	_, err := tx.ExecContext(ctx, sqlbuilder.
		InsertInto(t.Name()).
		Cols(
			t.ColumnOccurredAt(),
			t.ColumnType(),
			t.ColumnVersion(),
			t.ColumnExtraCfmId(),
			t.ColumnExtraFinishType(),
			t.ColumnCreatedAt(),
		).
		Values(
			t.ColumnOccurredAt(),
			t.ColumnType(),
			t.ColumnVersion(),
			t.ColumnExtraCfmId(),
			t.ColumnExtraFinishType(),
			t.ColumnCreatedAt(),
		).
		String(),
		// --
		row.OccurredAt,
		row.Type,
		row.Version,
		row.ExtraCfmId,
		row.ExtraFinishType,
		row.CreatedAt,
	)

	return err
}

// ---------------------------------------------------------------------------------------------------------------------
