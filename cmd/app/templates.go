package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

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
