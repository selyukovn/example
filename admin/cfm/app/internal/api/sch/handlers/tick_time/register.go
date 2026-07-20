package tick_time

import (
	"context"
	"example/admin/cfm/internal/api/sch/kernel"
	"example/admin/cfm/internal/opera/use_cases/tick_time"
	"github.com/robfig/cron/v3"
	"github.com/selyukovn/go-std"
	"time"
)

func Register(
	c *cron.Cron,
	timeout time.Duration,
	ucTickTime tick_time.Command,
	hMws []kernel.Middleware,
) {
	handler := newHandler(ucTickTime)

	for i := len(hMws) - 1; i >= 0; i-- {
		handler = hMws[i](handler)
	}

	std.Must(c.AddFunc("@every "+timeout.String(), func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer func() { cancel() }()

		handler(ctx)
	}))
}
