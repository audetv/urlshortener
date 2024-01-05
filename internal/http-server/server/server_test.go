package server_test

import (
	"github.com/audetv/urlshortener/internal/config"
	"github.com/audetv/urlshortener/internal/http-server/server"
	"net/http"
	"testing"
	"time"
)

// TestNew tests the creation of a new server instance.
// TestNew created with codeium
func TestNew(t *testing.T) {
	// Create a new configuration with test environment settings.
	cfg := &config.Config{
		Env: "test",
		HTTPServer: config.HTTPServer{
			Address:     "localhost:8082",
			Timeout:     10 * time.Second,
			IdleTimeout: 5 * time.Minute,
		},
	}

	// Create a new router instance.
	router := http.NewServeMux()

	// Create a new server instance with the provided configuration and router.
	srv := server.New(cfg, router)

	// Check if the server address is set correctly.
	if srv.Addr != cfg.Address {
		t.Errorf("Expected server address to be %s, but got %s", cfg.Address, srv.Addr)
	}

	// Check if the server handler is set correctly.
	if srv.Handler != router {
		t.Error("Expected server handler to be the provided router")
	}

	// Check if the server read timeout is set correctly.
	if srv.ReadTimeout != cfg.HTTPServer.Timeout {
		t.Errorf("Expected read timeout to be %s, but got %s", cfg.HTTPServer.Timeout, srv.ReadTimeout)
	}

	// Check if the server write timeout is set correctly.
	if srv.WriteTimeout != cfg.HTTPServer.Timeout {
		t.Errorf("Expected write timeout to be %s, but got %s", cfg.HTTPServer.Timeout, srv.WriteTimeout)
	}

	// Check if the server idle timeout is set correctly.
	if srv.IdleTimeout != cfg.HTTPServer.IdleTimeout {
		t.Errorf("Expected idle timeout to be %s, but got %s", cfg.HTTPServer.IdleTimeout, srv.IdleTimeout)
	}
}
