package app

import "net/http"

// Liveness reports that the process is alive. It performs no external checks.
func Liveness(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, HealthResponse{Status: "OK"})
}

// Readiness reports that the service can accept traffic. Add dependency pings
// here as they are introduced (database, caches, ...).
func Readiness(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, HealthResponse{Status: "READY"})
}
