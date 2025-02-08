package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// Middleware represents custom middleware functions
type Middleware struct {
	cfg *Config
}

// NewMiddleware creates a new Middleware instance
func NewMiddleware(cfg *Config) *Middleware {
	return &Middleware{cfg: cfg}
}

// SetupCommonMiddleware returns a slice of middleware handlers
func (m *Middleware) SetupCommonMiddleware() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		middleware.Logger,    // Log API request calls
		middleware.Recoverer, // Recover from panics without crashing server
		middleware.RequestID, // Injects a request ID into the context of each request
		middleware.RealIP,    // Sets a http.Request's RemoteAddr to either X-Real-IP or X-Forwarded-For
		middleware.Timeout(
			60 * time.Second,
		), // Timeout requests after 60 seconds
	}
}
