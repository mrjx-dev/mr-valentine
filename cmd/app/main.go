package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ShutdownRequest struct {
	Immediate bool `json:"immediate"`
}

func main() {
	port := "8080"
	url := fmt.Sprintf("http://localhost:%s", port)

	// Create a channel to signal when the server is ready
	serverReady := make(chan bool, 1)
	
	// Create channels for shutdown coordination
	shutdownChan := make(chan bool)
	shutdownComplete := make(chan bool)
	cancelShutdown := make(chan bool)

	// Create server instance
	srv := &http.Server{
		Addr: ":" + port,
	}

	// Start the server in a goroutine
	go func() {
		r := chi.NewRouter()

		// Middleware
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)

		// Static files - serve from ui/static
		workDir, _ := os.Getwd()
		filesDir := http.Dir(filepath.Join(workDir, "ui/static"))
		r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(filesDir)))

		// Routes
		r.Get("/", handleIndex(workDir))
		
		// Add shutdown trigger endpoint
		r.Post("/trigger-shutdown", func(w http.ResponseWriter, r *http.Request) {
			var req ShutdownRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				req.Immediate = false
			}

			w.WriteHeader(http.StatusOK)
			
			if req.Immediate {
				// Trigger immediate shutdown
				shutdownChan <- true
				return
			}

			// Start delayed shutdown in a goroutine
			go func() {
				select {
				case <-time.After(30 * time.Second):
					shutdownChan <- true
				case <-cancelShutdown:
					log.Println("Shutdown cancelled")
					return
				}
			}()
		})

		// Add cancel shutdown endpoint
		r.Post("/cancel-shutdown", func(w http.ResponseWriter, r *http.Request) {
			cancelShutdown <- true
			w.WriteHeader(http.StatusOK)
		})

		srv.Handler = r

		// Signal that we're about to start the server
		serverReady <- true

		log.Printf("Server starting on %s", url)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for server to start
	<-serverReady
	time.Sleep(100 * time.Millisecond) // Give the server a moment to fully start

	// Open the browser
	openBrowser(url)

	// Handle graceful shutdown
	go func() {
		<-shutdownChan
		log.Println("Initiating graceful shutdown...")
		
		// Create a context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Error during server shutdown: %v", err)
		}
		
		log.Println("Server shutdown complete")
		shutdownComplete <- true
	}()

	// Wait for shutdown to complete
	<-shutdownComplete
	log.Println("Application exiting")
	os.Exit(0)
} 
