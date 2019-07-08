package kiterr

import (
	"context"

	"emperror.dev/emperror"
	"github.com/go-kit/kit/transport"
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
