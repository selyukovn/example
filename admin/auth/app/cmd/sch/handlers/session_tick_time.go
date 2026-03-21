package handlers

import (
	"context"
	"example/admin/auth/cmd/sch/container"
	"example/admin/auth/internal/opera/use_cases/session_tick_time"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
)

func NewSessionTickTime(ctr *container.Container) func(ctx context.Context) {
	return func(ctx context.Context) {
		err := ctr.UseCases.SessionTickTime.Execute(session_tick_time.NewArgs(ctx, 100))
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
