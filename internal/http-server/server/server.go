package server

import (
	"github.com/audetv/urlshortener/internal/config"
	"net/http"
)

// New creates a new HTTP server with the provided configuration and router.
// It returns a pointer to the created server.
func New(cfg *config.Config, router http.Handler) *http.Server {
	// Create a new server object with the provided configuration
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	return srv
}
