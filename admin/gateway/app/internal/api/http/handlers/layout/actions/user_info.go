package actions

import (
	"context"
	"example/admin/gateway/internal/api/http/components/security"
	"example/admin/gateway/internal/api/http/handlers/layout/openapi"
	"net/http"
)

// ---------------------------------------------------------------------------------------------------------------------

type UserInfo = func(context.Context, openapi.GetLayoutUserInfoRequestObject) (openapi.GetLayoutUserInfoResponseObject, error)

// ---------------------------------------------------------------------------------------------------------------------

func NewUserInfo(sec security.Security) UserInfo {
	return func(ctx context.Context, r openapi.GetLayoutUserInfoRequestObject) (openapi.GetLayoutUserInfoResponseObject, error) {
		user := sec.AssociatedUser(ctx)
		if user.IsGuest() {
			return openapi.GetLayoutUserInfo422JSONResponse{
				Code:    http.StatusUnauthorized,
				Message: http.StatusText(http.StatusUnauthorized),
			}, nil
		}

		// TODO : ...
		avatarUrl := ""
		username := user.AccountId()

		return openapi.GetLayoutUserInfo200JSONResponse{
			AvatarUrl: &avatarUrl,
			Username:  &username,
		}, nil
	}
}

// ---------------------------------------------------------------------------------------------------------------------
