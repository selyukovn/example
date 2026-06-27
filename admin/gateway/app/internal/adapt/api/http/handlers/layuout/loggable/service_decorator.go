package loggable

import (
	"context"
	"example/admin/gateway/internal/api/http/handlers/layout/openapi"
	"github.com/selyukovn/go-std/logger"
	assert "github.com/selyukovn/go-wm-assert"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

var _ openapi.StrictServerInterface = sDecoratorLoggable{}

type sDecoratorLoggable struct {
	origin openapi.StrictServerInterface
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewDecorator(origin openapi.StrictServerInterface) openapi.StrictServerInterface {
	assert.NotNilDeepMust(origin)

	return sDecoratorLoggable{
		origin: origin,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (s sDecoratorLoggable) GetLayoutUserInfo(ctx context.Context, request openapi.GetLayoutUserInfoRequestObject) (openapi.GetLayoutUserInfoResponseObject, error) {
	logger.InfoFf(ctx, "%T: %T=%+v", s, request, struct{}{})

	resp, err := s.origin.GetLayoutUserInfo(ctx, request)

	switch vResp := resp.(type) {
	case openapi.GetLayoutUserInfo200JSONResponse:
		logger.InfoFf(ctx, "%T: %T=%+v", s, resp, struct {
			Username  string
			AvatarUrl string
		}{
			Username:  *vResp.Username,
			AvatarUrl: *vResp.AvatarUrl,
		})
	default:
		logger.InfoFf(ctx, "%T: %T=%+v", s, resp, resp)
	}

	return resp, err
}

// ---------------------------------------------------------------------------------------------------------------------
