package sch

import (
	"example/admin/cfm/internal/api/sch/kernel"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Router struct {
	m map[string]struct {
		timeout time.Duration
		handler kernel.Handler
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

func NewRouter() Router {
	return Router{
		m: make(map[string]struct {
			timeout time.Duration
			handler kernel.Handler
		}),
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (r Router) Register(name string, timeout time.Duration, handler kernel.Handler) {
	r.m[name] = struct {
		timeout time.Duration
		handler kernel.Handler
	}{
		timeout: timeout,
		handler: handler,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (r Router) Routes() []string {
	s := make([]string, 0, len(r.m))
	for k := range r.m {
		s = append(s, k)
	}
	return s
}

func (r Router) RouteTimeout(route string) (time.Duration, bool) {
	v, ok := r.m[route]
	if !ok {
		return 0, false
	}
	return v.timeout, true
}

func (r Router) RouteHandler(route string) (kernel.Handler, bool) {
	v, ok := r.m[route]
	if !ok {
		return nil, false
	}
	return v.handler, true
}

// ---------------------------------------------------------------------------------------------------------------------
