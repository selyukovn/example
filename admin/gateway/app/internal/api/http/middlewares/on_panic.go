package middlewares

import (
	"net/http"
	"runtime/debug"
)

func OnPanic(
	fnOnPanic func(w http.ResponseWriter, r *http.Request, panicValue any, debugStack []byte),
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if pv := recover(); pv != nil {
					fnOnPanic(w, r, pv, debug.Stack())
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
