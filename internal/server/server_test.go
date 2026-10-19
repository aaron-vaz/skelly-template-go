package server

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTestServer(t *testing.T) (*Server, *bytes.Buffer) {
	t.Helper()
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, nil))
	return New(logger, time.Second), &logs
}

func TestHandlers(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{name: "liveness", path: "/healthz", wantStatus: http.StatusOK, wantBody: `"status":"ok"`},
		{name: "readiness", path: "/readyz", wantStatus: http.StatusOK, wantBody: `"status":"ok"`},
		{name: "version", path: "/version", wantStatus: http.StatusOK, wantBody: `"version":"dev"`},
		{name: "path parameter", path: "/hello/Gopher", wantStatus: http.StatusOK, wantBody: `"message":"Hello, Gopher!"`},
		{name: "unknown route", path: "/nope", wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			srv, _ := newTestServer(t)
			rec := httptest.NewRecorder()

			// When
			srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))

			// Then
			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			// And
			if !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Errorf("body = %q, want it to contain %q", rec.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestHandlers_MethodNotAllowed(t *testing.T) {
	// Given
	srv, _ := newTestServer(t)
	rec := httptest.NewRecorder()

	// When
	srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/healthz", nil))

	// Then
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestMiddleware_RequestID(t *testing.T) {
	// Given
	srv, logs := newTestServer(t)

	t.Run("generates an id", func(t *testing.T) {
		// When
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

		// Then
		if rec.Header().Get(requestIDHeader) == "" {
			t.Error("Expected a generated request id header")
		}
	})

	t.Run("keeps the id the caller sent and logs the request", func(t *testing.T) {
		// Given
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		req.Header.Set(requestIDHeader, "abc123")
		rec := httptest.NewRecorder()

		// When
		srv.Handler().ServeHTTP(rec, req)

		// Then
		if got := rec.Header().Get(requestIDHeader); got != "abc123" {
			t.Errorf("request id = %q, want abc123", got)
		}

		// And
		if !strings.Contains(logs.String(), "request_id=abc123") || !strings.Contains(logs.String(), "status=200") {
			t.Errorf("Expected an access log line with the id and status, got %q", logs.String())
		}
	})
}

func TestMiddleware_RecoversFromPanic(t *testing.T) {
	// Given
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, nil))
	panicking := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { panic("boom") })
	handler := withMiddleware(logger, panicking)
	rec := httptest.NewRecorder()

	// When
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	// Then
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}

	// And
	if !strings.Contains(logs.String(), "panic while handling request") {
		t.Errorf("Expected the panic to be logged, got %q", logs.String())
	}
}

func TestServer_RunStopsOnCancel(t *testing.T) {
	// Given
	srv, _ := newTestServer(t)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- srv.Run(ctx, ln) }()

	// When it is serving requests
	resp, err := http.Get("http://" + ln.Addr().String() + "/healthz") //nolint:noctx // plain smoke request
	if err != nil {
		t.Fatalf("GET /healthz error = %v", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	// And the context is cancelled
	cancel()

	// Then it shuts down cleanly
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Run() error = %v, want nil", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run() did not return after the context was cancelled")
	}
}
