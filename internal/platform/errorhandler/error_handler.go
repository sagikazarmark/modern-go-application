package errorhandler

import (
	"github.com/goph/emperror"
	"github.com/goph/logur"
)

// New returns a new error handler.
func New(logger logur.ErrorLogger) emperror.Handler {
	return logur.NewErrorHandler(logger)
}
