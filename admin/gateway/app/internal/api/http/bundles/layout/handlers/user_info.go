package handlers

import (
	"context"
	"example/admin/gateway/internal/api/http/bundles/layout/openapi"
	"example/admin/gateway/internal/api/http/components/security"
	"net/http"
)

func NewUserInfo(
	sec security.Security,
) func(context.Context, openapi.GetLayoutUserInfoRequestObject) (openapi.GetLayoutUserInfoResponseObject, error) {
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
