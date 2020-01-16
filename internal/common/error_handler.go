package common

import (
	"context"
)

// ErrorHandler handles an error.
type ErrorHandler interface {
	Handle(err error)
	HandleContext(ctx context.Context, err error)
}

// NoopErrorHandler is an error handler that discards every error.
type NoopErrorHandler struct{}

func (NoopErrorHandler) Handle(_ error)                           {}
func (NoopErrorHandler) HandleContext(_ context.Context, _ error) {}
