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

type EventRow = struct {
	OccurredAt              time.Time
	Type                    string
	Version                 uint
	ExtraAccEmail           sql.NullString
	ExtraAccId              sql.NullString
	ExtraAccIpWhitelistJson sql.NullString
	ExtraSessId             sql.NullString
	CreatedAt               time.Time
	OutboxGroupId           string
	OutboxOperationId       string
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
func (t eventTable) ColumnExtraAccEmail() string   { return "extra_account_email" }
func (t eventTable) ColumnExtraAccId() string      { return "extra_account_id" }
func (t eventTable) ColumnExtraAccIpWhitelistJson() string {
	return "extra_account_ip_whitelist_json"
}
func (t eventTable) ColumnExtraSessId() string       { return "extra_session_id" }
func (t eventTable) ColumnCreatedAt() string         { return "created_at" }
func (t eventTable) ColumnOutboxGroupId() string     { return "outbox_group_id" }
func (t eventTable) ColumnOutboxOperationId() string { return "outbox_operation_id" }

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
			t.ColumnExtraAccEmail(),
			t.ColumnExtraAccId(),
			t.ColumnExtraAccIpWhitelistJson(),
			t.ColumnExtraSessId(),
			t.ColumnCreatedAt(),
			t.ColumnOutboxGroupId(),
			t.ColumnOutboxOperationId(),
		).
		Values(
			t.ColumnOccurredAt(),
			t.ColumnType(),
			t.ColumnVersion(),
			t.ColumnExtraAccEmail(),
			t.ColumnExtraAccId(),
			t.ColumnExtraAccIpWhitelistJson(),
			t.ColumnExtraSessId(),
			t.ColumnCreatedAt(),
			t.ColumnOutboxGroupId(),
			t.ColumnOutboxOperationId(),
		).
		String(),
		// --
		row.OccurredAt,
		row.Type,
		row.Version,
		row.ExtraAccEmail,
		row.ExtraAccId,
		row.ExtraAccIpWhitelistJson,
		row.ExtraSessId,
		row.CreatedAt,
		row.OutboxGroupId,
		row.OutboxOperationId,
	)

	return err
}

// ---------------------------------------------------------------------------------------------------------------------
