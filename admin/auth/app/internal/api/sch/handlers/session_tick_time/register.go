package tick_time

import (
	"context"
	"example/admin/auth/internal/api/sch/kernel"
	"example/admin/auth/internal/opera/use_cases/session_tick_time"
	"github.com/robfig/cron/v3"
	"github.com/selyukovn/go-std"
	"time"
)

func Register(
	c *cron.Cron,
	timeout time.Duration,
	ucSessionTickTime session_tick_time.Command,
	hMws []kernel.Middleware,
) {
	handler := newSessionTickTime(ucSessionTickTime)

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
