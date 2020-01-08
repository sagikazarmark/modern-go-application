package appkit

import (
	"emperror.dev/errors"
)

// IsClientError checks if an error should be returned to the client for processing.
func IsClientError(err error) bool {
	var e interface {
		ClientError() bool
	}

	if errors.As(err, &e) {
		return e.ClientError()
	}

	return false
}

// IsNotFoundError checks if an error is related to a resource being not found.
func IsNotFoundError(err error) bool {
	var e interface {
		NotFound() bool
	}

	if errors.As(err, &e) {
		return e.NotFound()
	}

	return false
}

// IsValidationError checks if an error is related to a resource being invalid.
func IsValidationError(err error) bool {
	var e interface {
		Validation() bool
	}

	if errors.As(err, &e) {
		return e.Validation()
	}

	return false
}
