package cfm

import (
	"context"
	"example/admin/auth/internal/domain/cfm"
	client "example/admin/auth/internal/infra/clients/cfm"
	infra_client_cfm "example/admin/auth/internal/infra/clients/cfm"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type ServiceImplCfmService struct {
	client infra_client_cfm.ClientInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewServiceImplCfmService
//
// Паникует при нулевых аргументах.
func NewServiceImplCfmService(client infra_client_cfm.ClientInterface) *ServiceImplCfmService {
	assert.NotNilDeepMust(client)

	return &ServiceImplCfmService{
		client: client,
	}
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
func (s *ServiceImplCfmService) CreateForEmail(ctx context.Context, email std.Email) (cfm.ServiceResultCreate, error) {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(email.IsNil())

	nilRes := cfm.ServiceResultCreateNil

	clRes, clErr := s.client.CreateForEmail(ctx, email)

	// error
	// ----------------

	if clErr != nil {
		return nilRes, std.WrapErrorToRuntime(clErr, s, "CreateForEmail")
	}

	// success
	// ----------------

	cfmId, err := cfm.IdFromString(clRes.CfmId)
	if err != nil {
		return nilRes, std.WrapErrorToRuntime(err, s, "CreateForEmail", "CfmId")
	}

	return cfm.NewServiceResultCreate(email, cfmId, clRes.ExpireAt), nil
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
func (s *ServiceImplCfmService) Request(ctx context.Context, cfmId cfm.Id) (cfm.ServiceResultRequest, error) {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(cfmId.IsNil())

	nilRes := cfm.ServiceResultRequestNil

	clRes, clErr := s.client.Request(ctx, cfmId.String())
	switch vErr := clErr.(type) {
	case nil:
	case std.ErrorNotFound:
		return nilRes, vErr
	case client.ErrorFinished:
		if vErr.IsAsPassed {
			return nilRes, cfm.NewErrorFinishedAsPassed(cfmId, vErr.FinishedAt)
		} else if vErr.IsAsFailed {
			return nilRes, cfm.NewErrorFinishedAsFailed(cfmId, vErr.FinishedAt)
		} else if vErr.IsAsExpired {
			return nilRes, cfm.NewErrorFinishedAsExpired(cfmId, vErr.FinishedAt)
		} else {
			panic(vErr)
		}
	case client.ErrorNoAttemptsLeft:
		return nilRes, cfm.NewErrorNoAttemptsLeft(cfmId)
	case client.ErrorRequestsFrequency:
		return nilRes, cfm.NewErrorRequestsFrequency(cfmId, vErr.CanReqAfter, vErr.CanReqAttemptsLeft)
	case std.ErrorRuntime:
		return nilRes, std.WrapErrorToRuntime(vErr, s, "Request", "Internal")
	default:
		panic(clErr)
	}

	// success
	// ----------------

	if err := assert.Str().Eq(cfmId.String()).Check(clRes.CfmId); err != nil {
		return nilRes, std.WrapErrorToRuntime(err, s, "Request", "success", "CfmId")
	}

	// can again
	if clRes.CanReqAgain {
		return cfm.NewServiceResultRequestCanAgain(cfmId, clRes.CanReqAttemptsLeft, clRes.CanReqAfter), nil
	}

	// last
	return cfm.NewServiceResultRequestLast(cfmId), nil
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
func (s *ServiceImplCfmService) Confirm(ctx context.Context, cfmId cfm.Id, code cfm.Code) (cfm.ServiceResultConfirm, error) {
	assert.NotNilDeepMust(ctx)
	assert.FalseMust(cfmId.IsNil())
	assert.FalseMust(code.IsNil())

	nilRes := cfm.ServiceResultConfirmNil

	clRes, clErr := s.client.Confirm(ctx, cfmId.String(), code.String())
	switch vErr := clErr.(type) {
	case nil:
	case std.ErrorNotFound:
		return nilRes, vErr
	case client.ErrorFinished:
		if vErr.IsAsPassed {
			return nilRes, cfm.NewErrorFinishedAsPassed(cfmId, vErr.FinishedAt)
		} else if vErr.IsAsFailed {
			return nilRes, cfm.NewErrorFinishedAsFailed(cfmId, vErr.FinishedAt)
		} else if vErr.IsAsExpired {
			return nilRes, cfm.NewErrorFinishedAsExpired(cfmId, vErr.FinishedAt)
		} else {
			panic(vErr)
		}
	case std.ErrorUnprocessable:
		return nilRes, vErr
	case std.ErrorRuntime:
		return nilRes, std.WrapErrorToRuntime(vErr, s, "Confirm", "Internal")
	default:
		panic(vErr)
	}

	// success
	// ----------------

	// can again
	if !clRes.IsFinishedAsPassed && !clRes.IsFinishedAsFailed {
		return cfm.NewServiceResultConfirmCanAgain(cfmId, clRes.FailsLeft), nil
	}

	// last
	return cfm.NewServiceResultConfirmLast(cfmId, clRes.FinishedAt, clRes.IsFinishedAsPassed), nil
}

// ---------------------------------------------------------------------------------------------------------------------
