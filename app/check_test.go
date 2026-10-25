package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoints(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string
	}{
		{name: "liveness ok", method: http.MethodGet, path: "/v1/liveness", wantStatus: http.StatusOK, wantBody: "OK"},
		{name: "readiness ok", method: http.MethodGet, path: "/v1/readiness", wantStatus: http.StatusOK, wantBody: "READY"},
		{name: "wrong method 404", method: http.MethodPost, path: "/v1/liveness", wantStatus: http.StatusNotFound},
		{name: "unknown route 404", method: http.MethodGet, path: "/nope", wantStatus: http.StatusNotFound},
	}

	mux := NewMux()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequestWithContext(t.Context(), tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if got := rec.Result().StatusCode; got != tt.wantStatus {
				t.Fatalf("status code: got %d want %d", got, tt.wantStatus)
			}
			if tt.wantBody == "" {
				return
			}

			var resp HealthResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("decoding response: %v", err)
			}
			if resp.Status != tt.wantBody {
				t.Errorf("status: got %q want %q", resp.Status, tt.wantBody)
			}
		})
	}
}
