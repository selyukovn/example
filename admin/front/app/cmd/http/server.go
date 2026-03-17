package http

import (
	"context"
	"example/admin/front/cmd/http/kernel"
	"example/admin/front/internal/infra/clients/gateway"
	"example/admin/front/internal/infra/logger"
	assert "github.com/selyukovn/go-wm-assert"
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
	logger *logger.Logger,
	apiClient *gateway.ApiClient,
	appName string,
	baseUrl string,
	sessionCookieName string,
) *Server {
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
		logger,
		apiClient,
		mux,
		appName,
	)

	// --

	s := &http.Server{
		Addr:    ":80",
		Handler: mux,
	}

	return &Server{
		s: s,
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Actions
// ---------------------------------------------------------------------------------------------------------------------

func (s *Server) Start() error {
	return s.s.ListenAndServe()
}

func (s *Server) Stop() error {
	// todo : возможно, есть смысл ограничить по времени
	return s.s.Shutdown(context.Background())
}

// ---------------------------------------------------------------------------------------------------------------------
