package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

// setupRoutes configures all application routes
func (app *Application) setupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Add middleware
	for _, mw := range app.middleware.SetupCommonMiddleware() {
		r.Use(mw)
	}

	// Static files
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "ui/static"))
	r.Handle(
		"/static/*",
		http.StripPrefix("/static/", http.FileServer(filesDir)),
	)

	// Routes
	r.Get("/", app.templates.HandleIndex())
	r.Post("/trigger-shutdown", app.handleShutdown())
	r.Post("/cancel-shutdown", app.handleCancelShutdown())

	return r
}
