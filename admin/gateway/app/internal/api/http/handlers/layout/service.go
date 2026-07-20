package layout

import (
	"context"
	"example/admin/gateway/internal/api/http/handlers/layout/actions"
	"example/admin/gateway/internal/api/http/handlers/layout/openapi"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ openapi.StrictServerInterface = Service{}

type Service struct {
	userInfo func(context.Context, openapi.GetLayoutUserInfoRequestObject) (openapi.GetLayoutUserInfoResponseObject, error)
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewService(
	userInfo actions.UserInfo,
) Service {
	assert.NotNilDeepMust(userInfo)

	return Service{
		userInfo: userInfo,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (s Service) GetLayoutUserInfo(ctx context.Context, request openapi.GetLayoutUserInfoRequestObject) (openapi.GetLayoutUserInfoResponseObject, error) {
	return s.userInfo(ctx, request)
}

// ---------------------------------------------------------------------------------------------------------------------
