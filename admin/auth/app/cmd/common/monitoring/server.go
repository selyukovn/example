package monitoring

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func NewMonitoringServer() Server {
	mux := http.NewServeMux()

	mux.Handle("/prometheus-metrics", promhttp.Handler())

	s := &http.Server{
		Addr:    ":8080",
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
