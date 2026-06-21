package sch

import (
	"context"
	"example/admin/auth/internal/api/sch/kernel"
	"github.com/robfig/cron/v3"
	"github.com/selyukovn/go-std"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Scheduler struct {
	c  *cron.Cron
	wg *sync.WaitGroup
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewScheduler(
	router Router,
	fnInterceptors ...func(string, time.Duration) func(kernel.Handler) kernel.Handler,
) Scheduler {
	c := cron.New()

	for _, name := range router.Routes() {
		timeout, _ := router.RouteTimeout(name)

		handler, _ := router.RouteHandler(name)
		for i := len(fnInterceptors) - 1; i >= 0; i-- {
			handler = fnInterceptors[i](name, timeout)(handler)
		}

		std.Must(c.AddFunc("@every "+timeout.String(), func() {
			ctx := context.Background()
			ctx, cancelCtx := context.WithTimeout(ctx, timeout)
			defer cancelCtx()

			handler(ctx)
		}))
	}

	return Scheduler{
		c:  c,
		wg: new(sync.WaitGroup),
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (s Scheduler) Start(ctx context.Context) error {
	// todo : использовать контекст как базовый
	s.c.Start()
	s.wg.Add(1)
	s.wg.Wait()
	return nil
}

func (s Scheduler) Stop(ctx context.Context) error {
	// todo : возможно, есть смысл ограничить по времени
	s.wg.Done()
	<-s.c.Stop().Done()
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
