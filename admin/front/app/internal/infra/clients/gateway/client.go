package gateway

import (
	"context"
	"example/admin/front/internal/infra/clients/gateway/auth"
	"fmt"
	"github.com/selyukovn/go-std"
	assert "github.com/selyukovn/go-wm-assert"
	"net/http"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type ApiClient struct {
	auth auth.ClientWithResponsesInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewApiClient(baseUrl string) ApiClient {
	authClient, err := auth.NewClientWithResponses(baseUrl)
	assert.TrueMust(err == nil)

	return ApiClient{
		auth: authClient,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (c ApiClient) fnAddHeaders(
	fromIp string,
	fromUserAgent string,
	sessionId string,
) func(ctx context.Context, req *http.Request) error {
	return func(ctx context.Context, r *http.Request) error {
		r.Header.Set("X-Client-Ip", fromIp)
		r.Header.Set("User-Agent", fromUserAgent)

		if sessionId != "" {
			r.Header.Set("Authorization", "Bearer "+sessionId)
		}

		return nil
	}
}

// Auth
// ---------------------------------------------------------------------------------------------------------------------

func (c ApiClient) AuthSignInRequest(
	fromIp string,
	fromUserAgent string,
	email string,
) (*auth.PostAuthSignInRequestResponse, error) {
	resp, err := c.auth.PostAuthSignInRequestWithResponse(context.Background(), auth.PostAuthSignInRequestJSONRequestBody{
		Email: &email,
	}, c.fnAddHeaders(fromIp, fromUserAgent, ""))
	if err != nil {
		return nil, err
	} else if sCode := resp.StatusCode(); sCode != http.StatusOK && sCode != http.StatusUnprocessableEntity {
		return nil, std.WrapErrorToRuntime(fmt.Errorf("%d : %s", sCode, resp.Status()), c, "AuthSignInRequest")
	}
	return resp, nil
}

func (c ApiClient) AuthSignInRequestRetry(
	fromIp string,
	fromUserAgent string,
	signInId string,
) (*auth.PutAuthSignInRequestRetryResponse, error) {
	resp, err := c.auth.PutAuthSignInRequestRetryWithResponse(context.Background(), auth.PutAuthSignInRequestRetryJSONRequestBody{
		SignInId: &signInId,
	}, c.fnAddHeaders(fromIp, fromUserAgent, ""))
	if err != nil {
		return nil, err
	} else if sCode := resp.StatusCode(); sCode != http.StatusOK && sCode != http.StatusUnprocessableEntity {
		return nil, std.WrapErrorToRuntime(fmt.Errorf("%d : %s", sCode, resp.Status()), c, "AuthSignInRequestRetry")
	}
	return resp, nil
}

func (c ApiClient) AuthSignInConfirm(
	fromIp string,
	fromUserAgent string,
	signInId string,
	code string,
) (*auth.PutAuthSignInConfirmResponse, error) {
	resp, err := c.auth.PutAuthSignInConfirmWithResponse(context.Background(), auth.PutAuthSignInConfirmJSONRequestBody{
		SignInId: &signInId,
		Code:     &code,
	}, c.fnAddHeaders(fromIp, fromUserAgent, ""))
	if err != nil {
		return nil, err
	} else if sCode := resp.StatusCode(); sCode != http.StatusOK && sCode != http.StatusUnprocessableEntity {
		return nil, std.WrapErrorToRuntime(fmt.Errorf("%d : %s", sCode, resp.Status()), c, "AuthSignInConfirm")
	}
	return resp, nil
}

func (c ApiClient) AuthSignOut(
	fromIp string,
	fromUserAgent string,
	sessId string,
) (*auth.DeleteAuthSignOutResponse, error) {
	resp, err := c.auth.DeleteAuthSignOutWithResponse(context.Background(), c.fnAddHeaders(fromIp, fromUserAgent, sessId))
	if err != nil {
		return nil, err
	} else if sCode := resp.StatusCode(); sCode != http.StatusOK && sCode != http.StatusUnprocessableEntity {
		return nil, std.WrapErrorToRuntime(fmt.Errorf("%d : %s", sCode, resp.Status()), c, "AuthSignOut")
	}
	return resp, nil
}

// ---------------------------------------------------------------------------------------------------------------------
