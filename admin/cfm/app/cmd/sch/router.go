package sch

import (
	"context"
	"example/admin/cfm/cmd/common/container"
	"example/admin/cfm/cmd/sch/handlers"
	"time"
)

func routes(ctr *container.Container) map[string]struct {
	timeout time.Duration
	handler func(context.Context)
} {
	tickTime := handlers.NewTickTime(ctr)

	return map[string]struct {
		timeout time.Duration
		handler func(context.Context)
	}{
		"tickTime": {10 * time.Minute, tickTime},
	}
}
