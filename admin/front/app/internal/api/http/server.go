package http

import (
	"context"
	"example/admin/front/internal/api/http/kernel"
	"example/admin/front/internal/infra/clients/gateway"
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
func NewServer(
	apiClient gateway.ApiClient,
	appName string,
	baseUrl string,
	sessionCookieName string,
) Server {
	assert.Str().NotEmpty().Must(appName)
	assert.Str().NotEmpty().Must(baseUrl)
	assert.Str().NotEmpty().Must(sessionCookieName)

	mux := http.NewServeMux()

	// --

	kernel.Configure(
		baseUrl,
		sessionCookieName,
	)

	registerRoutes(
		apiClient,
		mux,
		appName,
	)

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
