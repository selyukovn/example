package loggable

import (
	"context"
	"example/admin/gateway/internal/api/http/handlers/auth/openapi"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ openapi.StrictServerInterface = Decorator{}

type Decorator struct {
	origin openapi.StrictServerInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewDecorator(origin openapi.StrictServerInterface) Decorator {
	assert.NotNilDeepMust(origin)

	return Decorator{
		origin: origin,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (d Decorator) PostAuthSignInRequest(ctx context.Context, request openapi.PostAuthSignInRequestRequestObject) (openapi.PostAuthSignInRequestResponseObject, error) {
	logger.InfoFf(ctx, "%T: %+v", request, struct {
		Email string
	}{
		Email: *request.Body.Email,
	})

	resp, err := d.origin.PostAuthSignInRequest(ctx, request)

	switch vResp := resp.(type) {
	case openapi.PostAuthSignInRequest200JSONResponse:
		logger.InfoFf(ctx, "%T: %+v", resp, struct {
			CanRetryAt  string
			ExpireAt    string
			RetriesLeft int
			SignInId    string
		}{
			CanRetryAt:  *vResp.CanRetryAt,
			ExpireAt:    *vResp.ExpireAt,
			RetriesLeft: *vResp.RetriesLeft,
			SignInId:    *vResp.SignInId,
		})
	default:
		logger.InfoFf(ctx, "%T: %+v", resp, resp)
	}

	return resp, err
}

func (d Decorator) PutAuthSignInRequestRetry(ctx context.Context, request openapi.PutAuthSignInRequestRetryRequestObject) (openapi.PutAuthSignInRequestRetryResponseObject, error) {
	logger.InfoFf(ctx, "%T: %+v", request, struct {
		SignInId string
	}{
		SignInId: *request.Body.SignInId,
	})

	resp, err := d.origin.PutAuthSignInRequestRetry(ctx, request)

	switch vResp := resp.(type) {
	case openapi.PutAuthSignInRequestRetry200JSONResponse:
		logger.InfoFf(ctx, "%T: %+v", resp, struct {
			CanRetryAt  string
			RetriesLeft int
			SignInId    string
		}{
			CanRetryAt:  *vResp.CanRetryAt,
			RetriesLeft: *vResp.RetriesLeft,
			SignInId:    *vResp.SignInId,
		})
	default:
		logger.InfoFf(ctx, "%T: %+v", resp, resp)
	}

	return resp, err
}

func (d Decorator) PutAuthSignInConfirm(ctx context.Context, request openapi.PutAuthSignInConfirmRequestObject) (openapi.PutAuthSignInConfirmResponseObject, error) {
	logger.InfoFf(ctx, "%T: %+v", request, struct {
		SignInId string
		Code     string
	}{
		SignInId: *request.Body.SignInId,
		Code:     std.MaskStrNotFirstLast(*request.Body.Code),
	})

	resp, err := d.origin.PutAuthSignInConfirm(ctx, request)

	switch vResp := resp.(type) {
	case openapi.PutAuthSignInConfirm200JSONResponse:
		logger.InfoFf(ctx, "%T: %+v", resp, struct {
			AttemptsLeft    int
			IsPassed        bool
			SessionId       string
			SessionExpireAt string
		}{
			AttemptsLeft:    *vResp.AttemptsLeft,
			IsPassed:        *vResp.IsPassed,
			SessionId:       std.MaskStrNotFirstLast(*vResp.SessionId),
			SessionExpireAt: *vResp.SessionExpireAt,
		})
	default:
		logger.InfoFf(ctx, "%T: %+v", resp, resp)
	}

	return resp, err
}

func (d Decorator) DeleteAuthSignOut(ctx context.Context, request openapi.DeleteAuthSignOutRequestObject) (openapi.DeleteAuthSignOutResponseObject, error) {
	logger.InfoFf(ctx, "%T: %+v", request, request)

	resp, err := d.origin.DeleteAuthSignOut(ctx, request)

	switch vResp := resp.(type) {
	case openapi.DeleteAuthSignOut200JSONResponse:
		logger.InfoFf(ctx, "%T: %+v", resp, struct {
			Success bool
		}{
			Success: *vResp.Success,
		})
	default:
		logger.InfoFf(ctx, "%T: %+v", resp, resp)
	}

	return resp, err
}

// ---------------------------------------------------------------------------------------------------------------------
