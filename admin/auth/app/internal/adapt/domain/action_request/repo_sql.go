package action_request

import (
	"context"
	"database/sql"
	"errors"
	"example/admin/auth/internal/domain/account"
	"example/admin/auth/internal/domain/action_request"
	"example/admin/auth/internal/domain/cfm"
	"example/admin/auth/internal/infra/sql/db"
	"fmt"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-txr"
	assert "github.com/selyukovn/go-wm-assert"
	"strings"
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

const (
	rowTypeSignIn uint = 1
)

func (r *RepositoryImplSql) mapSignInToDbRow(signIn *action_request.SignIn) *db.ActionRequestRow {
	dbRow := &db.ActionRequestRow{}

	sId, sAccId, sCfmId, sReqAt := action_request.ReflectExtract(signIn)

	dbRow.Id = sId.String()
	dbRow.Type = rowTypeSignIn
	dbRow.AccId = sAccId.String()
	dbRow.CfmId = sCfmId.String()
	dbRow.RequestedAt = sReqAt

	return dbRow
}

func (r *RepositoryImplSql) mapDbRowToSignIn(dbRow *db.ActionRequestRow) (*action_request.SignIn, error) {
	var err error

	sId, err := action_request.IdFromString(dbRow.Id)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "mapDbRowToSignIn", "id")
	}

	sAccId, err := account.IdFromString(dbRow.AccId)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "mapDbRowToSignIn", "AccId")
	}

	sCfmId, err := cfm.IdFromString(dbRow.CfmId)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "mapDbRowToSignIn", "CfmId")
	}

	sReqAt := dbRow.RequestedAt

	signIn := action_request.ReflectRestore(sId, sAccId, sCfmId, sReqAt)

	return signIn, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// Add
//
// Паникует, если:
//   - ctx == nil
//   - actReq == nil или не входит в набор типов: *SignIn
//
// Ошибки:
//   - std.ErrorAlreadyDone -- если с таким id уже существует
//   - std.ErrorRuntime
func (r *RepositoryImplSql) Add(ctx context.Context, actReq any) error {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)

	var dbRow *db.ActionRequestRow

	switch actReq.(type) {
	case *action_request.SignIn:
		dbRow = r.mapSignInToDbRow(actReq.(*action_request.SignIn))
	default:
		panic(fmt.Errorf(
			"%T.Add требует actReq один из типов: %s, а не %T",
			r,
			strings.Join([]string{"*SignIn"}, ","),
			actReq,
		))
	}

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	now := time.Now()

	dbRow.CreatedAt = now
	dbRow.UpdatedAt = now

	err := db.ActionRequestTable.Insert(ctx, tx, dbRow)

	if r.fnIsDuplicateKeyError(err) {
		return std.NewErrorAlreadyDoneFf("Запрос на действие %q уже существует: %v", dbRow.Id, err)
	} else if err != nil {
		return std.WrapErrorToRuntime(err, r, "Add")
	}

	return nil
}

// Update
//
// Паникует, если:
//   - ctx == nil
//   - actReq == nil или не входит в набор типов: *SignIn
//
// Ошибки:
//   - std.ErrorRuntime
func (r *RepositoryImplSql) Update(ctx context.Context, actReq any) error {
	var dbRow *db.ActionRequestRow

	switch actReq.(type) {
	case *action_request.SignIn:
		dbRow = r.mapSignInToDbRow(actReq.(*action_request.SignIn))
	default:
		panic(fmt.Errorf(
			"%T.Update требует actReq один из типов: %s, а не %T",
			r,
			strings.Join([]string{"*SignIn"}, ","),
			actReq,
		))
	}

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	now := time.Now()

	dbRow.UpdatedAt = now

	err := db.ActionRequestTable.Update(ctx, tx, dbRow)

	if err != nil {
		return std.WrapErrorToRuntime(err, r, "Update")
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

// GetSignIn
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - std.ErrorRuntime
func (r *RepositoryImplSql) GetSignIn(ctx context.Context, actReqId action_request.Id) (*action_request.SignIn, error) {
	assert.Cmp[context.Context]().NotEq(nil).Must(ctx)
	assert.Cmp[action_request.Id]().NotEq(action_request.IdNil).Must(actReqId)

	tx := txr.TxFromCtx(ctx).(*sql.Tx)

	dbRow, err := db.ActionRequestTable.QueryWhereIdAndType(ctx, tx, actReqId.String(), rowTypeSignIn)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, std.NewErrorNotFoundFf("SignIn %q не найден: %v", actReqId, err)
	} else if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "GetSignIn")
	}

	actReq, err := r.mapDbRowToSignIn(dbRow)
	if err != nil {
		return nil, std.WrapErrorToRuntime(err, r, "GetSignIn")
	}

	return actReq, nil
}

// ---------------------------------------------------------------------------------------------------------------------
