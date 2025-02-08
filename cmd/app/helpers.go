package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// handleShutdown handles the shutdown request
func (app *Application) handleShutdown() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ShutdownRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			req.Immediate = false
		}

		w.WriteHeader(http.StatusOK)

		if req.Immediate {
			app.shutdown <- true
			return
		}

		go func() {
			select {
			case <-time.After(app.cfg.Server.GracePeriod):
				app.shutdown <- true
			case <-app.cancel:
				log.Println("Shutdown cancelled")
			}
		}()
	}
}

type ShutdownRequest struct {
	Immediate bool `json:"immediate"`
}

// handleCancelShutdown handles the shutdown cancellation request
func (app *Application) handleCancelShutdown() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.cancel <- true
		w.WriteHeader(http.StatusOK)
	}
}

// Start initializes and starts the application
func (app *Application) Start() error {
	app.middleware = NewMiddleware(app.cfg)

	// Initialize server
	app.server = &http.Server{
		Addr:    ":" + app.cfg.Server.Port,
		Handler: app.setupRoutes(),
	}

	// Start server
	go func() {
		log.Printf("Server starting on %s", app.cfg.BaseURL())
		if err := app.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Open browser
	if err := app.browser.Open(app.cfg.BaseURL()); err != nil {
		log.Printf("Warning: failed to open browser: %v", err)
	}

	// Handle shutdown
	go app.handleGracefulShutdown()

	// Wait for shutdown to complete
	<-app.done
	log.Println("Application exiting")
	return nil
}

// handleGracefulShutdown manages the graceful shutdown process
func (app *Application) handleGracefulShutdown() {
	<-app.shutdown
	log.Println("Initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		app.cfg.Server.ShutdownTimeout,
	)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("Server shutdown complete")
	app.done <- true
}

// TemplateHandler manages template rendering
type TemplateHandler struct {
	cfg       *Config
	templates map[string]*template.Template
}

// NewTemplateHandler creates a new template handler instance
func NewTemplateHandler(cfg *Config) (*TemplateHandler, error) {
	th := &TemplateHandler{
		cfg:       cfg,
		templates: make(map[string]*template.Template),
	}

	// Load templates
	if err := th.loadTemplates(); err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	return th, nil
}

// loadTemplates loads all templates from the template directory
func (th *TemplateHandler) loadTemplates() error {
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Load index template
	indexPath := filepath.Join(workDir, "ui/templates/index.html")
	tmpl, err := template.ParseFiles(indexPath)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", indexPath, err)
	}

	th.templates["index"] = tmpl
	return nil
}

// RenderTemplate renders a template by name with the given data
func (th *TemplateHandler) RenderTemplate(
	w http.ResponseWriter,
	name string,
	data interface{},
) error {
	tmpl, ok := th.templates[name]
	if !ok {
		return fmt.Errorf("template %s not found", name)
	}

	return tmpl.Execute(w, data)
}

// HandleIndex returns a handler for the index page
func (th *TemplateHandler) HandleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := th.RenderTemplate(w, "index", nil); err != nil {
			log.Printf("Error rendering index template: %v", err)
			http.Error(
				w,
				"Internal Server Error",
				http.StatusInternalServerError,
			)
		}
	}
}

// BrowserOpener handles browser operations
type BrowserOpener struct {
	cfg *Config
}

// NewBrowserOpener creates a new browser opener instance
func NewBrowserOpener(cfg *Config) *BrowserOpener {
	return &BrowserOpener{cfg: cfg}
}

// Open opens the default browser to the specified URL
func (bo *BrowserOpener) Open(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}

	return nil
}
