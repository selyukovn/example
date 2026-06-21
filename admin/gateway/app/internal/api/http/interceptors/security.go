package interceptors

import (
	"example/admin/gateway/internal/api/http/components/security"
	"github.com/selyukovn/go-std/logger"
	"net/http"
)

func Security(sec security.Security) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return sec.Middleware(
			// Имеет смысл искать логи по аккаунту,
			// но обогащение контекста для логгера не входит в обязанности security-middleware --
			// нужна еще одна обертка.
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()

				user := sec.AssociatedUser(ctx)
				if user.IsAuthenticated() {
					ctx = logger.AddAttrToCtx(ctx, "account_id", user.AccountId())
				} else {
					// todo : аналогично сделать для гостей с каким-то guest-id ???
				}

				r = r.WithContext(ctx)

				next.ServeHTTP(w, r)
			}),
		)
	}
}
