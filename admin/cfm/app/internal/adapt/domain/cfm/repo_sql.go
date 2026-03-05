package cfm

import (
	"context"
	"database/sql"
	"errors"
	"example/admin/cfm/internal/domain/cfm"
	"example/admin/cfm/internal/domain/cfm/code"
	"example/admin/cfm/internal/infra/sql/db"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-txr"
	assert "github.com/selyukovn/go-wm-assert"
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

// NewRepositoryImplSql
//
// Паникует при нулевых аргументах.
func NewRepositoryImplSql(fnIsDuplicateKeyError func(error) bool) *RepositoryImplSql {
	assert.NotNilDeepMust(fnIsDuplicateKeyError)

	return &RepositoryImplSql{
		fnIsDuplicateKeyError: fnIsDuplicateKeyError,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Mapping
// ---------------------------------------------------------------------------------------------------------------------

func (r *RepositoryImplSql) mapCfmToDbRows(c *cfm.Cfm) (
	*db.CfmRow,
	[]*db.CfmRequestRow,
) {
	cId, cEmail, cExpiredAt, cFinishedAt, cFinishType, cRequests, cFailsMade := cfm.ReflectExtract(c)

	// --

	cDbRow := &db.CfmRow{}

	cDbRow.Id = cId.String()
	cDbRow.Email = cEmail.String()
	cDbRow.ExpireAt = cExpiredAt

	cDbRow.FinishedAt = sql.NullTime{}
	if !cFinishedAt.IsZero() {
		cDbRow.FinishedAt = sql.NullTime{cFinishedAt, true}
	}

	cDbRow.FinishType = uint(cFinishType)

	crDbRows := make([]*db.CfmRequestRow, len(cRequests))
	for i, req := range cRequests {
		crDbRows[i] = &db.CfmRequestRow{
			CfmId:       cId.String(),
			Number:      req.Number,
			CodeHash:    req.CodeHash.String(),
			RequestedAt: req.RequestedAt,
		}
	}

	cDbRow.FailsMade = cFailsMade

	return cDbRow, crDbRows
}

func (r *RepositoryImplSql) mapDbRowToCfm(
	cDbRow *db.CfmRow,
	crDbRows []*db.CfmRequestRow,
) (*cfm.Cfm, error) {
	cId, err := cfm.IdFromString(cDbRow.Id)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "mapDbRowToCfm", "cId")
	}

	cEmail, err := std.EmailFromString(cDbRow.Email)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "mapDbRowToCfm", "cEmail")
	}

	cExpiredAt := cDbRow.ExpireAt

	cFinishedAt := time.Time{}
	if cDbRow.FinishedAt.Valid {
		cFinishedAt = cDbRow.FinishedAt.Time
	}

	cFinishType := cfm.FinishType(cDbRow.FinishType)

	cRequests := make([]*cfm.CfmRequest, len(crDbRows))
	for i, cr := range crDbRows {
		codeHash, err := code.HashFromString(cr.CodeHash)
		if err != nil {
			return nil, std.WrapErrorToRuntime(err, r, "mapDbRowToCfm", "cRequests", "codeHash")
		}

		cRequests[i] = &cfm.CfmRequest{
			Number:      cr.Number,
			CodeHash:    codeHash,
			RequestedAt: cr.RequestedAt,
		}
	}

	cFailsMade := cDbRow.FailsMade

	c := cfm.ReflectRestore(
		cId,
		cEmail,
		cExpiredAt,
		cFinishedAt,
		cFinishType,
		cRequests,
		cFailsMade,
	)

	return c, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Add
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorAlreadyDone -- если с таким id уже существует
//   - std.ErrorRuntime
func (r *RepositoryImplSql) Add(ctx context.Context, c *cfm.Cfm) error {
	assert.NotNilDeepMust(ctx)
	assert.NotNilDeepMust(c)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	now := time.Now()
	cDbRow, crRows := r.mapCfmToDbRows(c)
	cDbRow.CreatedAt = now
	cDbRow.UpdatedAt = now

	// --

	err := db.CfmTable.Insert(ctx, tx, cDbRow)
	if r.fnIsDuplicateKeyError(err) {
		return std.NewErrorAlreadyDoneFf("Конфирмация с id %q уже существует: %v", cDbRow.Id, err)
	} else if err != nil {
		return std.WrapErrorToRuntime(err, r, "Add", "CfmTable")
	}

	// --

	err = db.CfmRequestTable.InsertBulk(ctx, tx, crRows)
	if err != nil {
		return std.WrapErrorToRuntime(err, r, "Add", "CfmRequestTable", "InsertBulk")
	}

	// --

	return nil
}

// Update
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (r *RepositoryImplSql) Update(ctx context.Context, c *cfm.Cfm) error {
	assert.NotNilDeepMust(ctx)
	assert.NotNilDeepMust(c)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	now := time.Now()
	cDbRow, crRows := r.mapCfmToDbRows(c)
	cDbRow.UpdatedAt = now

	// --

	err := db.CfmTable.Update(ctx, tx, cDbRow)
	if err != nil {
		return std.WrapErrorToRuntime(err, r, "Update", "CfmTable")
	}

	// --

	err = db.CfmRequestTable.DeleteWhereCfmId(ctx, tx, cDbRow.Id)
	if err != nil {
		return std.WrapErrorToRuntime(err, r, "Update", "CfmRequestTable", "DeleteWhereCfmId")
	}

	err = db.CfmRequestTable.InsertBulk(ctx, tx, crRows)
	if err != nil {
		return std.WrapErrorToRuntime(err, r, "Update", "CfmRequestTable", "InsertBulk")
	}

	// --

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
func (r *RepositoryImplSql) GetById(ctx context.Context, id cfm.Id) (*cfm.Cfm, error) {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(id.IsNil())

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	cRow, err := db.CfmTable.QueryWhereId(ctx, tx, id.String())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, std.NewErrorNotFoundFf("Конфирмация %q не найдена: %v", id, err)
	} else if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "GetById", "CfmTable")
	}

	crRows, err := db.CfmRequestTable.QueryWhereCfmId(ctx, tx, id.String())
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "GetById", "CfmRequestTable")
	}

	c, err := r.mapDbRowToCfm(cRow, crRows)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "GetById", "mapDbRowToCfm")
	}

	return c, nil
}

// GetIdsOfGoingToExpire
//
// Cм. Cfm.IsFinished()
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (r *RepositoryImplSql) GetIdsOfGoingToExpire(
	ctx context.Context,
	now time.Time,
	limit uint,
) ([]cfm.Id, error) {
	assert.NotNilDeepMust(ctx)
	assert.NotZeroMust(now)
	assert.NotZeroMust(limit)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	dbIds, err := db.CfmTable.QueryIdWhereFinishTypeEqAndExpireLess(
		ctx,
		tx,
		uint(cfm.FinishNotYet),
		now,
		limit,
	)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "GetIdsOfGoingToExpire")
	}

	ids := make([]cfm.Id, len(dbIds))

	for i, dbId := range dbIds {
		id, err := cfm.IdFromString(dbId)
		if err != nil {
			return nil, std.WrapErrorToRuntime(err, r, "GetIdsOfGoingToExpire", "id="+dbId)
		}

		ids[i] = id
	}

	return ids, nil
}
