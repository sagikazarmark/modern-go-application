package health

import (
	"github.com/InVisionApp/go-health"
	"github.com/go-kit/kit/log"
)

// Config is a type alias to make package importing easier.
type Config = health.Config

// New returns a new health checker.
func New(logger log.Logger) *health.Health {
	h := health.New()

	h.Logger = &loggerShim{logger}

	return h
}
