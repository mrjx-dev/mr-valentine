package main

import (
	"log"
	"net/http"
	"os"
)

// App represents the main application
type Application struct {
	cfg        *Config
	server     *http.Server
	middleware *Middleware
	templates  *TemplateHandler
	browser    *BrowserOpener
	shutdown   chan bool
	done       chan bool
	cancel     chan bool
}

// NewApp creates a new application instance
func NewApp() (*Application, error) {
	cfg := NewConfig()

	// Initialize template handler
	templates, err := NewTemplateHandler(cfg)
	if err != nil {
		return nil, err
	}

	return &Application{
		cfg:       cfg,
		templates: templates,
		browser:   NewBrowserOpener(cfg),
		shutdown:  make(chan bool),
		done:      make(chan bool),
		cancel:    make(chan bool),
	}, nil
}

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
