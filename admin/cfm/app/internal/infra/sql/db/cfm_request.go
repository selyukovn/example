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

type CfmRequestRow = struct {
	CfmId       string
	Number      uint
	CodeHash    string
	RequestedAt time.Time
}

// #####################################################################################################################
// TABLE
// #####################################################################################################################

var CfmRequestTable = cfmRequestsTable{}

type cfmRequestsTable struct{}

func (t cfmRequestsTable) Name() string              { return "confirmation_request" }
func (t cfmRequestsTable) ColumnCfmId() string       { return "confirmation_id" }
func (t cfmRequestsTable) ColumnNumber() string      { return "number" }
func (t cfmRequestsTable) ColumnCodeHash() string    { return "code_hash" }
func (t cfmRequestsTable) ColumnRequestedAt() string { return "requested_at" }

// ---------------------------------------------------------------------------------------------------------------------

func (t cfmRequestsTable) InsertBulk(
	ctx context.Context,
	tx *sql.Tx,
	rows []*CfmRequestRow,
) error {
	if len(rows) == 0 {
		return nil
	}

	query := sqlbuilder.
		InsertInto(t.Name()).
		Cols(
			t.ColumnCfmId(),
			t.ColumnNumber(),
			t.ColumnCodeHash(),
			t.ColumnRequestedAt(),
		)

	valuesAsArgs := make([]interface{}, len(rows)*4)
	for i, row := range rows {
		query = query.Values(
			t.ColumnCfmId(),
			t.ColumnNumber(),
			t.ColumnCodeHash(),
			t.ColumnRequestedAt(),
		)
		valuesAsArgs[i*4+0] = row.CfmId
		valuesAsArgs[i*4+1] = row.Number
		valuesAsArgs[i*4+2] = row.CodeHash
		valuesAsArgs[i*4+3] = row.RequestedAt
	}

	_, err := tx.ExecContext(ctx, query.String(), valuesAsArgs...)

	return err
}

func (t cfmRequestsTable) DeleteWhereCfmId(
	ctx context.Context,
	tx *sql.Tx,
	cfmId string,
) error {
	_, err := tx.ExecContext(ctx, sqlbuilder.
		DeleteFrom(t.Name()).
		Where(fmt.Sprintf("%s = ?", t.ColumnCfmId())).
		String(),
		// --
		cfmId,
	)

	return err
}

// ---------------------------------------------------------------------------------------------------------------------

func (t cfmRequestsTable) QueryWhereCfmId(
	ctx context.Context,
	tx *sql.Tx,
	cfmId string,
) ([]*CfmRequestRow, error) {
	iter, err := tx.QueryContext(ctx, sqlbuilder.
		Select(
			t.ColumnNumber(),
			t.ColumnCodeHash(),
			t.ColumnRequestedAt(),
		).
		From(t.Name()).
		Where(fmt.Sprintf("%s = ?", t.ColumnCfmId())).
		String(),
		// --
		cfmId,
	)
	defer func() {
		if iter != nil {
			_ = iter.Close()
		}
	}()

	rows := make([]*CfmRequestRow, 0)

	if errors.Is(err, sql.ErrNoRows) {
		return rows, nil
	} else if err != nil {
		return rows, err
	}

	for iter.Next() {
		row := &CfmRequestRow{}
		row.CfmId = cfmId
		err = iter.Scan(
			&row.Number,
			&row.CodeHash,
			&row.RequestedAt,
		)
		rows = append(rows, row)
	}

	return rows, nil
}

// #####################################################################################################################
