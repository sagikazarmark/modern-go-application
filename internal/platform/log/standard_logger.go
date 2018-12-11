package log

import (
	"log"

	"github.com/goph/logur"
)

// NewStandardErrorLogger returns a new standard logger logging on error level.
func NewStandardErrorLogger(logger logur.Logger) *log.Logger {
	return logur.NewStandardLogger(logger, logur.Error, "", 0)
}
