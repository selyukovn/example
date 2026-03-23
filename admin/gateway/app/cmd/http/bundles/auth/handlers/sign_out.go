package handlers

import (
	"context"
	"example/admin/gateway/cmd/http/bundles/auth/openapi"
	"example/admin/gateway/cmd/http/components/security"
	"github.com/selyukovn/go-std"
	"net/http"
)

func NewSignOut(
	sec security.Security,
) func(context.Context, openapi.DeleteAuthSignOutRequestObject) (openapi.DeleteAuthSignOutResponseObject, error) {
	return func(ctx context.Context, r openapi.DeleteAuthSignOutRequestObject) (openapi.DeleteAuthSignOutResponseObject, error) {
		user := sec.AssociatedUser(ctx)
		if user.IsGuest() {
			return openapi.DeleteAuthSignOut422JSONResponse{
				Code:    http.StatusUnauthorized,
				Message: http.StatusText(http.StatusUnauthorized),
			}, nil
		}

		// --

		err := sec.UnAuthenticate(ctx, user)
		switch err.(type) {
		case nil:
		case std.ErrorRuntime:
			return nil, std.WrapErrorToRuntime(err, "handlers", "NewSignOut", "UnAuthenticate")
		default:
			panic(err)
		}

		// --

		xTrue := true
		return openapi.DeleteAuthSignOut200JSONResponse{Success: &xTrue}, nil
	}
}
