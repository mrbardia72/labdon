// Package app is the HTTP layer. Request/response structs here use primitive
// JSON types only; strong domain types never cross this boundary directly —
// they arrive through named converters (see ARCHITECTURE.md).
package app

import (
	"encoding/json"
	"net/http"
)

// HealthResponse is the primitive JSON payload returned by the health checks.
type HealthResponse struct {
	Status string `json:"status"`
}

// WriteJSON encodes v as a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}
