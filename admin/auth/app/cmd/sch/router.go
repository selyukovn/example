package sch

import (
	"context"
	"example/admin/auth/cmd/sch/container"
	"example/admin/auth/cmd/sch/handlers"
	"time"
)

func routes(ctr *container.Container) map[string]struct {
	timeout time.Duration
	handler func(context.Context)
} {
	sessionTickTime := handlers.NewSessionTickTime(ctr)

	return map[string]struct {
		timeout time.Duration
		handler func(context.Context)
	}{
		"sessionTickTime": {10 * time.Minute, sessionTickTime},
	}
}
