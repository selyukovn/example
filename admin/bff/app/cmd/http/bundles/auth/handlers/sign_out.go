package handlers

import (
	"context"
	"example/admin/bff/cmd/http/bundles/auth/openapi"
	"example/admin/bff/cmd/http/components/security"
	"github.com/selyukovn/go-std"
	"net/http"
)

func NewSignOut(
	sec *security.Security,
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

		redirectUrl := openapi.UrlSignInWelcome
		return signOutSuccessResponse{
			openapi.DeleteAuthSignOut200JSONResponse{
				RedirectUrl: &redirectUrl,
			},
			func(w http.ResponseWriter) error {
				err := sec.UnAuthenticate(ctx, user, w)
				switch err.(type) {
				case nil:
				case std.ErrorRuntime:
					return std.WrapErrorToRuntime(err, "handlers", "NewSignOut", "UnAuthenticate")
				default:
					panic(err)
				}
				return nil
			},
		}, nil
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type signOutSuccessResponse struct {
	openapi.DeleteAuthSignOutResponseObject
	fnPreResponse func(http.ResponseWriter) error
}

func (r signOutSuccessResponse) VisitDeleteAuthSignOutResponse(w http.ResponseWriter) error {
	if err := r.fnPreResponse(w); err != nil {
		return err
	}
	return r.DeleteAuthSignOutResponseObject.VisitDeleteAuthSignOutResponse(w)
}

// ---------------------------------------------------------------------------------------------------------------------
