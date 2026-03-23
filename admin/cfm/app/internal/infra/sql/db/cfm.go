package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"time"
)

// #####################################################################################################################
// ROW
// #####################################################################################################################

type CfmRow = struct {
	Id         string
	Email      string
	ExpireAt   time.Time
	FinishedAt sql.NullTime
	FinishType uint
	FailsMade  uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// #####################################################################################################################
// TABLE
// #####################################################################################################################

var CfmTable = cfmTable{}

type cfmTable struct{}

func (t cfmTable) Name() string             { return "confirmation" }
func (t cfmTable) ColumnId() string         { return "id" }
func (t cfmTable) ColumnEmail() string      { return "email" }
func (t cfmTable) ColumnExpireAt() string   { return "expire_at" }
func (t cfmTable) ColumnFinishedAt() string { return "finished_at" }
func (t cfmTable) ColumnFinishType() string { return "finish_type" }
func (t cfmTable) ColumnFailsMade() string  { return "fails_made" }
func (t cfmTable) ColumnCreatedAt() string  { return "created_at" }
func (t cfmTable) ColumnUpdatedAt() string  { return "updated_at" }

// ---------------------------------------------------------------------------------------------------------------------

func (t cfmTable) Insert(ctx context.Context, tx *sql.Tx, row *CfmRow) error {
	_, err := tx.ExecContext(ctx, sqlbuilder.
		InsertInto(t.Name()).
		Cols(
			t.ColumnId(),
			t.ColumnEmail(),
			t.ColumnExpireAt(),
			t.ColumnFinishedAt(),
			t.ColumnFinishType(),
			t.ColumnFailsMade(),
			t.ColumnCreatedAt(),
			t.ColumnUpdatedAt(),
		).
		Values(
			t.ColumnId(),
			t.ColumnEmail(),
			t.ColumnExpireAt(),
			t.ColumnFinishedAt(),
			t.ColumnFinishType(),
			t.ColumnFailsMade(),
			t.ColumnCreatedAt(),
			t.ColumnUpdatedAt(),
		).
		String(),
		// --
		row.Id,
		row.Email,
		row.ExpireAt,
		row.FinishedAt,
		row.FinishType,
		row.FailsMade,
		row.CreatedAt,
		row.UpdatedAt,
	)

	return err
}

func (t cfmTable) Update(ctx context.Context, tx *sql.Tx, row *CfmRow) error {
	_, err := tx.ExecContext(ctx, sqlbuilder.
		Update(t.Name()).
		Set(
			fmt.Sprintf("%s = ?", t.ColumnEmail()),
			fmt.Sprintf("%s = ?", t.ColumnExpireAt()),
			fmt.Sprintf("%s = ?", t.ColumnFinishedAt()),
			fmt.Sprintf("%s = ?", t.ColumnFinishType()),
			fmt.Sprintf("%s = ?", t.ColumnFailsMade()),
			fmt.Sprintf("%s = ?", t.ColumnUpdatedAt()),
		).
		Where(fmt.Sprintf("%s = ?", t.ColumnId())).
		String(),
		// --
		row.Email,
		row.ExpireAt,
		row.FinishedAt,
		row.FinishType,
		row.FailsMade,
		row.UpdatedAt,
		// --
		row.Id,
	)

	return err
}

// ---------------------------------------------------------------------------------------------------------------------

func (t cfmTable) QueryWhereId(ctx context.Context, tx *sql.Tx, id string) (*CfmRow, error) {
	row := &CfmRow{}
	err := tx.
		QueryRowContext(ctx, sqlbuilder.
			Select(
				t.ColumnId(),
				t.ColumnEmail(),
				t.ColumnExpireAt(),
				t.ColumnFinishedAt(),
				t.ColumnFinishType(),
				t.ColumnFailsMade(),
				t.ColumnCreatedAt(),
				t.ColumnUpdatedAt(),
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
			&row.ExpireAt,
			&row.FinishedAt,
			&row.FinishType,
			&row.FailsMade,
			&row.CreatedAt,
			&row.UpdatedAt,
		)

	return row, err
}

func (t cfmTable) QueryIdWhereFinishTypeEqAndExpireLess(
	ctx context.Context,
	tx *sql.Tx,
	finishType uint,
	expireMax time.Time,
	limit uint,
) ([]string, error) {
	iter, err := tx.QueryContext(ctx, sqlbuilder.
		Select(t.ColumnId()).
		From(t.Name()).
		Where(
			fmt.Sprintf("%s = ?", t.ColumnFinishType()),
			fmt.Sprintf("%s < ?", t.ColumnExpireAt()),
		).
		Limit(int(limit)).
		String(),
		// --
		finishType,
		expireMax,
		limit,
	)
	defer func() {
		if iter != nil {
			_ = iter.Close()
		}
	}()

	ids := make([]string, 0)

	if errors.Is(err, sql.ErrNoRows) {
		return ids, nil
	} else if err != nil {
		return nil, err
	}

	var id string
	for iter.Next() {
		err = iter.Scan(&id)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

// #####################################################################################################################
