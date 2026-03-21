package handlers

import (
	"context"
	"example/admin/cfm/cmd/common/container"
	"example/admin/cfm/internal/opera/use_cases/tick_time"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
)

func NewTickTime(ctr *container.Container) func(ctx context.Context) {
	return func(ctx context.Context) {
		err := ctr.UseCases.TickTime.Execute(tick_time.NewArgs(ctx, 100))
		switch err.(type) {
		case nil:
		case std.ErrorRuntime:
			err = std.WrapErrorToRuntime(err, "sch.handlers", "TickTime")
			logger.ErrorFf(ctx, err.Error())
		default:
			panic(err)
		}
	}
}
