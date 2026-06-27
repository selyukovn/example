package middlewares

import (
	"context"
	"example/admin/cfm/internal/api/sch/kernel"
	"runtime/debug"
)

func OnPanic(fnOnPanic func(ctx context.Context, panicValue any, debugStack []byte)) kernel.Middleware {
	return func(next kernel.Handler) kernel.Handler {
		return func(ctx context.Context) {
			defer func() {
				if pv := recover(); pv != nil {
					fnOnPanic(ctx, pv, debug.Stack())
				}
			}()

			next(ctx)
		}
	}
}
