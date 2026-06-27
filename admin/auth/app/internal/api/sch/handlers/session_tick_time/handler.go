package tick_time

import (
	"context"
	"example/admin/auth/internal/api/sch/kernel"
	"example/admin/auth/internal/opera/use_cases/session_tick_time"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
)

func newSessionTickTime(sessionTickTime session_tick_time.Command) kernel.Handler {
	return func(ctx context.Context) {
		err := sessionTickTime.Execute(ctx, 100)
		switch err.(type) {
		case nil:
		case std.ErrorRuntime:
			err = std.WrapErrorToRuntime(err, "sch.handlers", "SessionTickTime")
			logger.ErrorFf(ctx, err.Error())
		default:
			panic(err)
		}
	}
}
