package handlers

import (
	"context"
	"example/admin/gateway/cmd/http/bundles/layout/openapi"
	"example/admin/gateway/cmd/http/components/security"
	"example/admin/gateway/cmd/http/container"
	"net/http"
)

func NewUserInfo(
	ctr *container.Container,
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
