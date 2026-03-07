package root

import (
	"example/admin/bff/cmd/http/bundles/root/config"
	"example/admin/bff/cmd/http/components/security"
	"net/http"
)

// #####################################################################################################################
// CONSTANTS
// #####################################################################################################################

const (
	UrlRoot = "/"
)

// #####################################################################################################################
// ROUTER
// #####################################################################################################################

func Register(
	mux *http.ServeMux,
	middlewares []func(http.Handler) http.Handler,
	sec *security.Security,
	urlForGuest string,
	urlForAuthorized string,
) {
	cfg := config.New(
		urlForGuest,
		urlForAuthorized,
	)

	rootHandler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sec.AssociatedUser(r.Context()).IsGuest() {
			http.Redirect(w, r, cfg.UrlForGuest(), http.StatusTemporaryRedirect)
		} else {
			http.Redirect(w, r, cfg.UrlForAuthorized(), http.StatusTemporaryRedirect)
		}
	}))

	unknownRouteHandler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))

	for i := len(middlewares) - 1; i >= 0; i-- {
		rootHandler = middlewares[i](rootHandler)
		unknownRouteHandler = middlewares[i](unknownRouteHandler)
	}

	// --

	mux.Handle("GET /{$}", rootHandler)
	mux.Handle("/", unknownRouteHandler)
}

// #####################################################################################################################
