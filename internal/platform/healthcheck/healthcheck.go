package healthcheck

import (
	"net/http"

	health "github.com/InVisionApp/go-health"
	"github.com/InVisionApp/go-health/handlers"
	"github.com/goph/logur"
	"github.com/goph/logur/integrations/invisionlog"
)

// New returns a new health checker instance.
func New(logger logur.Logger) *health.Health {
	healthChecker := health.New()
	healthChecker.Logger = invisionlog.New(logur.WithFields(logger, map[string]interface{}{"component": "healthcheck"}))

	return healthChecker
}

// Handler returns a new HTTP handler for a health checker.
func Handler(healthChecker health.IHealth) http.Handler {
	return handlers.NewJSONHandlerFunc(healthChecker, nil)
}
