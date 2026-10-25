// Package errs holds the error vocabulary used by the Business layer.
package errs

import (
	"errors"
	"fmt"
)

// FieldError describes a validation failure on a single field.
type FieldError struct {
	Field   string
	Message string
}

// FieldErrors collects validation failures for a request payload.
type FieldErrors map[string]error

// Add records a failure for the given field, creating the map when needed.
func (fe FieldErrors) Add(field, message string) FieldErrors {
	if fe == nil {
		fe = make(FieldErrors)
	}
	fe[field] = errors.New(message)
	return fe
}

// Err returns the collection as a single error, or nil when empty.
func (fe FieldErrors) Err() error {
	if len(fe) == 0 {
		return nil
	}
	return fe
}

// Error implements the error interface.
func (fe FieldErrors) Error() string {
	errs := make([]error, 0, len(fe))
	for _, err := range fe {
		errs = append(errs, err)
	}
	return fmt.Sprintf("field errors: %v", errors.Join(errs...))
}

// New builds an error that wraps the standard error kinds with context.
func New(text string) error {
	return errors.New(text)
}

// Newf builds a formatted error.
func Newf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

// Is reports whether any error in err's chain matches target.
func Is(err, target error) bool { return errors.Is(err, target) }

// As finds the first error in err's chain that matches target.
func As(err error, target any) bool { return errors.As(err, target) }

// Join combines non-nil errors into one.
func Join(errs ...error) error { return errors.Join(errs...) }
