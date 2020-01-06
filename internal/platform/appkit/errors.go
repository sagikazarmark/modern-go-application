package appkit

import (
	"emperror.dev/errors"
)

// IsClientError checks if an error should be returned to the client for processing.
func IsClientError(err error) bool {
	var clientError interface {
		ClientError() bool
	}

	if errors.As(err, &clientError) {
		return clientError.ClientError()
	}

	return false
}
