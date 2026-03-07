package handlers

import (
	"bytes"
	"context"
	"example/admin/bff/cmd/http/bundles/auth/config"
	"example/admin/bff/cmd/http/bundles/auth/openapi"
	"example/admin/bff/cmd/http/components/security"
	"html/template"
)

func NewSignInWelcome(
	sec *security.Security,
	cfg *config.Config,
) func(context.Context, openapi.GetAuthSignInWelcomeRequestObject) (openapi.GetAuthSignInWelcomeResponseObject, error) {
	tpl := template.Must(template.ParseFiles(cfg.StaticBasePath() + "/sign_in_welcome.html"))

	return func(ctx context.Context, r openapi.GetAuthSignInWelcomeRequestObject) (openapi.GetAuthSignInWelcomeResponseObject, error) {
		user := sec.AssociatedUser(ctx)
		if user.IsAuthenticated() {
			return openapi.GetAuthSignInWelcome307Response{
				Headers: openapi.GetAuthSignInWelcome307ResponseHeaders{Location: cfg.UrlRedirectToOnSuccess()},
			}, nil
		}

		// --

		bb := bytes.NewBuffer(nil)
		if err := tpl.Execute(bb, struct {
			AppName               string
			UrlSignInRequest      string
			UrlSignInRequestRetry string
			UrlSignInConfirm      string
			StaticBaseUrl         string
		}{
			AppName:               cfg.AppName(),
			UrlSignInRequest:      openapi.UrlSignInRequest,
			UrlSignInRequestRetry: openapi.UrlSignInRequestRetry,
			UrlSignInConfirm:      openapi.UrlSignInConfirm,
			StaticBaseUrl:         cfg.StaticBaseUrl(),
		}); err != nil {
			return nil, err
		}
		return openapi.GetAuthSignInWelcome200TexthtmlResponse{Body: bb}, nil
	}
}
