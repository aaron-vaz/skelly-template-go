package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

const requestIDHeader = "X-Request-ID"

type loggerKey struct{}

// loggerFrom returns the request-scoped logger set by the logging middleware,
// so every line logged while handling a request carries its request id.
func loggerFrom(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// withMiddleware wraps next so that panics are recovered, and every request
// gets an id and one access log line. The order matters: logging is outermost
// so that a recovered panic is still logged with its 500.
func withMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return withLogging(logger, withRecovery(next))
}

func withLogging(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		id := r.Header.Get(requestIDHeader)
		if id == "" {
			id = newRequestID()
		}
		w.Header().Set(requestIDHeader, id)

		reqLogger := logger.With("request_id", id)
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r.WithContext(context.WithValue(r.Context(), loggerKey{}, reqLogger)))

		reqLogger.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", float64(time.Since(start).Microseconds())/1000,
		)
	})
}

func withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			// http.ErrAbortHandler is how a handler deliberately aborts a response,
			// and net/http expects it to keep propagating
			if err, ok := rec.(error); ok && errors.Is(err, http.ErrAbortHandler) {
				panic(rec)
			}
			loggerFrom(r.Context()).Error("panic while handling request", "panic", rec)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}()
		next.ServeHTTP(w, r)
	})
}

// statusRecorder remembers the status code so it can be logged.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func newRequestID() string {
	buf := make([]byte, 8)
	// crypto/rand.Read never returns an error on supported platforms
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}
