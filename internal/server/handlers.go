package server

import (
	"encoding/json"
	"net/http"

	"github.com/{{.Inputs.owner.Value}}/{{.Inputs.project_name.Value}}/internal/version"
)

// routes registers every endpoint. Add new ones here.
func routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handleHealth)
	mux.HandleFunc("GET /readyz", handleHealth)
	mux.HandleFunc("GET /version", handleVersion)
	mux.HandleFunc("GET /hello/{name}", handleHello)

	return mux
}

// handleHealth backs both the liveness and readiness probes. Make /readyz check
// the service's dependencies once it has any.
func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, r, http.StatusOK, map[string]string{"status": "ok"})
}

func handleVersion(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, r, http.StatusOK, map[string]string{"version": version.Version})
}

// handleHello is an example endpoint using a path parameter; replace it with your own.
func handleHello(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, r, http.StatusOK, map[string]string{"message": "Hello, " + r.PathValue("name") + "!"})
}

func writeJSON(w http.ResponseWriter, r *http.Request, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		// The status line is already sent, so all that is left is to record it
		loggerFrom(r.Context()).Error("failed to write response", "error", err)
	}
}
