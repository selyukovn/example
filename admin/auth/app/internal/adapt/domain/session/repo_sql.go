package session

import (
	"context"
	"database/sql"
	"errors"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/action_request"
	"example/admin/auth/internal/domain/client"
	"example/admin/auth/internal/domain/session"
	"example/admin/auth/internal/infra/sql/db"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-txr"
	assert "github.com/selyukovn/go-wm-assert"
	"net/netip"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ session.RepositoryInterface = RepositoryImplSql{}

type RepositoryImplSql struct {
	fnIsDuplicateKeyError func(error) bool
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewRepositoryImplSql(fnIsDuplicateKeyError func(error) bool) RepositoryImplSql {
	return RepositoryImplSql{
		fnIsDuplicateKeyError: fnIsDuplicateKeyError,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Mapping
// ---------------------------------------------------------------------------------------------------------------------

func (r RepositoryImplSql) mapSessionToDbRow(s *session.Session) *db.SessionRow {
	dbRow := &db.SessionRow{}

	sId,
		sAccId,
		sSignInId,
		sInitCl,
		sInitAt,
		sExpAt,
		sClosedAt := session.ReflectExtract(s)

	dbRow.Id = sId.String()
	dbRow.AccId = sAccId.String()
	dbRow.SignInId = sSignInId.String()
	dbRow.InitialClientUserAgent = sInitCl.UserAgent().String()
	dbRow.InitialClientIp = sInitCl.IpAddress().String()
	dbRow.InitiatedAt = sInitAt
	dbRow.ExpireAt = sExpAt

	dbRow.IsClosed = false
	dbRow.ClosedAt = sql.NullTime{}
	if !sClosedAt.IsZero() {
		dbRow.IsClosed = true
		dbRow.ClosedAt = sql.NullTime{Time: sClosedAt, Valid: true}
	}

	return dbRow
}

func (r RepositoryImplSql) mapDbRowToSession(dbRow *db.SessionRow) (*session.Session, error) {
	var err error

	sId, err := session.IdFromString(dbRow.Id)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "mapDbRowToSession", "id")
	}

	sAccId, err := account.IdFromString(dbRow.AccId)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "mapDbRowToSession", "AccId")
	}

	sSignInId, err := action_request.IdFromString(dbRow.SignInId)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "mapDbRowToSession", "SignInId")
	}

	clUa, err := client.UserAgentFromString(dbRow.InitialClientUserAgent)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "mapDbRowToSession", "InitialClientUserAgent")
	}

	clIp, err := netip.ParseAddr(dbRow.InitialClientIp)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "mapDbRowToSession", "InitialClientIp")
	}

	sInitCl, err := client.NewClient(clUa, clIp)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "mapDbRowToSession", "initialClient")
	}

	sInitAt := dbRow.InitiatedAt
	sExpAt := dbRow.ExpireAt

	sClosedAt := time.Time{}
	if dbRow.ClosedAt.Valid {
		sClosedAt = dbRow.ClosedAt.Time
	}

	s := session.ReflectRestore(sId, sAccId, sSignInId, sInitCl, sInitAt, sExpAt, sClosedAt)

	return s, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Add
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorAlreadyDone -- если с таким id или signInRequestId уже существует
//   - std.ErrorRuntime
func (r RepositoryImplSql) Add(ctx context.Context, s *session.Session) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*session.Session]().NotEq(nil).Must(s)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	now := time.Now()

	dbRow := r.mapSessionToDbRow(s)
	dbRow.CreatedAt = now
	dbRow.UpdatedAt = now

	err := db.SessionTable.Insert(ctx, tx, dbRow)

	if r.fnIsDuplicateKeyError(err) {
		return std.NewErrorAlreadyDoneFf(
			"Сессия с id %q или signInRequestId %q уже существует: %v",
			dbRow.Id,
			dbRow.SignInId,
			err,
		)
	} else if err != nil {
		return std.WrapErrorToRuntime(err, r, "Add")
	}

	return nil
}

// Update
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (r RepositoryImplSql) Update(ctx context.Context, s *session.Session) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*session.Session]().NotEq(nil).Must(s)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	now := time.Now()

	dbRow := r.mapSessionToDbRow(s)
	dbRow.UpdatedAt = now

	err := db.SessionTable.Update(ctx, tx, dbRow)
	if err != nil {
		return std.WrapErrorToRuntime(err, r, "Update")
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

// GetById
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - std.ErrorRuntime
func (r RepositoryImplSql) GetById(ctx context.Context, id session.Id) (*session.Session, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[session.Id]().NotEq(session.IdNil).Must(id)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	dbRow, err := db.SessionTable.QueryWhereId(ctx, tx, id.String())

	if errors.Is(err, sql.ErrNoRows) {
		return nil, std.NewErrorNotFoundFf("Сессия %q не найдена: %v", id, err)
	} else if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "GetById")
	}

	sess, err := r.mapDbRowToSession(dbRow)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "GetById")
	}

	return sess, nil
}

// GetIdsOfGoingToExpire
//
// Cм. Session.IsClosed
//
// Паникует при нулевых аргументах:
// - ctx
//
// Ошибки:
//   - std.ErrorRuntime
func (r RepositoryImplSql) GetIdsOfGoingToExpire(ctx context.Context, now time.Time, limit uint) ([]session.Id, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	dbIds, err := db.SessionTable.QueryIdWhereNoClosedAtAndExpireLess(ctx, tx, now, limit)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "GetIdsOfGoingToExpire")
	}

	ids := make([]session.Id, len(dbIds))

	for i, id := range dbIds {
		ids[i], err = session.IdFromString(id)
		if err != nil {
			return []session.Id{}, std.WrapErrorToRuntime(err, r, "GetIdsOfGoingToExpire")
		}
	}

	return ids, nil
}

// HasBySignInRequest
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (r RepositoryImplSql) HasBySignInRequest(ctx context.Context, signInRequestId action_request.Id) (bool, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[action_request.Id]().NotEq(action_request.IdNil).Must(signInRequestId)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	sessId, err := db.SessionTable.QuerySessIdWhereSignInRequestId(ctx, tx, signInRequestId.String())

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, std.WrapErrorToRuntime(err, r, "HasBySignInRequest")
	}

	return sessId != "", nil
}

// GetAccIdById
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - std.ErrorRuntime
func (r RepositoryImplSql) GetAccIdById(ctx context.Context, id session.Id) (account.Id, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[session.Id]().NotEq(session.IdNil).Must(id)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	accIdStr, err := db.SessionTable.QueryAccIdWhereId(ctx, tx, id.String())

	if errors.Is(err, sql.ErrNoRows) {
		return account.IdNil, std.NewErrorNotFoundFf("Аккаунт не найден для Сессии %q", id)
	} else if err != nil {
		return account.IdNil, std.WrapErrorToRuntime(err, r, "GetAccIdById")
	}

	accId, err := account.IdFromString(accIdStr)
	if err != nil {
		return account.IdNil, std.WrapErrorToRuntime(err, r, "GetAccIdById", "accId")
	}

	return accId, nil
}

// GetAccIdAndExpAtById
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - std.ErrorRuntime
func (r RepositoryImplSql) GetAccIdAndExpAtById(ctx context.Context, id session.Id) (account.Id, time.Time, error) {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(id.IsNil())

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	accIdStr, expAt, err := db.SessionTable.QueryAccIdAndExpireAtWhereId(ctx, tx, id.String())

	if errors.Is(err, sql.ErrNoRows) {
		return account.IdNil, time.Time{}, std.NewErrorNotFoundFf("Аккаунт не найден для Сессии %q", id)
	} else if err != nil {
		return account.IdNil, time.Time{}, std.WrapErrorToRuntime(err, r, "GetAccIdById")
	}

	accId, err := account.IdFromString(accIdStr)
	if err != nil {
		return account.IdNil, time.Time{}, std.WrapErrorToRuntime(err, r, "GetAccIdById", "accId")
	}

	return accId, expAt, nil
}

// ---------------------------------------------------------------------------------------------------------------------
