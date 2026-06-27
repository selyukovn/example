package auth

import (
	"context"
	"example/admin/gateway/internal/api/http/handlers/auth/actions"
	"example/admin/gateway/internal/api/http/handlers/auth/openapi"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ openapi.StrictServerInterface = sDefault{}

type sDefault struct {
	signInRequest      actions.SignInRequest
	signInRequestRetry actions.SignInRequestRetry
	signInConfirm      actions.SignInConfirm
	signOut            actions.SignOut
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func newServiceDefault(
	signInRequest actions.SignInRequest,
	signInRequestRetry actions.SignInRequestRetry,
	signInConfirm actions.SignInConfirm,
	signOut actions.SignOut,
) sDefault {
	assert.NotNilDeepMust(signInRequest)
	assert.NotNilDeepMust(signInRequestRetry)
	assert.NotNilDeepMust(signInConfirm)
	assert.NotNilDeepMust(signOut)

	return sDefault{
		signInRequest:      signInRequest,
		signInRequestRetry: signInRequestRetry,
		signInConfirm:      signInConfirm,
		signOut:            signOut,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (r sDefault) PostAuthSignInRequest(ctx context.Context, request openapi.PostAuthSignInRequestRequestObject) (openapi.PostAuthSignInRequestResponseObject, error) {
	return r.signInRequest(ctx, request)
}

func (r sDefault) PutAuthSignInRequestRetry(ctx context.Context, request openapi.PutAuthSignInRequestRetryRequestObject) (openapi.PutAuthSignInRequestRetryResponseObject, error) {
	return r.signInRequestRetry(ctx, request)
}

func (r sDefault) PutAuthSignInConfirm(ctx context.Context, request openapi.PutAuthSignInConfirmRequestObject) (openapi.PutAuthSignInConfirmResponseObject, error) {
	return r.signInConfirm(ctx, request)
}

func (r sDefault) DeleteAuthSignOut(ctx context.Context, request openapi.DeleteAuthSignOutRequestObject) (openapi.DeleteAuthSignOutResponseObject, error) {
	return r.signOut(ctx, request)
}

// ---------------------------------------------------------------------------------------------------------------------
