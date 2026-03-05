package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// #####################################################################################################################
// ROW
// #####################################################################################################################

type SessionRow struct {
	Id                     string
	AccId                  string
	SignInId               string
	InitialClientUserAgent string
	InitialClientIp        string
	InitiatedAt            time.Time
	ExpireAt               time.Time
	IsClosed               bool
	ClosedAt               sql.NullTime
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// #####################################################################################################################
// TABLE
// #####################################################################################################################

var SessionTable = sessionTable{}

// ---------------------------------------------------------------------------------------------------------------------

type sessionTable struct{}

func (t sessionTable) Name() string                         { return "session" }
func (t sessionTable) ColumnId() string                     { return "id" }
func (t sessionTable) ColumnAccId() string                  { return "account_id" }
func (t sessionTable) ColumnSignInId() string               { return "sign_in_request_id" }
func (t sessionTable) ColumnInitialClientUserAgent() string { return "initial_client_user_agent" }
func (t sessionTable) ColumnInitialClientIp() string        { return "initial_client_ip" }
func (t sessionTable) ColumnInitiatedAt() string            { return "initiated_at" }
func (t sessionTable) ColumnExpireAt() string               { return "expire_at" }
func (t sessionTable) ColumnIsClosed() string               { return "is_closed" }
func (t sessionTable) ColumnClosedAt() string               { return "closed_at" }
func (t sessionTable) ColumnCreatedAt() string              { return "created_at" }
func (t sessionTable) ColumnUpdatedAt() string              { return "updated_at" }

// ---------------------------------------------------------------------------------------------------------------------

// Insert
//
// Паникует при нулевых аргументах.
func (t sessionTable) Insert(ctx context.Context, tx *sql.Tx, row *SessionRow) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*sql.Tx]().NotEq(nil).Must(tx)
	assert.Cmp[*SessionRow]().NotEq(nil).Must(row)

	_, err := tx.ExecContext(ctx, sqlbuilder.
		InsertInto(t.Name()).
		Cols(
			t.ColumnId(),
			t.ColumnAccId(),
			t.ColumnSignInId(),
			t.ColumnInitialClientUserAgent(),
			t.ColumnInitialClientIp(),
			t.ColumnInitiatedAt(),
			t.ColumnExpireAt(),
			t.ColumnIsClosed(),
			t.ColumnClosedAt(),
			t.ColumnCreatedAt(),
			t.ColumnUpdatedAt(),
		).
		Values(
			t.ColumnId(),
			t.ColumnAccId(),
			t.ColumnSignInId(),
			t.ColumnInitialClientUserAgent(),
			t.ColumnInitialClientIp(),
			t.ColumnInitiatedAt(),
			t.ColumnExpireAt(),
			t.ColumnIsClosed(),
			t.ColumnClosedAt(),
			t.ColumnCreatedAt(),
			t.ColumnUpdatedAt(),
		).
		String(),
		// --
		row.Id,
		row.AccId,
		row.SignInId,
		row.InitialClientUserAgent,
		row.InitialClientIp,
		row.InitiatedAt,
		row.ExpireAt,
		row.IsClosed,
		row.ClosedAt,
		row.CreatedAt,
		row.UpdatedAt,
	)

	return err
}

// Update
//
// Паникует при нулевых аргументах.
func (t sessionTable) Update(ctx context.Context, tx *sql.Tx, row *SessionRow) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*sql.Tx]().NotEq(nil).Must(tx)
	assert.Cmp[*SessionRow]().NotEq(nil).Must(row)

	_, err := tx.ExecContext(ctx, sqlbuilder.
		Update(t.Name()).
		Set(
			fmt.Sprintf("%s = ?", t.ColumnAccId()),
			fmt.Sprintf("%s = ?", t.ColumnSignInId()),
			fmt.Sprintf("%s = ?", t.ColumnInitialClientUserAgent()),
			fmt.Sprintf("%s = ?", t.ColumnInitialClientIp()),
			fmt.Sprintf("%s = ?", t.ColumnInitiatedAt()),
			fmt.Sprintf("%s = ?", t.ColumnExpireAt()),
			fmt.Sprintf("%s = ?", t.ColumnIsClosed()),
			fmt.Sprintf("%s = ?", t.ColumnClosedAt()),
			fmt.Sprintf("%s = ?", t.ColumnUpdatedAt()),
		).
		Where(fmt.Sprintf("%s = ?", t.ColumnId())).
		String(),
		// --
		row.AccId,
		row.SignInId,
		row.InitialClientUserAgent,
		row.InitialClientIp,
		row.InitiatedAt,
		row.ExpireAt,
		row.IsClosed,
		row.ClosedAt,
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
func (t sessionTable) QueryWhereId(ctx context.Context, tx *sql.Tx, id string) (*SessionRow, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*sql.Tx]().NotEq(nil).Must(tx)

	row := &SessionRow{}
	err := tx.
		QueryRowContext(ctx, sqlbuilder.
			Select(
				t.ColumnId(),
				t.ColumnAccId(),
				t.ColumnSignInId(),
				t.ColumnInitialClientUserAgent(),
				t.ColumnInitialClientIp(),
				t.ColumnInitiatedAt(),
				t.ColumnExpireAt(),
				t.ColumnIsClosed(),
				t.ColumnClosedAt(),
			).
			From(t.Name()).
			Where(fmt.Sprintf("%s = ?", t.ColumnId())).
			String(),
			// --
			id,
		).
		Scan(
			&row.Id,
			&row.AccId,
			&row.SignInId,
			&row.InitialClientUserAgent,
			&row.InitialClientIp,
			&row.InitiatedAt,
			&row.ExpireAt,
			&row.IsClosed,
			&row.ClosedAt,
		)

	return row, err
}

// QueryIdWhereNoClosedAtAndExpireLess
//
// Паникует при нулевых аргументах:
//   - ctx
//   - tx
func (t sessionTable) QueryIdWhereNoClosedAtAndExpireLess(
	ctx context.Context,
	tx *sql.Tx,
	expireMax time.Time,
	limit uint,
) ([]string, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*sql.Tx]().NotEq(nil).Must(tx)

	iter, err := tx.QueryContext(ctx, sqlbuilder.
		Select(t.ColumnId()).
		From(t.Name()).
		Where(
			fmt.Sprintf("%s = ?", t.ColumnIsClosed()),
			fmt.Sprintf("%s < ?", t.ColumnExpireAt()),
		).
		Limit(int(limit)).
		String(),
		// --
		false,
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

// QuerySessIdWhereSignInRequestId
//
// Паникует при нулевых аргументах:
//   - ctx
//   - tx
func (t sessionTable) QuerySessIdWhereSignInRequestId(
	ctx context.Context,
	tx *sql.Tx,
	signInRequestId string,
) (string, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*sql.Tx]().NotEq(nil).Must(tx)

	row := &SessionRow{}
	err := tx.
		QueryRowContext(ctx, sqlbuilder.
			Select(t.ColumnId()).
			From(t.Name()).
			Where(fmt.Sprintf("%s = ?", t.ColumnSignInId())).
			String(),
			// --
			signInRequestId,
		).
		Scan(&row.Id)

	return row.Id, err
}

// QueryAccIdWhereId
//
// Паникует при нулевых аргументах:
//   - ctx
//   - tx
func (t sessionTable) QueryAccIdWhereId(ctx context.Context, tx *sql.Tx, id string) (string, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*sql.Tx]().NotEq(nil).Must(tx)

	row := &SessionRow{}
	err := tx.
		QueryRowContext(ctx, sqlbuilder.
			Select(t.ColumnAccId()).
			From(t.Name()).
			Where(fmt.Sprintf("%s = ?", t.ColumnId())).
			String(),
			// --
			id,
		).
		Scan(&row.AccId)

	return row.AccId, err
}

// QueryAccIdAndExpireAtWhereId
//
// Паникует при нулевых аргументах:
//   - ctx
//   - tx
func (t sessionTable) QueryAccIdAndExpireAtWhereId(ctx context.Context, tx *sql.Tx, id string) (string, time.Time, error) {
	assert.NotNilDeepMust(ctx)
	assert.NotNilDeepMust(tx)

	row := &SessionRow{}
	err := tx.
		QueryRowContext(ctx, sqlbuilder.
			Select(t.ColumnAccId(), t.ColumnExpireAt()).
			From(t.Name()).
			Where(fmt.Sprintf("%s = ?", t.ColumnId())).
			String(),
			// --
			id,
		).
		Scan(&row.AccId, &row.ExpireAt)

	return row.AccId, row.ExpireAt, err
}

// ---------------------------------------------------------------------------------------------------------------------
