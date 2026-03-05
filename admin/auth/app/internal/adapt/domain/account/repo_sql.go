package account

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"example/admin/auth/internal/domain/account"
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

type RepositoryImplSql struct {
	fnIsDuplicateKeyError func(error) bool
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewRepositoryImplSql(fnIsDuplicateKeyError func(error) bool) *RepositoryImplSql {
	return &RepositoryImplSql{
		fnIsDuplicateKeyError: fnIsDuplicateKeyError,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Mapping
// ---------------------------------------------------------------------------------------------------------------------

// mapAccountToDbRow
//
//   - std.ErrorRuntime
func (r *RepositoryImplSql) mapAccountToDbRow(a *account.Account) (*db.AccountRow, error) {
	var err error
	dbRow := &db.AccountRow{}

	aId, aEmail, aDeactivatedAt, aIpWhitelist := account.ReflectExtract(a)

	dbRow.Id = aId.String()
	dbRow.Email = aEmail.String()
	dbRow.IsActive = aDeactivatedAt.IsZero()

	dbRow.DeactivatedAt = sql.NullTime{}
	if !aDeactivatedAt.IsZero() {
		dbRow.DeactivatedAt = sql.NullTime{aDeactivatedAt, true}
	}

	dbRow.IpWhitelistJson, err = func(wl account.IpWhitelist) (string, error) {
		subnets := wl.Subnets()
		strings := make([]string, len(subnets))
		for i, s := range subnets {
			strings[i] = s.String()
		}
		jBytes, err := json.Marshal(strings)
		if err != nil {
			return "", std.WrapErrorToRuntime(err, r, "mapAccountToDbRow", "IpWhitelist")
		}
		return string(jBytes), nil
	}(aIpWhitelist)
	if err != nil {
		return nil, err
	}

	return dbRow, nil
}

// mapDbRowToAccount
//
//   - std.ErrorRuntime
func (r *RepositoryImplSql) mapDbRowToAccount(dbRow *db.AccountRow) (*account.Account, error) {
	var err error

	id, err := account.IdFromString(dbRow.Id)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "mapDbRowToAccount", "Id")
	}

	email, err := std.EmailFromString(dbRow.Email)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "mapDbRowToAccount", "Email")
	}

	deactivatedAt := time.Time{}
	if dbRow.DeactivatedAt.Valid {
		deactivatedAt = dbRow.DeactivatedAt.Time
	}

	ipWhitelist, err := func(IpWhitelistJson string) (account.IpWhitelist, error) {
		var strings []string
		b := ([]byte)(IpWhitelistJson)
		err := json.Unmarshal(b, &strings)
		if err != nil {
			return account.IpWhitelistEmpty, err
		}
		subnets := make([]netip.Prefix, len(strings))
		for i, s := range strings {
			p, err := netip.ParsePrefix(s)
			if err != nil {
				return account.IpWhitelistEmpty, err
			}
			subnets[i] = p
		}
		iwl, err := account.IpWhitelistFromPrefixes(subnets)
		if err != nil {
			return account.IpWhitelistEmpty, err
		}
		return iwl, nil
	}(dbRow.IpWhitelistJson)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "mapDbRowToAccount", "IpWhitelist")
	}

	a := account.ReflectRestore(id, email, deactivatedAt, ipWhitelist)

	return a, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Add
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorAlreadyDone -- если с таким id или email уже существует
//   - std.ErrorRuntime
func (r *RepositoryImplSql) Add(ctx context.Context, a *account.Account) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*account.Account]().NotEq(nil).Must(a)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	now := time.Now()

	dbRow, err := r.mapAccountToDbRow(a)
	if err != nil {
		return std.WrapErrorToRuntime(err, r, "Add")
	}
	dbRow.CreatedAt = now
	dbRow.UpdatedAt = now

	err = db.AccountTable.Insert(ctx, tx, dbRow)

	if r.fnIsDuplicateKeyError(err) {
		return std.NewErrorAlreadyDoneFf("Аккаунт %q или %q уже сушествует: %v", dbRow.Id, dbRow.Email, err)
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
func (r *RepositoryImplSql) Update(ctx context.Context, a *account.Account) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[*account.Account]().NotEq(nil).Must(a)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	now := time.Now()

	dbRow, err := r.mapAccountToDbRow(a)
	if err != nil {
		return std.WrapErrorToRuntime(err, r, "Update")
	}
	dbRow.UpdatedAt = now

	err = db.AccountTable.Update(ctx, tx, dbRow)
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
func (r *RepositoryImplSql) GetById(ctx context.Context, id account.Id) (*account.Account, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[account.Id]().NotEq(account.IdNil).Must(id)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	dbRow, err := db.AccountTable.QueryWhereId(ctx, tx, id.String())

	if errors.Is(err, sql.ErrNoRows) {
		return nil, std.NewErrorNotFoundFf("Аккаунт %q не найден: %v", id, err)
	} else if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "GetById")
	}

	acc, err := r.mapDbRowToAccount(dbRow)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "GetById")
	}

	return acc, nil
}

// GetByEmail
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - std.ErrorRuntime
func (r *RepositoryImplSql) GetByEmail(ctx context.Context, email std.Email) (*account.Account, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[std.Email]().NotEq(std.EmailNil).Must(email)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	dbRow, err := db.AccountTable.QueryWhereEmail(ctx, tx, email.String())

	if errors.Is(err, sql.ErrNoRows) {
		return nil, std.NewErrorNotFoundFf("Аккаунт %q не найден: %v", email, err)
	} else if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "GetByEmail")
	}

	acc, err := r.mapDbRowToAccount(dbRow)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "GetByEmail")
	}

	return acc, nil
}

// GetEmailById
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - std.ErrorRuntime
func (r *RepositoryImplSql) GetEmailById(ctx context.Context, id account.Id) (std.Email, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[account.Id]().NotEq(account.IdNil).Must(id)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	str, err := db.AccountTable.QueryEmailWhereId(ctx, tx, id.String())

	if errors.Is(err, sql.ErrNoRows) {
		return std.EmailNil, std.NewErrorNotFoundFf("Аккаунт %q не найден: %v", id, err)
	} else if err != nil {
		return std.EmailNil, std.WrapErrorToRuntime(err, r, "GetEmailById")
	}

	email, err := std.EmailFromString(str)
	if err != nil {
		return std.EmailNil, std.WrapErrorToRuntime(err, r, "GetEmailById")
	}

	return email, nil
}

// ---------------------------------------------------------------------------------------------------------------------
