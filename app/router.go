package app

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// ErrNotFound is returned when no route matches a request.
var ErrNotFound = errors.New("route not found")

// handlerFunc adapts an (w, r) error-returning function to http.HandlerFunc,
// centralizing the status mapping for errors.
func handlerFunc(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

// NewMux builds the service's HTTP router.
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/liveness", handlerFunc(Liveness))
	mux.HandleFunc("GET /v1/readiness", handlerFunc(Readiness))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, ErrNotFound.Error(), http.StatusNotFound)
	})
	return mux
}

// Server wraps http.Server with sane timeouts and graceful shutdown.
type Server struct {
	httpServer *http.Server
}

// NewServer returns a Server listening on addr with conservative timeouts.
func NewServer(addr string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           NewMux(),
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 3 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
	}
}

// Start begins serving. It returns once ListenAndServe exits, which after
// Shutdown is ErrServerClosed — callers should treat that as clean exit.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server within the ctx deadline.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
