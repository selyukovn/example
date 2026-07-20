package tick_time

import (
	"context"
	"example/admin/cfm/internal/opera/use_cases/tick_time"
	"github.com/selyukovn/go-std"
	"github.com/selyukovn/go-std/logger"
)

// ---------------------------------------------------------------------------------------------------------------------

type TickTime = func(ctx context.Context)

// ---------------------------------------------------------------------------------------------------------------------

func newHandler(ucTickTime tick_time.Command) TickTime {
	return func(ctx context.Context) {
		err := ucTickTime.Execute(ctx, 100)
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

// ---------------------------------------------------------------------------------------------------------------------
