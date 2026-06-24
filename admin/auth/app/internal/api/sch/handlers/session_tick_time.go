package handlers

import (
	"context"
	"example/admin/auth/internal/opera/use_cases/session_tick_time"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
)

func NewSessionTickTime(sessionTickTime session_tick_time.Command) func(ctx context.Context) {
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
