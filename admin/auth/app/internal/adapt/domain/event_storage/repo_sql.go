package event_storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/event_storage"
	"example/admin/auth/internal/domain/session"
	"example/admin/auth/internal/infra/sql/db"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-txr"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ event_storage.RepositoryInterface = &RepositoryImplSql{}

type RepositoryImplSql struct {
}

const (
	EventTypeAccountCreated            = "account_created"
	EventTypeAccountDeactivated        = "account_deactivated"
	EventTypeAccountIpWhitelistChanged = "account_ip_whitelist_changed"
	EventTypeSessionCreated            = "session_created"
	EventTypeSessionClosed             = "session_closed"
)

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewRepositoryImplSql() *RepositoryImplSql {
	return &RepositoryImplSql{}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// AddAccountCreated
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (r *RepositoryImplSql) AddAccountCreated(ctx context.Context, e account.EventCreated) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[account.EventCreated]().NotEq(account.EventCreated{}).Must(e)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	dbRow := &db.EventRow{}
	dbRow.CreatedAt = time.Now()
	dbRow.Type = EventTypeAccountCreated
	dbRow.Version = 1
	dbRow.OccurredAt = e.OccurredAt()
	// --
	dbRow.ExtraAccEmail = sql.NullString{String: e.Email().String(), Valid: true}
	dbRow.ExtraAccId = sql.NullString{String: e.AccId().String(), Valid: true}

	err := db.EventTable.Insert(ctx, tx, dbRow)
	if err != nil {
		return std.WrapErrorToRuntime(err, r, "AddAccountCreated")
	}

	return nil
}

// AddAccountDeactivated
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (r *RepositoryImplSql) AddAccountDeactivated(ctx context.Context, e account.EventDeactivated) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[account.EventDeactivated]().NotEq(account.EventDeactivated{}).Must(e)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	dbRow := &db.EventRow{}
	dbRow.CreatedAt = time.Now()
	dbRow.Type = EventTypeAccountDeactivated
	dbRow.Version = 1
	dbRow.OccurredAt = e.OccurredAt()
	// --
	dbRow.ExtraAccId = sql.NullString{String: e.AccId().String(), Valid: true}

	err := db.EventTable.Insert(ctx, tx, dbRow)
	if err != nil {
		return std.WrapErrorToRuntime(err, r, "AddAccountDeactivated")
	}

	return nil
}

// AddAccountIpWhitelistChanged
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (r *RepositoryImplSql) AddAccountIpWhitelistChanged(ctx context.Context, e account.EventIpWhitelistChanged) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[account.EventIpWhitelistChanged]().NotEq(account.EventIpWhitelistChanged{}).Must(e)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	dbRow := &db.EventRow{}
	dbRow.CreatedAt = time.Now()
	dbRow.Type = EventTypeAccountIpWhitelistChanged
	dbRow.Version = 1
	dbRow.OccurredAt = e.OccurredAt()
	// --
	dbRow.ExtraAccId = sql.NullString{String: e.AccId().String(), Valid: true}

	extraAccountIpWhitelistJson, err := func() (string, error) {
		subnets := e.NewList().Subnets()
		r := make([]string, len(subnets))
		for i, sn := range subnets {
			r[i] = sn.String()
		}
		rd, err := json.Marshal(r)
		if err != nil {
			return "", err
		}
		return string(rd), nil
	}()
	if err != nil {
		return std.WrapErrorToRuntime(err, r, "AddAccountIpWhitelistChanged", "IpWhitelistJson")
	}
	dbRow.ExtraAccIpWhitelistJson = sql.NullString{String: extraAccountIpWhitelistJson, Valid: true}

	err = db.EventTable.Insert(ctx, tx, dbRow)
	if err != nil {
		return std.WrapErrorToRuntime(err, r, "AddAccountIpWhitelistChanged")
	}

	return nil
}

// AddSessionCreated
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (r *RepositoryImplSql) AddSessionCreated(ctx context.Context, e session.EventCreated) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[session.EventCreated]().NotEq(session.EventCreated{}).Must(e)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	dbRow := &db.EventRow{}
	dbRow.CreatedAt = time.Now()
	dbRow.Type = EventTypeSessionCreated
	dbRow.Version = 1
	dbRow.OccurredAt = e.OccurredAt()
	// --
	dbRow.ExtraAccId = sql.NullString{String: e.AccId().String(), Valid: true}
	dbRow.ExtraSessId = sql.NullString{String: e.SessId().String(), Valid: true}

	err := db.EventTable.Insert(ctx, tx, dbRow)
	if err != nil {
		return std.WrapErrorToRuntime(err, r, "AddSessionCreated")
	}

	return nil
}

// AddSessionClosed
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (r *RepositoryImplSql) AddSessionClosed(ctx context.Context, e session.EventClosed) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[session.EventClosed]().NotEq(session.EventClosed{}).Must(e)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	dbRow := &db.EventRow{}
	dbRow.CreatedAt = time.Now()
	dbRow.Type = EventTypeSessionClosed
	dbRow.Version = 1
	dbRow.OccurredAt = e.OccurredAt()
	// --
	dbRow.ExtraAccId = sql.NullString{String: e.AccId().String(), Valid: true}
	dbRow.ExtraSessId = sql.NullString{String: e.SessId().String(), Valid: true}

	err := db.EventTable.Insert(ctx, tx, dbRow)
	if err != nil {
		return std.WrapErrorToRuntime(err, r, "AddSessionClosed")
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
