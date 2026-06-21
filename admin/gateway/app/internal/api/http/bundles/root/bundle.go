package root

import "net/http"

func Register(mux *http.ServeMux) {
	mux.Handle("/", http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})))
}
