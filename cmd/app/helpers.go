package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

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
