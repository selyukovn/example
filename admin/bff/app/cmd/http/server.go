package http

import (
	"context"
	"example/admin/bff/cmd/http/container"
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
	ctr *container.Container,
	appName string,
	baseUrl string,
	sessionCookieName string,
) *Server {
	assert.NotNilDeepMust(ctr)
	assert.Str().NotEmpty().Must(appName)
	assert.Str().NotEmpty().Must(baseUrl)
	assert.Str().NotEmpty().Must(sessionCookieName)

	mux := http.NewServeMux()

	// --

	registerRoutes(
		mux,
		ctr,
		appName,
		baseUrl,
		sessionCookieName,
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
