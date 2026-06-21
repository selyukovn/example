package root

import (
	"example/admin/front/internal/api/http/kernel"
	"net/http"
)

const UrlRoot = "/"

func Register(mux *http.ServeMux, urlForGuest string, urlForAuthorized string) {
	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		if kernel.CookieHasSessId(r) {
			kernel.Redirect307(w, r, urlForAuthorized)
		} else {
			kernel.Redirect307(w, r, urlForGuest)
		}
	})
}
