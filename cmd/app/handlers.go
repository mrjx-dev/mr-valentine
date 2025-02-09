package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// ShutdownRequest represents a request to shutdown the server
type ShutdownRequest struct {
	Immediate bool `json:"immediate"`
}

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

// handleCancelShutdown handles the shutdown cancellation request
func (app *Application) handleCancelShutdown() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.cancel <- true
		w.WriteHeader(http.StatusOK)
	}
}
