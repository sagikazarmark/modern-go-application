package common

import (
	"context"
)

// ErrorHandler handles an error.
type ErrorHandler interface {
	// Handle handles an error.
	Handle(ctx context.Context, err error)
}

type noopErrorHandler struct{}

// NewNoopErrorHandler returns an error handler that discards every error.
func NewNoopErrorHandler() ErrorHandler { return noopErrorHandler{} }

func (noopErrorHandler) Handle(ctx context.Context, err error) {}
