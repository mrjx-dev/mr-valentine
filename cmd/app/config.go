package main

import "time"

// Config holds the application configuration
type Config struct {
	// Server configuration
	Server struct {
		Port            string
		ShutdownTimeout time.Duration
		GracePeriod     time.Duration
	}

	// Static file configuration
	Static struct {
		Dir string
	}

	// Template configuration
	Template struct {
		Dir string
	}
}

// NewConfig returns a new Config instance with default values
func NewConfig() *Config {
	cfg := &Config{}

	// Set default values
	cfg.Server.Port = "14334"
	cfg.Server.ShutdownTimeout = 10 * time.Second
	cfg.Server.GracePeriod = 30 * time.Second

	return cfg
}

// BaseURL returns the base URL for the application
func (c *Config) BaseURL() string {
	return "http://localhost:" + c.Server.Port
}
