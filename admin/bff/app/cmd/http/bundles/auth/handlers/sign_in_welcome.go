package handlers

import (
	"example/admin/bff/cmd/http/bundles/auth/config"
	"example/admin/bff/cmd/http/components/security"
	"example/admin/bff/cmd/http/container"
	"example/admin/bff/cmd/http/kernel"
	"html/template"
	"net/http"
)

func NewSignInWelcome(ctr *container.Container, cfg *config.Config, sec *security.Security) http.Handler {
	tpl := template.Must(template.ParseFiles(cfg.StaticBasePath() + "/sign_in_welcome.html"))

	return sec.AllowAny(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := sec.AssociatedUser(r)

		if user.IsAuthorized() {
			kernel.Redirect307(w, r, cfg.UrlRedirectToOnSuccess())
			return
		}

		// --

		ctx := r.Context()

		// --

		if err := tpl.Execute(w, struct {
			AppName               string
			UrlSignInRequest      string
			UrlSignInRequestRetry string
			UrlSignInConfirm      string
			StaticBaseUrl         string
		}{
			AppName:               cfg.AppName(),
			UrlSignInRequest:      cfg.UrlSignInRequest(),
			UrlSignInRequestRetry: cfg.UrlSignInRequestRetry(),
			UrlSignInConfirm:      cfg.UrlSignInConfirm(),
			StaticBaseUrl:         cfg.StaticBaseUrl(),
		}); err != nil {
			ctr.Logger.CtxErrorFf(ctx, err.Error())
			kernel.Error500(w)
		}
	}))
}
