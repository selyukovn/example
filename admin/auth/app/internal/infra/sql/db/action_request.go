package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// #####################################################################################################################
// ROW
// #####################################################################################################################

type ActionRequestRow struct {
	Id          string
	Type        uint
	AccId       string
	CfmId       string
	RequestedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// #####################################################################################################################
// TABLE
// #####################################################################################################################

var ActionRequestTable = actionRequestTable{}

// ---------------------------------------------------------------------------------------------------------------------

type actionRequestTable struct{}

func (t actionRequestTable) Name() string              { return "action_request" }
func (t actionRequestTable) ColumnId() string          { return "id" }
func (t actionRequestTable) ColumnType() string        { return "type" }
func (t actionRequestTable) ColumnAccId() string       { return "account_id" }
func (t actionRequestTable) ColumnCfmId() string       { return "confirmation_id" }
func (t actionRequestTable) ColumnRequestedAt() string { return "requested_at" }
func (t actionRequestTable) ColumnCreatedAt() string   { return "created_at" }
func (t actionRequestTable) ColumnUpdatedAt() string   { return "updated_at" }

// ---------------------------------------------------------------------------------------------------------------------

// Insert
//
// Паникует при нулевых аргументах.
func (t actionRequestTable) Insert(ctx context.Context, tx *sql.Tx, row *ActionRequestRow) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*sql.Tx]().NotEq(nil).Must(tx)
	assert.Cmp[*ActionRequestRow]().NotEq(nil).Must(row)

	_, err := tx.ExecContext(ctx, sqlbuilder.
		InsertInto(t.Name()).
		Cols(
			t.ColumnId(),
			t.ColumnType(),
			t.ColumnAccId(),
			t.ColumnCfmId(),
			t.ColumnRequestedAt(),
			t.ColumnCreatedAt(),
			t.ColumnUpdatedAt(),
		).
		Values(
			row.Id,
			row.Type,
			row.AccId,
			row.CfmId,
			row.RequestedAt,
			row.CreatedAt,
			row.UpdatedAt,
		).
		String(),
		// --
		row.Id,
		row.Type,
		row.AccId,
		row.CfmId,
		row.RequestedAt,
		row.CreatedAt,
		row.UpdatedAt,
	)

	return err
}

// Update
//
// Паникует при нулевых аргументах.
func (t actionRequestTable) Update(ctx context.Context, tx *sql.Tx, row *ActionRequestRow) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*sql.Tx]().NotEq(nil).Must(tx)
	assert.Cmp[*ActionRequestRow]().NotEq(nil).Must(row)

	_, err := tx.ExecContext(ctx, sqlbuilder.
		Update(t.Name()).
		Set(
			fmt.Sprintf("%s = ?", t.ColumnType()),
			fmt.Sprintf("%s = ?", t.ColumnAccId()),
			fmt.Sprintf("%s = ?", t.ColumnCfmId()),
			fmt.Sprintf("%s = ?", t.ColumnRequestedAt()),
			fmt.Sprintf("%s = ?", t.ColumnUpdatedAt()),
		).
		Where(fmt.Sprintf("%s = ?", t.ColumnId())).
		String(),
		// --
		row.Type,
		row.AccId,
		row.CfmId,
		row.RequestedAt,
		row.UpdatedAt,
		// --
		row.Id,
	)

	return err
}

// ---------------------------------------------------------------------------------------------------------------------

// QueryWhereIdAndType
//
// Паникует при нулевых аргументах:
//   - ctx
//   - tx
func (t actionRequestTable) QueryWhereIdAndType(
	ctx context.Context,
	tx *sql.Tx,
	id string,
	tType uint,
) (*ActionRequestRow, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*sql.Tx]().NotEq(nil).Must(tx)

	row := &ActionRequestRow{}
	err := tx.
		QueryRowContext(ctx, sqlbuilder.
			Select(
				t.ColumnId(),
				t.ColumnType(),
				t.ColumnAccId(),
				t.ColumnCfmId(),
				t.ColumnRequestedAt(),
			).
			From(t.Name()).
			Where(fmt.Sprintf("%s = ?", t.ColumnId())).
			Where(fmt.Sprintf("%s = ?", t.ColumnType())).
			String(),
			// --
			id,
			tType,
		).
		Scan(
			&row.Id,
			&row.Type,
			&row.AccId,
			&row.CfmId,
			&row.RequestedAt,
		)

	return row, err
}

// ---------------------------------------------------------------------------------------------------------------------
