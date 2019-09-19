package common

import (
	"context"
)

// ErrorHandler handles an error.
type ErrorHandler interface {
	// Handle handles an error.
	Handle(ctx context.Context, err error)
}
