// Package server is the HTTP server: routes, handlers, middleware and the
// graceful start/stop lifecycle.
package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"
)

// Server wraps an http.Server with a graceful shutdown.
type Server struct {
	http            *http.Server
	logger          *slog.Logger
	shutdownTimeout time.Duration
}

// Run serves on ln until ctx is cancelled, then stops accepting connections and
// gives in-flight requests up to the shutdown timeout to finish. A server that
// stopped because ctx was cancelled returns nil.
func (s *Server) Run(ctx context.Context, ln net.Listener) error {
	errs := make(chan error, 1)
	go func() {
		s.logger.Info("server listening", "addr", ln.Addr().String())
		errs <- s.http.Serve(ln)
	}()

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
	}

	s.logger.Info("server shutting down")
	// The caller's context is already cancelled, so shutdown needs one of its own
	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), s.shutdownTimeout)
	defer cancel()

	if err := s.http.Shutdown(shutdownCtx); err != nil {
		return err
	}
	if err := <-errs; !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Handler returns the server's root handler.
func (s *Server) Handler() http.Handler {
	return s.http.Handler
}

// New builds a Server with the routes and middleware already wired up.
func New(logger *slog.Logger, shutdownTimeout time.Duration) *Server {
	return &Server{
		http: &http.Server{
			Handler:           withMiddleware(logger, routes()),
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
		logger:          logger,
		shutdownTimeout: shutdownTimeout,
	}
}
