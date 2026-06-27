package kernel

import (
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func RootMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := strings.Replace(uuid.Must(uuid.NewRandom()).String(), "-", "", -1)
			enrichRequest(r, requestId)

			w = wrapResponseWriter(w)

			next.ServeHTTP(w, r)
		})
	}
}
