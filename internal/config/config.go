// Package config loads the service configuration from the environment.
package config

import (
	"fmt"
	"log/slog"
	"strings"
	"time"
)

const (
	defaultAddr            = ":8080"
	defaultLogLevel        = "info"
	defaultShutdownTimeout = 10 * time.Second
)

// Config holds everything the service needs to start.
type Config struct {
	// Addr is the address the HTTP server listens on.
	Addr string
	// LogLevel is the minimum level that is logged.
	LogLevel slog.Level
	// ShutdownTimeout is how long in-flight requests get to finish on shutdown.
	ShutdownTimeout time.Duration
}

// Load reads ADDR, LOG_LEVEL and SHUTDOWN_TIMEOUT through getenv, falling back
// to defaults for anything that is unset.
func Load(getenv func(string) string) (Config, error) {
	cfg := Config{
		Addr:            valueOr(getenv("ADDR"), defaultAddr),
		ShutdownTimeout: defaultShutdownTimeout,
	}

	if err := cfg.LogLevel.UnmarshalText([]byte(strings.ToLower(valueOr(getenv("LOG_LEVEL"), defaultLogLevel)))); err != nil {
		return Config{}, fmt.Errorf("invalid LOG_LEVEL: %w", err)
	}

	if raw := getenv("SHUTDOWN_TIMEOUT"); raw != "" {
		timeout, err := time.ParseDuration(raw)
		if err != nil || timeout <= 0 {
			return Config{}, fmt.Errorf("invalid SHUTDOWN_TIMEOUT %q: want a positive duration such as 30s", raw)
		}
		cfg.ShutdownTimeout = timeout
	}

	return cfg, nil
}

func valueOr(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
