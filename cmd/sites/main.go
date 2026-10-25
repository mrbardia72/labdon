// Command sites runs the labdon HTTP service.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"

	"github.com/mrbardia72/labdon/app"
	"github.com/mrbardia72/labdon/foundation/inner/buildinfo"
)

const defaultHTTPPort = 3000

// buildInfo carries values injected via -ldflags at link time.
var build = "dev"

func main() {
	if err := run(); err != nil {
		slog.Error("shutdown", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// ── Parse configuration ──────────────────────────────────────────────
	var traceEnabled bool
	flag.BoolVar(&traceEnabled, "trace", false, "enable tracing")
	flag.Parse()

	var cfg struct {
		Web struct {
			HTTPPort         int
			HTTPReadTimeout  time.Duration
			HTTPWriteTimeout time.Duration
			ShutdownTimeout  time.Duration
		}
	}
	cfg.Web.HTTPPort = defaultHTTPPort
	cfg.Web.HTTPReadTimeout = 5 * time.Second
	cfg.Web.HTTPWriteTimeout = 10 * time.Second
	cfg.Web.ShutdownTimeout = 10 * time.Second

	help, err := conf.Parse("LABDON", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("rendering config: %w", err)
	}
	slog.Info("startup", "build", buildinfo.Get().Version, "injected", build, "config", out)

	// ── Start HTTP service ───────────────────────────────────────────────
	addr := fmt.Sprintf(":%d", cfg.Web.HTTPPort)
	server := app.NewServer(addr)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("service starting", "addr", addr)
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	// ── Wait for termination signal ──────────────────────────────────────
	select {
	case err := <-serverErr:
		if err != nil {
			return fmt.Errorf("starting server: %w", err)
		}
	case sig := <-shutdown:
		slog.Info("shutdown started", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}
		if err := <-serverErr; err != nil {
			return fmt.Errorf("server exited with error: %w", err)
		}
		slog.Info("shutdown complete")
	}

	return nil
}
