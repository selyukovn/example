package sch

import (
	"context"
	"example/admin/auth/cmd/sch/container"
	"example/admin/auth/cmd/sch/interceptors"
	"github.com/robfig/cron/v3"
	assert "github.com/selyukovn/go-wm-assert"
	"sync"
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

// NewScheduler
//
// Паникует при нулевых аргументах.
// Паникует при ошибке регистрации задачи.
func NewScheduler(ctr *container.Container) Scheduler {
	assert.NotNilDeepMust(ctr)

	c := cron.New()

	fnIntercepts := []func(func(context.Context), string, *container.Container) func(context.Context){
		interceptors.NewBoundary(),
	}

	for jName, jDef := range routes(ctr) {
		handler := jDef.handler
		for i := len(fnIntercepts) - 1; i >= 0; i-- {
			handler = fnIntercepts[i](handler, jName, ctr)
		}

		timeout := jDef.timeout
		cronEntryId, err := c.AddFunc("@every "+timeout.String(), func() {
			ctx := context.Background()
			ctx, cancelCtx := context.WithTimeout(ctx, timeout)
			defer cancelCtx()

			handler(ctx)
		})
		if err != nil {
			panic(err)
		}

		_ = cronEntryId
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
