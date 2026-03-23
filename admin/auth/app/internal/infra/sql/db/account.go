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

type AccountRow = struct {
	Id              string
	Email           string
	IsActive        bool
	DeactivatedAt   sql.NullTime
	IpWhitelistJson string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// #####################################################################################################################
// TABLE
// #####################################################################################################################

var AccountTable = accountTable{}

// ---------------------------------------------------------------------------------------------------------------------

type accountTable struct{}

func (t accountTable) Name() string                  { return "account" }
func (t accountTable) ColumnId() string              { return "id" }
func (t accountTable) ColumnEmail() string           { return "email" }
func (t accountTable) ColumnIsActive() string        { return "is_active" }
func (t accountTable) ColumnDeactivatedAt() string   { return "deactivated_at" }
func (t accountTable) ColumnIpWhitelistJson() string { return "ip_whitelist_json" }
func (t accountTable) ColumnCreatedAt() string       { return "created_at" }
func (t accountTable) ColumnUpdatedAt() string       { return "updated_at" }

// ---------------------------------------------------------------------------------------------------------------------

// Insert
//
// Паникует при нулевых аргументах.
func (t accountTable) Insert(ctx context.Context, tx *sql.Tx, row *AccountRow) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*sql.Tx]().NotEq(nil).Must(tx)
	assert.Cmp[*AccountRow]().NotEq(nil).Must(row)

	_, err := tx.ExecContext(ctx, sqlbuilder.
		InsertInto(t.Name()).
		Cols(
			t.ColumnId(),
			t.ColumnEmail(),
			t.ColumnIsActive(),
			t.ColumnDeactivatedAt(),
			t.ColumnIpWhitelistJson(),
			t.ColumnCreatedAt(),
			t.ColumnUpdatedAt(),
		).
		Values(
			row.Id,
			row.Email,
			row.IsActive,
			row.DeactivatedAt,
			row.IpWhitelistJson,
			row.CreatedAt,
			row.UpdatedAt,
		).
		String(),
		// --
		row.Id,
		row.Email,
		row.IsActive,
		row.DeactivatedAt,
		row.IpWhitelistJson,
		row.CreatedAt,
		row.UpdatedAt,
	)

	return err
}

// Update
//
// Паникует при нулевых аргументах.
func (t accountTable) Update(ctx context.Context, tx *sql.Tx, row *AccountRow) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*sql.Tx]().NotEq(nil).Must(tx)
	assert.Cmp[*AccountRow]().NotEq(nil).Must(row)

	_, err := tx.ExecContext(ctx, sqlbuilder.
		Update(t.Name()).
		Set(
			fmt.Sprintf("%s = ?", t.ColumnEmail()),
			fmt.Sprintf("%s = ?", t.ColumnIsActive()),
			fmt.Sprintf("%s = ?", t.ColumnDeactivatedAt()),
			fmt.Sprintf("%s = ?", t.ColumnIpWhitelistJson()),
			fmt.Sprintf("%s = ?", t.ColumnUpdatedAt()),
		).
		Where(fmt.Sprintf("%s = ?", t.ColumnId())).
		String(),
		// --
		row.Email,
		row.IsActive,
		row.DeactivatedAt,
		row.IpWhitelistJson,
		row.UpdatedAt,
		// --
		row.Id,
	)

	return err
}

// ---------------------------------------------------------------------------------------------------------------------

// QueryWhereId
//
// Паникует при нулевых аргументах:
//   - ctx
//   - tx
func (t accountTable) QueryWhereId(ctx context.Context, tx *sql.Tx, id string) (*AccountRow, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*sql.Tx]().NotEq(nil).Must(tx)

	row := &AccountRow{}
	err := tx.
		QueryRowContext(ctx, sqlbuilder.
			Select(
				t.ColumnId(),
				t.ColumnEmail(),
				t.ColumnIsActive(),
				t.ColumnDeactivatedAt(),
				t.ColumnIpWhitelistJson(),
			).
			From(t.Name()).
			Where(fmt.Sprintf("%s = ?", t.ColumnId())).
			String(),
			// --
			id,
		).
		Scan(
			&row.Id,
			&row.Email,
			&row.IsActive,
			&row.DeactivatedAt,
			&row.IpWhitelistJson,
		)

	return row, err
}

// QueryWhereEmail
//
// Паникует при нулевых аргументах:
//   - ctx
//   - tx
func (t accountTable) QueryWhereEmail(ctx context.Context, tx *sql.Tx, email string) (*AccountRow, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*sql.Tx]().NotEq(nil).Must(tx)

	row := &AccountRow{}
	err := tx.
		QueryRowContext(ctx, sqlbuilder.
			Select(
				t.ColumnId(),
				t.ColumnEmail(),
				t.ColumnIsActive(),
				t.ColumnDeactivatedAt(),
				t.ColumnIpWhitelistJson(),
			).
			From(t.Name()).
			Where(fmt.Sprintf("%s = ?", t.ColumnEmail())).
			String(),
			// --
			email,
		).
		Scan(
			&row.Id,
			&row.Email,
			&row.IsActive,
			&row.DeactivatedAt,
			&row.IpWhitelistJson,
		)

	return row, err
}

// QueryEmailWhereId
//
// Паникует при нулевых аргументах:
//   - ctx
//   - tx
func (t accountTable) QueryEmailWhereId(ctx context.Context, tx *sql.Tx, id string) (string, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*sql.Tx]().NotEq(nil).Must(tx)

	var str string
	err := tx.
		QueryRowContext(ctx, sqlbuilder.
			Select(t.ColumnEmail()).
			From(t.Name()).
			Where(fmt.Sprintf("%s = ?", t.ColumnId())).
			String(),
			// --
			id,
		).
		Scan(&str)

	return str, err
}

// ---------------------------------------------------------------------------------------------------------------------
