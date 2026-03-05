package cfm

import (
	"context"
	"example/admin/auth/internal/domain/cfm"
	"github.com/selyukovn/go-id/like_uuid"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type ServiceImplDummy struct{}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewServiceImplDummy() *ServiceImplDummy {
	return &ServiceImplDummy{}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

// CreateForEmail
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorRuntime
func (s *ServiceImplDummy) CreateForEmail(ctx context.Context, email std.Email) (cfm.ServiceResultCreate, error) {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(email.IsNil())

	idV, err := like_uuid.NewIdGeneratorUniqueRandom().Generate()
	if err != nil {
		return cfm.ServiceResultCreateNil, std.WrapErrorToRuntime(err, s, "CreateForEmail")
	}
	id, err := cfm.IdFromString(idV.String())
	if err != nil {
		return cfm.ServiceResultCreateNil, std.WrapErrorToRuntime(err, s, "CreateForEmail")
	}

	return cfm.NewServiceResultCreate(email, id, time.Now().Add(10*time.Minute)), nil
}

// Request
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - cfm.ErrorFinished
//   - cfm.ErrorNoAttemptsLeft
//   - cfm.ErrorRequestsFrequency
//   - std.ErrorRuntime
func (s *ServiceImplDummy) Request(ctx context.Context, cfmId cfm.Id) (cfm.ServiceResultRequest, error) {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(cfmId.IsNil())

	return cfm.NewServiceResultRequestCanAgain(
		cfmId,
		1,
		time.Now().Add(1*time.Minute),
	), nil
}

// Confirm
//
// Паникует при нулевых аргументах.
//
// Ошибки:
//   - std.ErrorNotFound
//   - cfm.ErrorFinished
//   - std.ErrorUnprocessable -- если не была запрошена
//   - std.ErrorRuntime
func (s *ServiceImplDummy) Confirm(ctx context.Context, cfmId cfm.Id, code cfm.Code) (cfm.ServiceResultConfirm, error) {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(cfmId.IsNil())
	assert.FalseMust(code.IsNil())

	return cfm.NewServiceResultConfirmLast(
		cfmId,
		time.Now().Add(-1*time.Second),
		true,
	), nil
}
