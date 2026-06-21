package handlers

import (
	"context"
	"example/admin/gateway/internal/api/http/bundles/auth/openapi"
	"example/admin/gateway/internal/api/http/components/security"
	"github.com/selyukovn/go-std"
	"net/http"
)

func NewSignOut(
	sec security.Security,
) func(context.Context, openapi.DeleteAuthSignOutRequestObject) (openapi.DeleteAuthSignOutResponseObject, error) {
	return func(ctx context.Context, r openapi.DeleteAuthSignOutRequestObject) (openapi.DeleteAuthSignOutResponseObject, error) {
		_o_ := "actions"
		_m_ := "NewSignOut"

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
			return nil, std.WrapErrorToRuntime(err, _o_, _m_, "UnAuthenticate")
		default:
			panic(err)
		}

		// --

		xTrue := true
		return openapi.DeleteAuthSignOut200JSONResponse{Success: &xTrue}, nil
	}
}
