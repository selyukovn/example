package http

import (
	"context"
	"example/admin/gateway/cmd/http/container"
	assert "github.com/selyukovn/go-wm-assert"
	"net"
	"net/http"
)

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Server struct {
	s *http.Server
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

// NewServer
//
// Паникует при нулевых аргументах.
func NewServer(ctr *container.Container) Server {
	assert.NotNilDeepMust(ctr)

	mux := http.NewServeMux()

	// --

	registerRoutes(mux, ctr)

	// --

	s := &http.Server{
		Addr:    ":80",
		Handler: mux,
	}

	return Server{
		s: s,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (s Server) Start(ctx context.Context) error {
	s.s.BaseContext = func(net.Listener) context.Context { return ctx }
	return s.s.ListenAndServe()
}

func (s Server) Stop(ctx context.Context) error {
	// todo : возможно, есть смысл ограничить по времени
	return s.s.Shutdown(ctx)
}

// ---------------------------------------------------------------------------------------------------------------------
