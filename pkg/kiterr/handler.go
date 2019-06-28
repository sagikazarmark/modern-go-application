package kiterr

import (
	"context"

	"github.com/go-kit/kit/transport"
	"github.com/goph/emperror"
)

type errorHandler struct {
	handler emperror.Handler
}

// NewHandler returns a new transport error handler.
func NewHandler(handler emperror.Handler) transport.ErrorHandler {
	return &errorHandler{handler}
}

// Handle implements the transport.ErrorHandler interface.
func (h *errorHandler) Handle(ctx context.Context, err error) {
	h.handler.Handle(err)
}
