// Package business holds domain logic and strong types. It must not import
// the app layer, and no two domains may import each other (see ARCHITECTURE.md).
package business

import (
	"context"
	"time"
)

// Status is a strong type describing service health. It never appears in API
// JSON structs or DB row structs directly — converters translate at crossings.
type Status int

// Supported statuses.
const (
	StatusUnknown Status = iota
	StatusOK
	StatusReady
)

// String implements fmt.Stringer.
func (s Status) String() string {
	switch s {
	case StatusOK:
		return "OK"
	case StatusReady:
		return "READY"
	default:
		return "UNKNOWN"
	}
}

// HealthCheck captures a point-in-time view of service health.
type HealthCheck struct {
	Status    Status
	Timestamp time.Time
}

// Service answers basic domain questions about the service itself.
type Service struct{}

// NewService constructs the core service.
func NewService() *Service {
	return &Service{}
}

// Liveness returns the current liveness state.
func (s *Service) Liveness(ctx context.Context) (HealthCheck, error) {
	return HealthCheck{Status: StatusOK, Timestamp: time.Now()}, nil
}

// Readiness returns the current readiness state. Extend with dependency pings.
func (s *Service) Readiness(ctx context.Context) (HealthCheck, error) {
	return HealthCheck{Status: StatusReady, Timestamp: time.Now()}, nil
}
